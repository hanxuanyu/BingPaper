package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/repo"
	"BingPaper/internal/storage"
	"BingPaper/internal/util"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type BingResponse struct {
	Images []BingImage `json:"images"`
}

type BingImage struct {
	Startdate     string `json:"startdate"`
	Fullstartdate string `json:"fullstartdate"`
	Enddate       string `json:"enddate"`
	URL           string `json:"url"`
	URLBase       string `json:"urlbase"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`
	Title         string `json:"title"`
	Quiz          string `json:"quiz"`
}

type Fetcher struct {
	httpClient *http.Client
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *Fetcher) Fetch(ctx context.Context, n int) error {
	util.Logger.Info("Starting fetch task", zap.Int("n", n))
	url := fmt.Sprintf("%s?format=js&idx=0&n=%d&uhd=1&mkt=%s", config.BingAPIBase, n, config.BingMkt)
	util.Logger.Debug("Requesting Bing API", zap.String("url", url))
	resp, err := f.httpClient.Get(url)
	if err != nil {
		util.Logger.Error("Failed to request Bing API", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	var bingResp BingResponse
	if err := json.NewDecoder(resp.Body).Decode(&bingResp); err != nil {
		util.Logger.Error("Failed to decode Bing API response", zap.Error(err))
		return err
	}

	util.Logger.Info("Fetched images from Bing", zap.Int("count", len(bingResp.Images)))

	for _, bingImg := range bingResp.Images {
		if err := f.processImage(ctx, bingImg); err != nil {
			util.Logger.Error("Failed to process image", zap.String("date", bingImg.Enddate), zap.Error(err))
		}
	}

	util.Logger.Info("Fetch task completed")
	return nil
}

func (f *Fetcher) processImage(ctx context.Context, bingImg BingImage) error {
	dateStr := fmt.Sprintf("%s-%s-%s", bingImg.Enddate[0:4], bingImg.Enddate[4:6], bingImg.Enddate[6:8])

	// 幂等检查
	var existing model.Image
	if err := repo.DB.Where("date = ?", dateStr).First(&existing).Error; err == nil {
		util.Logger.Info("Image already exists, skipping", zap.String("date", dateStr))
		return nil
	}

	util.Logger.Info("Processing new image", zap.String("date", dateStr), zap.String("title", bingImg.Title))

	// UHD 探测
	imgURL, variantName := f.probeUHD(bingImg.URLBase)

	imgData, err := f.downloadImage(imgURL)
	if err != nil {
		util.Logger.Error("Failed to download image", zap.String("url", imgURL), zap.Error(err))
		return err
	}

	// 解码图片用于缩放
	srcImg, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		util.Logger.Error("Failed to decode image data", zap.Error(err))
		return err
	}

	// 创建 DB 记录
	dbImg := model.Image{
		Date:      dateStr,
		Title:     bingImg.Title,
		Copyright: bingImg.Copyright,
		URLBase:   bingImg.URLBase,
		Quiz:      bingImg.Quiz,
	}

	if err := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}},
		DoNothing: true,
	}).Create(&dbImg).Error; err != nil {
		util.Logger.Error("Failed to create image record", zap.Error(err))
		return err
	}

	// 再次检查 dbImg.ID 是否被填充，如果没有被填充（说明由于冲突未插入），则需要查询出已有的 ID
	if dbImg.ID == 0 {
		var existing model.Image
		if err := repo.DB.Where("date = ?", dateStr).First(&existing).Error; err != nil {
			util.Logger.Error("Failed to query existing image record after conflict", zap.Error(err))
			return err
		}
		dbImg = existing
	}

	// 保存各种分辨率
	variants := []struct {
		name   string
		width  int
		height int
	}{
		{variantName, 0, 0}, // 原图 (UHD 或 1080p)
		{"1920x1080", 1920, 1080},
		{"1366x768", 1366, 768},
	}

	for _, v := range variants {
		// 如果是探测到的最高清版本，且我们已经有了数据，直接使用
		var currentImgData []byte
		if v.width == 0 {
			currentImgData = imgData
		} else {
			resized := imaging.Fill(srcImg, v.width, v.height, imaging.Center, imaging.Lanczos)
			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 90}); err != nil {
				util.Logger.Warn("Failed to encode jpeg", zap.String("variant", v.name), zap.Error(err))
				continue
			}
			currentImgData = buf.Bytes()
		}

		// 保存 JPG
		if err := f.saveVariant(ctx, &dbImg, v.name, "jpg", currentImgData); err != nil {
			util.Logger.Error("Failed to save variant", zap.String("variant", v.name), zap.Error(err))
		}

	}

	// 保存今日额外文件
	today := time.Now().Format("2006-01-02")
	if dateStr == today && config.GetConfig().Feature.WriteDailyFiles {
		f.saveDailyFiles(srcImg, imgData)
	}

	return nil
}

func (f *Fetcher) probeUHD(urlBase string) (string, string) {
	uhdURL := fmt.Sprintf("https://www.bing.com%s_UHD.jpg", urlBase)
	resp, err := f.httpClient.Head(uhdURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		return uhdURL, "UHD"
	}
	return fmt.Sprintf("https://www.bing.com%s_1920x1080.jpg", urlBase), "1920x1080"
}

func (f *Fetcher) downloadImage(url string) ([]byte, error) {
	resp, err := f.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (f *Fetcher) saveVariant(ctx context.Context, img *model.Image, variant, format string, data []byte) error {
	key := fmt.Sprintf("%s/%s_%s.%s", img.Date, img.Date, variant, format)
	contentType := "image/jpeg"
	if format == "webp" {
		contentType = "image/webp"
	}

	stored, err := storage.GlobalStorage.Put(ctx, key, bytes.NewReader(data), contentType)
	if err != nil {
		return err
	}

	vRecord := model.ImageVariant{
		ImageID:    img.ID,
		Variant:    variant,
		Format:     format,
		StorageKey: stored.Key,
		PublicURL:  stored.PublicURL,
		Size:       int64(len(data)),
	}

	return repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "image_id"}, {Name: "variant"}, {Name: "format"}},
		DoNothing: true,
	}).Create(&vRecord).Error
}

func (f *Fetcher) saveDailyFiles(srcImg image.Image, originalData []byte) {
	util.Logger.Info("Saving daily files")
	localRoot := config.GetConfig().Storage.Local.Root
	if config.GetConfig().Storage.Type != "local" {
		// 如果不是本地存储，保存在静态资源目录
		localRoot = "static"
	}

	if err := os.MkdirAll(localRoot, 0755); err != nil {
		util.Logger.Error("Failed to create directory", zap.String("path", localRoot), zap.Error(err))
		return
	}

	// daily.jpeg (quality 95)
	jpegPath := filepath.Join(localRoot, "daily.jpeg")
	fJpeg, err := os.Create(jpegPath)
	if err != nil {
		util.Logger.Error("Failed to create daily.jpeg", zap.Error(err))
	} else {
		jpeg.Encode(fJpeg, srcImg, &jpeg.Options{Quality: 95})
		fJpeg.Close()
	}

	// original.jpeg (quality 100)
	originalPath := filepath.Join(localRoot, "original.jpeg")
	if err := os.WriteFile(originalPath, originalData, 0644); err != nil {
		util.Logger.Error("Failed to write original.jpeg", zap.Error(err))
	}
}

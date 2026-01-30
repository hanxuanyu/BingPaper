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
	"strings"
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
	HSH           string `json:"hsh"`
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
	regions := config.GetConfig().Fetcher.Regions
	if len(regions) == 0 {
		regions = []string{config.GetConfig().GetDefaultMkt()}
	}

	for _, mkt := range regions {
		if err := f.FetchRegion(ctx, mkt); err != nil {
			util.Logger.Error("Failed to fetch region images", zap.String("mkt", mkt), zap.Error(err))
		}
	}

	util.Logger.Info("Fetch task completed")
	return nil
}

// FetchRegion 抓取指定地区的图片
func (f *Fetcher) FetchRegion(ctx context.Context, mkt string) error {
	if !util.IsValidRegion(mkt) {
		util.Logger.Warn("Skipping fetch for invalid region", zap.String("mkt", mkt))
		return fmt.Errorf("invalid region code: %s", mkt)
	}
	util.Logger.Info("Fetching images for region", zap.String("mkt", mkt))
	// 调用两次 API 获取最多两周的数据
	// 第一次 idx=0&n=8 (今天起往回数 8 张)
	if err := f.fetchByMkt(ctx, mkt, 0, 8); err != nil {
		util.Logger.Error("Failed to fetch images", zap.String("mkt", mkt), zap.Int("idx", 0), zap.Error(err))
		return err
	}
	// 第二次 idx=7&n=8 (7天前起往回数 8 张，与第一次有重叠，确保不漏)
	if err := f.fetchByMkt(ctx, mkt, 7, 8); err != nil {
		util.Logger.Error("Failed to fetch images", zap.String("mkt", mkt), zap.Int("idx", 7), zap.Error(err))
		// 第二次失败不一定返回错误，因为可能第一次已经拿到了
	}
	return nil
}

func (f *Fetcher) fetchByMkt(ctx context.Context, mkt string, idx int, n int) error {
	url := fmt.Sprintf("%s?format=js&idx=%d&n=%d&uhd=1&mkt=%s", config.BingAPIBase, idx, n, mkt)
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

	util.Logger.Info("Fetched images from Bing", zap.String("mkt", mkt), zap.Int("count", len(bingResp.Images)))

	for _, bingImg := range bingResp.Images {
		if err := f.processImage(ctx, bingImg, mkt); err != nil {
			util.Logger.Error("Failed to process image", zap.String("date", bingImg.Enddate), zap.String("mkt", mkt), zap.Error(err))
		}
	}

	return nil
}

func (f *Fetcher) processImage(ctx context.Context, bingImg BingImage, mkt string) error {
	dateStr := fmt.Sprintf("%s-%s-%s", bingImg.Enddate[0:4], bingImg.Enddate[4:6], bingImg.Enddate[6:8])

	// 幂等检查
	var existing model.Image
	if err := repo.DB.Where("date = ? AND mkt = ?", dateStr, mkt).First(&existing).Error; err == nil {
		util.Logger.Debug("Image already exists in DB, skipping", zap.String("date", dateStr), zap.String("mkt", mkt))
		return nil
	}

	imageName := f.extractImageName(bingImg.URLBase, bingImg.HSH)
	util.Logger.Info("Processing image", zap.String("date", dateStr), zap.String("mkt", mkt), zap.String("imageName", imageName))

	// 创建 DB 记录
	dbImg := model.Image{
		Date:          dateStr,
		Mkt:           mkt,
		Title:         bingImg.Title,
		Copyright:     bingImg.Copyright,
		CopyrightLink: bingImg.CopyrightLink,
		URLBase:       bingImg.URLBase,
		Quiz:          bingImg.Quiz,
		StartDate:     bingImg.Startdate,
		FullStartDate: bingImg.Fullstartdate,
		HSH:           bingImg.HSH,
	}

	if err := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}, {Name: "mkt"}},
		DoNothing: true,
	}).Create(&dbImg).Error; err != nil {
		util.Logger.Error("Failed to create image record", zap.Error(err))
		return err
	}

	if dbImg.ID == 0 {
		var existing model.Image
		if err := repo.DB.Where("date = ? AND mkt = ?", dateStr, mkt).First(&existing).Error; err != nil {
			util.Logger.Error("Failed to query existing image record after conflict", zap.Error(err))
			return err
		}
		dbImg = existing
	}

	// UHD 探测
	imgURL, variantName := f.probeUHD(bingImg.URLBase)

	// 保存各种分辨率
	targetVariants := []struct {
		name   string
		width  int
		height int
	}{
		{"1920x1080", 1920, 1080},
		{"1366x768", 1366, 768},
		{"1280x720", 1280, 720},
		{"1024x768", 1024, 768},
		{"800x600", 800, 600},
		{"800x480", 800, 480},
		{"640x480", 640, 480},
		{"640x360", 640, 360},
		{"480x360", 480, 360},
		{"400x240", 400, 240},
		{"320x240", 320, 240},
	}

	// 检查是否所有变体都已存在于存储中
	allExist := true
	// 检查 UHD/原图
	uhdKey := f.generateKey(imageName, variantName, "jpg")
	exists, _ := storage.GlobalStorage.Exists(ctx, uhdKey)
	if !exists {
		allExist = false
	} else {
		for _, v := range targetVariants {
			if v.name == variantName {
				continue
			}
			vKey := f.generateKey(imageName, v.name, "jpg")
			exists, _ := storage.GlobalStorage.Exists(ctx, vKey)
			if !exists {
				allExist = false
				break
			}
		}
	}

	if allExist {
		util.Logger.Debug("All image variants exist in storage, linking only", zap.String("imageName", imageName))
		// 只建立关联信息
		f.saveVariant(ctx, &dbImg, imageName, variantName, "jpg", nil)
		for _, v := range targetVariants {
			if v.name == variantName {
				continue
			}
			f.saveVariant(ctx, &dbImg, imageName, v.name, "jpg", nil)
		}
	} else {
		// 需要下载并处理
		util.Logger.Debug("Downloading and processing image", zap.String("url", imgURL))
		imgData, err := f.downloadImage(imgURL)
		if err != nil {
			util.Logger.Error("Failed to download image", zap.String("url", imgURL), zap.Error(err))
			return err
		}

		srcImg, _, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			util.Logger.Error("Failed to decode image data", zap.Error(err))
			return err
		}

		// 保存原图
		if err := f.saveVariant(ctx, &dbImg, imageName, variantName, "jpg", imgData); err != nil {
			util.Logger.Error("Failed to save original variant", zap.String("variant", variantName), zap.Error(err))
		}

		for _, v := range targetVariants {
			if v.name == variantName {
				continue
			}
			resized := imaging.Fill(srcImg, v.width, v.height, imaging.Center, imaging.Lanczos)
			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 100}); err != nil {
				util.Logger.Warn("Failed to encode jpeg", zap.String("variant", v.name), zap.Error(err))
				continue
			}
			currentImgData := buf.Bytes()
			if err := f.saveVariant(ctx, &dbImg, imageName, v.name, "jpg", currentImgData); err != nil {
				util.Logger.Error("Failed to save variant", zap.String("variant", v.name), zap.Error(err))
			}
		}

		// 保存今日额外文件
		today := time.Now().Format("2006-01-02")
		if dateStr == today && config.GetConfig().Feature.WriteDailyFiles {
			f.saveDailyFiles(srcImg, imgData, mkt)
		}
	}

	return nil
}

func (f *Fetcher) extractImageName(urlBase, hsh string) string {
	// 示例: /th?id=OHR.MilwaukeeHall_ROW0871854348
	start := 0
	if idx := strings.Index(urlBase, "OHR."); idx != -1 {
		start = idx + 4
	} else if idx := strings.Index(urlBase, "id="); idx != -1 {
		start = idx + 3
	}

	rem := urlBase[start:]
	end := strings.Index(rem, "_")
	if end == -1 {
		end = len(rem)
	}

	name := rem[:end]
	if name == "" {
		return hsh
	}
	return name
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

func (f *Fetcher) generateKey(imageName, variant, format string) string {
	return fmt.Sprintf("%s/%s_%s.%s", imageName, imageName, variant, format)
}

func (f *Fetcher) saveVariant(ctx context.Context, img *model.Image, imageName, variant, format string, data []byte) error {
	key := f.generateKey(imageName, variant, format)
	contentType := "image/jpeg"
	if format == "webp" {
		contentType = "image/webp"
	}

	var size int64
	var publicURL string

	exists, _ := storage.GlobalStorage.Exists(ctx, key)
	if exists {
		util.Logger.Debug("Variant already exists in storage, linking", zap.String("key", key))
		// 如果存在，我们需要获取它的大小和公共 URL (如果可能)
		// 但目前的 Storage 接口没有 Stat，我们可以尝试 Get 或者干脆 size 为 0
		// 为了简单，我们只从存储中获取公共 URL
		if pURL, ok := storage.GlobalStorage.PublicURL(key); ok {
			publicURL = pURL
		}
		// size 暂时设为 0 或者从 data 中取 (如果有的话)
		if data != nil {
			size = int64(len(data))
		}
	} else if data != nil {
		util.Logger.Debug("Saving variant to storage", zap.String("key", key))
		stored, err := storage.GlobalStorage.Put(ctx, key, bytes.NewReader(data), contentType)
		if err != nil {
			return err
		}
		publicURL = stored.PublicURL
		size = stored.Size
	} else {
		return fmt.Errorf("variant %s does not exist and no data provided", key)
	}

	vRecord := model.ImageVariant{
		ImageID:    img.ID,
		Variant:    variant,
		Format:     format,
		StorageKey: key,
		PublicURL:  publicURL,
		Size:       size,
	}

	return repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "image_id"}, {Name: "variant"}, {Name: "format"}},
		DoNothing: true,
	}).Create(&vRecord).Error
}

func (f *Fetcher) saveDailyFiles(srcImg image.Image, originalData []byte, mkt string) {
	util.Logger.Info("Saving daily files", zap.String("mkt", mkt))
	localRoot := config.GetConfig().Storage.Local.Root
	if localRoot == "" {
		localRoot = "data"
	}

	mktDir := filepath.Join(localRoot, mkt)
	if err := os.MkdirAll(mktDir, 0755); err != nil {
		util.Logger.Error("Failed to create directory", zap.String("path", mktDir), zap.Error(err))
		return
	}

	// daily.jpeg (quality 100)
	jpegPath := filepath.Join(mktDir, "daily.jpeg")
	fJpeg, err := os.Create(jpegPath)
	if err != nil {
		util.Logger.Error("Failed to create daily.jpeg", zap.Error(err))
	} else {
		jpeg.Encode(fJpeg, srcImg, &jpeg.Options{Quality: 100})
		fJpeg.Close()
	}

	// original.jpeg (quality 100)
	originalPath := filepath.Join(mktDir, "original.jpeg")
	if err := os.WriteFile(originalPath, originalData, 0644); err != nil {
		util.Logger.Error("Failed to write original.jpeg", zap.Error(err))
	}

	// 同时也保留一份在根目录下（兼容旧逻辑，且作为默认地区图片）
	// 如果是默认地区或者是第一个抓取的地区，可以覆盖根目录的文件
	if mkt == config.GetConfig().GetDefaultMkt() {
		jpegPathRoot := filepath.Join(localRoot, "daily.jpeg")
		fJpegRoot, err := os.Create(jpegPathRoot)
		if err == nil {
			jpeg.Encode(fJpegRoot, srcImg, &jpeg.Options{Quality: 100})
			fJpegRoot.Close()
		}
		originalPathRoot := filepath.Join(localRoot, "original.jpeg")
		os.WriteFile(originalPathRoot, originalData, 0644)
	}
}

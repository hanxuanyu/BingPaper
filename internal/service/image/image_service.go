package image

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/repo"
	"BingPaper/internal/service/fetcher"
	"BingPaper/internal/storage"
	"BingPaper/internal/util"

	"go.uber.org/zap"
)

var ErrFetchStarted = errors.New("on-demand fetch started")

func CleanupOldImages(ctx context.Context) error {
	days := config.GetConfig().Retention.Days
	if days <= 0 {
		return nil
	}

	threshold := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	util.Logger.Info("Starting cleanup task", zap.Int("retention_days", days), zap.String("threshold", threshold))

	var images []model.Image
	if err := repo.DB.Where("date < ?", threshold).Preload("Variants").Find(&images).Error; err != nil {
		util.Logger.Error("Failed to query old images for cleanup", zap.Error(err))
		return err
	}

	for _, img := range images {
		util.Logger.Info("Deleting old image", zap.String("date", img.Date))
		for _, v := range img.Variants {
			if err := storage.GlobalStorage.Delete(ctx, v.StorageKey); err != nil {
				util.Logger.Warn("Failed to delete storage object", zap.String("key", v.StorageKey), zap.Error(err))
			}
		}
		// 删除关联记录（逻辑外键控制）
		if err := repo.DB.Where("image_id = ?", img.ID).Delete(&model.ImageVariant{}).Error; err != nil {
			util.Logger.Error("Failed to delete variants", zap.Uint("image_id", img.ID), zap.Error(err))
		}
		// 删除主表记录
		if err := repo.DB.Delete(&img).Error; err != nil {
			util.Logger.Error("Failed to delete image", zap.Uint("id", img.ID), zap.Error(err))
		}
	}

	util.Logger.Info("Cleanup task completed", zap.Int("deleted_count", len(images)))
	return nil
}

func GetTodayImage(mkt string) (*model.Image, error) {
	today := time.Now().Format("2006-01-02")
	util.Logger.Debug("Getting today image", zap.String("mkt", mkt), zap.String("today", today))
	var img model.Image
	tx := repo.DB.Where("date = ?", today)
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	err := tx.Preload("Variants").First(&img).Error
	if err != nil && mkt != "" && config.GetConfig().API.EnableOnDemandFetch && util.IsValidRegion(mkt) {
		// 如果没找到，尝试异步按需抓取该地区
		util.Logger.Info("Image not found in DB, starting asynchronous on-demand fetch", zap.String("mkt", mkt))
		f := fetcher.NewFetcher()
		go func() {
			_ = f.FetchRegion(context.Background(), mkt)
		}()
		return nil, ErrFetchStarted
	}

	if err != nil {
		util.Logger.Debug("Today image not found, trying latest image", zap.String("mkt", mkt))
		// 如果今天还是没有，尝试获取最近的一张
		tx = repo.DB.Order("date desc")
		if mkt != "" {
			tx = tx.Where("mkt = ?", mkt)
		}
		err = tx.Preload("Variants").First(&img).Error
	}

	// 兜底逻辑：如果指定地区没找到，且开启了兜底开关，则尝试获取默认地区的图片
	if err != nil && mkt != "" && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultMkt()
		util.Logger.Debug("Image not found, trying fallback to default market", zap.String("mkt", mkt), zap.String("defaultMkt", defaultMkt))
		if mkt != defaultMkt {
			return GetTodayImage(defaultMkt)
		}
		return GetTodayImage("")
	}

	if err == nil {
		util.Logger.Debug("Found image", zap.String("date", img.Date), zap.String("mkt", img.Mkt))
	}
	return &img, err
}

func GetAllRegionsTodayImages() ([]model.Image, error) {
	regions := config.GetConfig().Fetcher.Regions
	if len(regions) == 0 {
		regions = []string{config.GetConfig().GetDefaultMkt()}
	}

	var images []model.Image
	for _, mkt := range regions {
		img, err := GetTodayImage(mkt)
		if err == nil {
			images = append(images, *img)
		}
	}
	return images, nil
}

func GetRandomImage(mkt string) (*model.Image, error) {
	util.Logger.Debug("Getting random image", zap.String("mkt", mkt))
	var img model.Image
	// SQLite 使用 RANDOM(), MySQL/Postgres 使用 RANDOM() 或 RAND()
	// 简单起见，先查总数再 Offset
	var count int64
	tx := repo.DB.Model(&model.Image{})
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	tx.Count(&count)
	if count == 0 && mkt != "" && config.GetConfig().API.EnableOnDemandFetch && util.IsValidRegion(mkt) {
		// 如果没找到，尝试异步按需抓取该地区
		util.Logger.Info("No images found in DB for region, starting asynchronous on-demand fetch", zap.String("mkt", mkt))
		f := fetcher.NewFetcher()
		go func() {
			_ = f.FetchRegion(context.Background(), mkt)
		}()
		return nil, ErrFetchStarted
	}

	if count == 0 {
		return nil, fmt.Errorf("no images found")
	}

	// 优化随机查询：使用 Offset 代替 ORDER BY RANDOM()
	// 注意：tx 包含了前面的 Where 条件
	offset := rand.Intn(int(count))
	util.Logger.Debug("Random image selection", zap.Int64("total", count), zap.Int("offset", offset))
	err := tx.Preload("Variants").Offset(offset).Limit(1).Find(&img).Error

	// 兜底逻辑
	if (err != nil || img.ID == 0) && mkt != "" && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultMkt()
		util.Logger.Debug("Random image not found, trying fallback", zap.String("mkt", mkt), zap.String("defaultMkt", defaultMkt))
		if mkt != defaultMkt {
			return GetRandomImage(defaultMkt)
		}
		return GetRandomImage("")
	}

	if err == nil && img.ID == 0 {
		return nil, fmt.Errorf("no images found")
	}

	if err == nil {
		util.Logger.Debug("Found random image", zap.String("date", img.Date), zap.String("mkt", img.Mkt))
	}

	return &img, err
}

func GetImageByDate(date string, mkt string) (*model.Image, error) {
	util.Logger.Debug("Getting image by date", zap.String("date", date), zap.String("mkt", mkt))
	var img model.Image
	tx := repo.DB.Where("date = ?", date)
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	err := tx.Preload("Variants").First(&img).Error
	if err != nil && mkt != "" && config.GetConfig().API.EnableOnDemandFetch && util.IsValidRegion(mkt) {
		// 如果没找到，尝试异步按需抓取该地区
		util.Logger.Info("Image not found in DB for date, starting asynchronous on-demand fetch", zap.String("mkt", mkt), zap.String("date", date))
		f := fetcher.NewFetcher()
		go func() {
			_ = f.FetchRegion(context.Background(), mkt)
		}()
		return nil, ErrFetchStarted
	}

	// 兜底逻辑
	if err != nil && mkt != "" && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultMkt()
		util.Logger.Debug("Image by date not found, trying fallback", zap.String("date", date), zap.String("mkt", mkt), zap.String("defaultMkt", defaultMkt))
		if mkt != defaultMkt {
			return GetImageByDate(date, defaultMkt)
		}
		return GetImageByDate(date, "")
	}

	if err == nil {
		util.Logger.Debug("Found image by date", zap.String("date", img.Date), zap.String("mkt", img.Mkt))
	}
	return &img, err
}

func GetImageList(limit int, offset int, month string, mkt string) ([]model.Image, error) {
	util.Logger.Debug("Getting image list", zap.Int("limit", limit), zap.Int("offset", offset), zap.String("month", month), zap.String("mkt", mkt))
	var images []model.Image
	tx := repo.DB.Model(&model.Image{})

	if month != "" {
		// 增强过滤：确保只处理 YYYY-MM 格式，防止注入或非法字符
		// 这里简单处理：只要不为空就增加 LIKE 过滤
		util.Logger.Debug("Filtering images by month", zap.String("month", month))
		tx = tx.Where("date LIKE ?", month+"%")
	}

	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}

	tx = tx.Order("date desc").Preload("Variants")

	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if offset > 0 {
		tx = tx.Offset(offset)
	}

	err := tx.Find(&images).Error
	if err != nil {
		util.Logger.Error("Failed to get image list", zap.Error(err), zap.String("month", month))
	}
	return images, err
}

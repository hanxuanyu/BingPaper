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
	"gorm.io/gorm"
)

var ErrFetchStarted = errors.New("on-demand fetch started")

func CleanupOldImages(ctx context.Context) error {
	days := config.GetConfig().Retention.Days
	if days <= 0 {
		return nil
	}

	threshold := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	util.Logger.Info("Starting cleanup task", zap.Int("retention_days", days), zap.String("threshold", threshold))

	var regionRecords []model.ImageRegion
	if err := repo.DB.Where("date < ?", threshold).Preload("Variants").Find(&regionRecords).Error; err != nil {
		util.Logger.Error("Failed to query old image regions for cleanup", zap.Error(err))
		return err
	}

	for _, m := range regionRecords {
		util.Logger.Info("Deleting old image region record", zap.String("date", m.Date), zap.String("mkt", m.Mkt))

		// 检查该图片名是否还有其他地区或日期在使用
		var count int64
		repo.DB.Model(&model.ImageRegion{}).Where("image_name = ? AND id != ?", m.ImageName, m.ID).Count(&count)

		if count == 0 {
			util.Logger.Info("Image content no longer referenced, deleting files and variants", zap.String("image_name", m.ImageName))
			for _, v := range m.Variants {
				if err := storage.GlobalStorage.Delete(ctx, v.StorageKey); err != nil {
					util.Logger.Warn("Failed to delete storage object", zap.String("key", v.StorageKey), zap.Error(err))
				}
			}
			// 删除变体记录
			if err := repo.DB.Where("image_name = ?", m.ImageName).Delete(&model.ImageVariant{}).Error; err != nil {
				util.Logger.Error("Failed to delete variants", zap.String("image_name", m.ImageName), zap.Error(err))
			}
		}

		// 删除地区记录
		if err := repo.DB.Delete(&m).Error; err != nil {
			util.Logger.Error("Failed to delete image region record", zap.Uint("id", m.ID), zap.Error(err))
		}
	}

	util.Logger.Info("Cleanup task completed", zap.Int("deleted_count", len(regionRecords)))
	return nil
}

func GetTodayImage(mkt string) (*model.ImageRegion, error) {
	if mkt == "" {
		mkt = config.GetConfig().GetDefaultRegion()
	}
	today := time.Now().Format("2006-01-02")
	util.Logger.Debug("Getting today image", zap.String("mkt", mkt), zap.String("today", today))
	var imgRegion model.ImageRegion
	tx := repo.DB.Where("date = ? AND mkt = ?", today, mkt)
	err := tx.Preload("Variants", func(db *gorm.DB) *gorm.DB {
		return db.Order("size asc")
	}).First(&imgRegion).Error
	if err != nil && config.GetConfig().API.EnableOnDemandFetch && util.IsValidRegion(mkt) {
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
		err = repo.DB.Where("mkt = ?", mkt).Order("date desc").Preload("Variants", func(db *gorm.DB) *gorm.DB {
			return db.Order("size asc")
		}).First(&imgRegion).Error
	}

	// 兜底逻辑
	if err != nil && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultRegion()
		util.Logger.Debug("Image not found, trying fallback to default region", zap.String("mkt", mkt), zap.String("defaultMkt", defaultMkt))
		if mkt != defaultMkt {
			return GetTodayImage(defaultMkt)
		}
	}

	if err == nil {
		util.Logger.Debug("Found image region record", zap.String("date", imgRegion.Date), zap.String("mkt", imgRegion.Mkt))
	}
	return &imgRegion, err
}

func GetAllRegionsTodayImages() ([]model.ImageRegion, error) {
	today := time.Now().Format("2006-01-02")
	regions := config.GetConfig().Fetcher.Regions
	if len(regions) == 0 {
		regions = []string{config.GetConfig().GetDefaultRegion()}
	}

	var images []model.ImageRegion
	err := repo.DB.Where("date = ? AND mkt IN ?", today, regions).
		Preload("Variants", func(db *gorm.DB) *gorm.DB {
			return db.Order("size asc")
		}).Find(&images).Error

	if err != nil {
		return nil, err
	}

	// 按照配置的 regions 顺序排序
	mktMap := make(map[string]model.ImageRegion)
	for _, img := range images {
		mktMap[img.Mkt] = img
	}

	var sortedImages []model.ImageRegion
	for _, r := range regions {
		if img, ok := mktMap[r]; ok {
			sortedImages = append(sortedImages, img)
		}
	}

	return sortedImages, nil
}

func GetRandomImage(mkt string) (*model.ImageRegion, error) {
	if mkt == "" {
		mkt = config.GetConfig().GetDefaultRegion()
	}
	util.Logger.Debug("Getting random image", zap.String("mkt", mkt))
	var imgRegion model.ImageRegion
	var count int64
	tx := repo.DB.Model(&model.ImageRegion{}).Where("mkt = ?", mkt)
	tx.Count(&count)
	if count == 0 && config.GetConfig().API.EnableOnDemandFetch && util.IsValidRegion(mkt) {
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

	offset := rand.Intn(int(count))
	util.Logger.Debug("Random image selection", zap.Int64("total", count), zap.Int("offset", offset))
	err := tx.Preload("Variants", func(db *gorm.DB) *gorm.DB {
		return db.Order("size asc")
	}).Offset(offset).Limit(1).Find(&imgRegion).Error

	if (err != nil || imgRegion.ID == 0) && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultRegion()
		util.Logger.Debug("Random image not found, trying fallback", zap.String("mkt", mkt), zap.String("defaultMkt", defaultMkt))
		if mkt != defaultMkt {
			return GetRandomImage(defaultMkt)
		}
	}

	if err == nil && imgRegion.ID == 0 {
		return nil, fmt.Errorf("no images found")
	}

	return &imgRegion, err
}

func GetImageByDate(date string, mkt string) (*model.ImageRegion, error) {
	if mkt == "" {
		mkt = config.GetConfig().GetDefaultRegion()
	}
	util.Logger.Debug("Getting image by date", zap.String("date", date), zap.String("mkt", mkt))
	var imgRegion model.ImageRegion
	err := repo.DB.Where("date = ? AND mkt = ?", date, mkt).Preload("Variants", func(db *gorm.DB) *gorm.DB {
		return db.Order("size asc")
	}).First(&imgRegion).Error
	if err != nil && config.GetConfig().API.EnableOnDemandFetch && util.IsValidRegion(mkt) {
		util.Logger.Info("Image not found in DB for date, starting asynchronous on-demand fetch", zap.String("mkt", mkt), zap.String("date", date))
		f := fetcher.NewFetcher()
		go func() {
			_ = f.FetchRegion(context.Background(), mkt)
		}()
		return nil, ErrFetchStarted
	}

	if err != nil && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultRegion()
		if mkt != defaultMkt {
			return GetImageByDate(date, defaultMkt)
		}
	}

	return &imgRegion, err
}

func GetImageList(limit int, offset int, month string, mkt string) ([]model.ImageRegion, error) {
	if mkt == "" {
		mkt = config.GetConfig().GetDefaultRegion()
	}
	var images []model.ImageRegion
	tx := repo.DB.Model(&model.ImageRegion{}).Where("mkt = ?", mkt)

	if month != "" {
		tx = tx.Where("date LIKE ?", month+"%")
	}

	tx = tx.Order("date desc").Preload("Variants", func(db *gorm.DB) *gorm.DB {
		return db.Order("size asc")
	})

	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if offset > 0 {
		tx = tx.Offset(offset)
	}

	err := tx.Find(&images).Error
	return images, err
}

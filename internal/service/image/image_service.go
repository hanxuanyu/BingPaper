package image

import (
	"context"
	"fmt"
	"time"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/repo"
	"BingPaper/internal/storage"
	"BingPaper/internal/util"

	"go.uber.org/zap"
)

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
	var img model.Image
	tx := repo.DB.Where("date = ?", today)
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	err := tx.Preload("Variants").First(&img).Error
	if err != nil {
		// 如果今天没有，尝试获取最近的一张
		tx = repo.DB.Order("date desc")
		if mkt != "" {
			tx = tx.Where("mkt = ?", mkt)
		}
		err = tx.Preload("Variants").First(&img).Error
	}

	// 兜底逻辑：如果指定地区没找到，且开启了兜底开关，则尝试获取默认地区的图片
	if err != nil && mkt != "" && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultMkt()
		if mkt != defaultMkt {
			return GetTodayImage(defaultMkt)
		}
		return GetTodayImage("")
	}

	return &img, err
}

func GetRandomImage(mkt string) (*model.Image, error) {
	var img model.Image
	// SQLite 使用 RANDOM(), MySQL/Postgres 使用 RANDOM() 或 RAND()
	// 简单起见，先查总数再 Offset
	var count int64
	tx := repo.DB.Model(&model.Image{})
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	tx.Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("no images found")
	}

	// 这种方法不适合海量数据，但对于 30 天的数据没问题
	tx = repo.DB.Order("RANDOM()")
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	err := tx.Preload("Variants").First(&img).Error
	if err != nil {
		// 适配 MySQL
		tx = repo.DB.Order("RAND()")
		if mkt != "" {
			tx = tx.Where("mkt = ?", mkt)
		}
		err = tx.Preload("Variants").First(&img).Error
	}

	// 兜底逻辑
	if err != nil && mkt != "" && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultMkt()
		if mkt != defaultMkt {
			return GetRandomImage(defaultMkt)
		}
		return GetRandomImage("")
	}

	return &img, err
}

func GetImageByDate(date string, mkt string) (*model.Image, error) {
	var img model.Image
	tx := repo.DB.Where("date = ?", date)
	if mkt != "" {
		tx = tx.Where("mkt = ?", mkt)
	}
	err := tx.Preload("Variants").First(&img).Error

	// 兜底逻辑
	if err != nil && mkt != "" && config.GetConfig().API.EnableMktFallback {
		defaultMkt := config.GetConfig().GetDefaultMkt()
		if mkt != defaultMkt {
			return GetImageByDate(date, defaultMkt)
		}
		return GetImageByDate(date, "")
	}

	return &img, err
}

func GetImageList(limit int, offset int, month string, mkt string) ([]model.Image, error) {
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

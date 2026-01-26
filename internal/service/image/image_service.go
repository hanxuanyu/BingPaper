package image

import (
	"context"
	"fmt"
	"time"

	"BingDailyImage/internal/config"
	"BingDailyImage/internal/model"
	"BingDailyImage/internal/repo"
	"BingDailyImage/internal/storage"
	"BingDailyImage/internal/util"

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
		return err
	}

	for _, img := range images {
		util.Logger.Info("Deleting old image", zap.String("date", img.Date))
		for _, v := range img.Variants {
			if err := storage.GlobalStorage.Delete(ctx, v.StorageKey); err != nil {
				util.Logger.Warn("Failed to delete storage object", zap.String("key", v.StorageKey), zap.Error(err))
			}
		}
		// 删除 DB 记录 (级联删除由代码处理，或者 GORM 会处理已加载的关联吗？)
		// 简单起见，手动删除关联
		repo.DB.Where("image_id = ?", img.ID).Delete(&model.ImageVariant{})
		repo.DB.Delete(&img)
	}

	util.Logger.Info("Cleanup task completed", zap.Int("deleted_count", len(images)))
	return nil
}

func GetTodayImage() (*model.Image, error) {
	today := time.Now().Format("2006-01-02")
	var img model.Image
	err := repo.DB.Where("date = ?", today).Preload("Variants").First(&img).Error
	if err != nil {
		// 如果今天没有，尝试获取最近的一张
		err = repo.DB.Order("date desc").Preload("Variants").First(&img).Error
	}
	return &img, err
}

func GetRandomImage() (*model.Image, error) {
	var img model.Image
	// SQLite 使用 RANDOM(), MySQL/Postgres 使用 RANDOM() 或 RAND()
	// 简单起见，先查总数再 Offset
	var count int64
	repo.DB.Model(&model.Image{}).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("no images found")
	}

	// 这种方法不适合海量数据，但对于 30 天的数据没问题
	err := repo.DB.Order("RANDOM()").Preload("Variants").First(&img).Error
	if err != nil {
		// 适配 MySQL
		err = repo.DB.Order("RAND()").Preload("Variants").First(&img).Error
	}
	return &img, err
}

func GetImageByDate(date string) (*model.Image, error) {
	var img model.Image
	err := repo.DB.Where("date = ?", date).Preload("Variants").First(&img).Error
	return &img, err
}

func GetImageList(limit int) ([]model.Image, error) {
	var images []model.Image
	db := repo.DB.Order("date desc").Preload("Variants")
	if limit > 0 {
		db = db.Limit(limit)
	}
	err := db.Find(&images).Error
	return images, err
}

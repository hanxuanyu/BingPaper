package repo

import (
	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/util"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateDataToNewDB 将数据从旧数据库迁移到新数据库
func MigrateDataToNewDB(oldDB *gorm.DB, newConfig *config.Config) error {
	util.Logger.Info("Starting data migration to new database",
		zap.String("new_type", newConfig.DB.Type),
		zap.String("new_dsn", newConfig.DB.DSN))

	// 1. 初始化新数据库连接
	dialector, err := GetDialector(newConfig.DB.Type, newConfig.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to get dialector for new DB: %w", err)
	}

	gormConfig := GetGormConfig(newConfig)
	newDB, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to new DB: %w", err)
	}

	// 2. 自动迁移结构
	if err := newDB.AutoMigrate(&model.ImageRegion{}, &model.ImageVariant{}, &model.Token{}); err != nil {
		return fmt.Errorf("failed to migrate schema in new DB: %w", err)
	}

	// 3. 清空新数据库中的现有数据（防止冲突）
	util.Logger.Info("Cleaning up destination database before migration")
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ImageVariant{}).Error; err != nil {
		return fmt.Errorf("failed to clear ImageVariants: %w", err)
	}
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ImageRegion{}).Error; err != nil {
		return fmt.Errorf("failed to clear ImageRegions: %w", err)
	}
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Token{}).Error; err != nil {
		return fmt.Errorf("failed to clear Tokens: %w", err)
	}

	// 4. 开始迁移数据
	// 使用事务确保迁移的原子性
	return newDB.Transaction(func(tx *gorm.DB) error {
		// 迁移 ImageRegions
		var regions []model.ImageRegion
		if err := oldDB.Find(&regions).Error; err != nil {
			return fmt.Errorf("failed to fetch image regions from old DB: %w", err)
		}
		if len(regions) > 0 {
			util.Logger.Info("Migrating image regions", zap.Int("count", len(regions)))
			if err := tx.Create(&regions).Error; err != nil {
				return fmt.Errorf("failed to insert image regions into new DB: %w", err)
			}
		}

		// 迁移 ImageVariants
		var variants []model.ImageVariant
		if err := oldDB.Find(&variants).Error; err != nil {
			return fmt.Errorf("failed to fetch variants from old DB: %w", err)
		}
		if len(variants) > 0 {
			util.Logger.Info("Migrating variants", zap.Int("count", len(variants)))
			if err := tx.Create(&variants).Error; err != nil {
				return fmt.Errorf("failed to insert variants into new DB: %w", err)
			}
		}

		// 迁移 Tokens
		var tokens []model.Token
		if err := oldDB.Find(&tokens).Error; err != nil {
			return fmt.Errorf("failed to fetch tokens from old DB: %w", err)
		}
		if len(tokens) > 0 {
			util.Logger.Info("Migrating tokens", zap.Int("count", len(tokens)))
			if err := tx.Create(&tokens).Error; err != nil {
				return fmt.Errorf("failed to insert tokens into new DB: %w", err)
			}
		}

		// 更新全局 DB 指针
		DB = newDB
		util.Logger.Info("Data migration completed successfully")
		return nil
	})
}

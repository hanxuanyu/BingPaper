package repo

import (
	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/util"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MigrationStats struct {
	ImageRegions  int `json:"image_regions"`
	ImageVariants int `json:"image_variants"`
	Tokens        int `json:"tokens"`
	ApiStats      int `json:"api_stats"`
}

var migrationMu sync.Mutex

func migrateTable[T any](source *gorm.DB, target *gorm.DB, modelName string) (int, error) {
	var rows []T
	if err := source.Unscoped().Order("id asc").Find(&rows).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch %s from source DB: %w", modelName, err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	util.Logger.Info("Migrating table",
		zap.String("model", modelName),
		zap.Int("count", len(rows)))

	if err := target.CreateInBatches(&rows, 200).Error; err != nil {
		return 0, fmt.Errorf("failed to insert %s into target DB: %w", modelName, err)
	}
	return len(rows), nil
}

// MigrateDataToNewDB 将数据从当前运行中的数据库迁移到目标数据库。
// 该过程只迁移数据，不会切换当前服务的活动数据库连接。
func MigrateDataToNewDB(oldDB *gorm.DB, baseCfg *config.Config, targetDB config.DBConfig) (MigrationStats, error) {
	var stats MigrationStats
	if oldDB == nil {
		return stats, fmt.Errorf("source database is not initialized")
	}

	migrationMu.Lock()
	defer migrationMu.Unlock()

	util.Logger.Info("Starting data migration to new database",
		zap.String("new_type", targetDB.Type),
		zap.String("new_dsn", targetDB.DSN))

	// 1. 初始化新数据库连接
	newDB, err := openDB(BuildDBRuntimeConfig(baseCfg, targetDB))
	if err != nil {
		return stats, fmt.Errorf("failed to connect to target DB: %w", err)
	}

	sqlDB, err := newDB.DB()
	if err != nil {
		return stats, fmt.Errorf("failed to get target SQL DB: %w", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return stats, fmt.Errorf("failed to ping target DB: %w", err)
	}

	// 2. 自动迁移结构
	if err := AutoMigrateModels(newDB); err != nil {
		return stats, fmt.Errorf("failed to migrate schema in target DB: %w", err)
	}

	// 3. 清空新数据库中的现有数据（防止冲突）
	util.Logger.Info("Cleaning up destination database before migration")
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ImageVariant{}).Error; err != nil {
		return stats, fmt.Errorf("failed to clear ImageVariants: %w", err)
	}
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ImageRegion{}).Error; err != nil {
		return stats, fmt.Errorf("failed to clear ImageRegions: %w", err)
	}
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Token{}).Error; err != nil {
		return stats, fmt.Errorf("failed to clear Tokens: %w", err)
	}
	if err := newDB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ApiStat{}).Error; err != nil {
		return stats, fmt.Errorf("failed to clear ApiStats: %w", err)
	}

	// 4. 开始迁移数据
	// 使用事务确保迁移的原子性
	if err := newDB.Transaction(func(tx *gorm.DB) error {
		stats.ImageRegions, err = migrateTable[model.ImageRegion](oldDB, tx, "ImageRegion")
		if err != nil {
			return err
		}

		stats.ImageVariants, err = migrateTable[model.ImageVariant](oldDB, tx, "ImageVariant")
		if err != nil {
			return err
		}

		stats.Tokens, err = migrateTable[model.Token](oldDB, tx, "Token")
		if err != nil {
			return err
		}

		stats.ApiStats, err = migrateTable[model.ApiStat](oldDB, tx, "ApiStat")
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return stats, err
	}

	util.Logger.Info("Data migration completed successfully",
		zap.Int("image_regions", stats.ImageRegions),
		zap.Int("image_variants", stats.ImageVariants),
		zap.Int("tokens", stats.Tokens),
		zap.Int("api_stats", stats.ApiStats))

	return stats, nil
}

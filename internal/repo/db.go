package repo

import (
	"BingDailyImage/internal/config"
	"BingDailyImage/internal/model"
	"BingDailyImage/internal/util"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	cfg := config.GetConfig()
	var dialector gorm.Dialector

	switch cfg.DB.Type {
	case "mysql":
		dialector = mysql.Open(cfg.DB.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DB.DSN)
	case "sqlite":
		dialector = sqlite.Open(cfg.DB.DSN)
	default:
		return fmt.Errorf("unsupported db type: %s", cfg.DB.Type)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return err
	}

	// 迁移
	if err := db.AutoMigrate(&model.Image{}, &model.ImageVariant{}, &model.Token{}); err != nil {
		return err
	}

	DB = db
	util.Logger.Info("Database initialized successfully", zap.String("type", cfg.DB.Type))
	return nil
}

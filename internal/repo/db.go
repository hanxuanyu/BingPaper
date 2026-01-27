package repo

import (
	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/util"
	"fmt"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return err
	}

	// 针对 MySQL 的额外处理：如果数据库不存在，GORM 的 mysql 驱动通常无法直接创建库。
	// 但此处假设 DSN 中指定的数据库已经存在。AutoMigrate 会负责创建表。

	// 迁移
	if err := db.AutoMigrate(&model.Image{}, &model.ImageVariant{}, &model.Token{}); err != nil {
		util.Logger.Error("Database migration failed", zap.Error(err))
		return err
	}

	DB = db
	util.Logger.Info("Database initialized successfully", zap.String("type", cfg.DB.Type))
	return nil
}

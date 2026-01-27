package repo

import (
	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/util"
	"context"
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type gormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{
		ZapLogger: l.ZapLogger,
		LogLevel:  level,
	}
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil && l.LogLevel >= logger.Error {
		l.ZapLogger.Error("SQL ERROR",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	} else if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn("SLOW SQL",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	} else if l.LogLevel >= logger.Info {
		l.ZapLogger.Info("SQL",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}

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

	gormLogLevel := logger.Info
	switch cfg.Log.DBLogLevel {
	case "debug":
		gormLogLevel = logger.Info // GORM 的 Info 级会输出所有 SQL
	case "info":
		gormLogLevel = logger.Info
	case "warn":
		gormLogLevel = logger.Warn
	case "error":
		gormLogLevel = logger.Error
	case "silent":
		gormLogLevel = logger.Silent
	}

	gormConfig := &gorm.Config{
		Logger: &gormLogger{
			ZapLogger: util.DBLogger,
			LogLevel:  gormLogLevel,
		},
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

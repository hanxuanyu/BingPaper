package cron

import (
	"context"

	"BingPaper/internal/config"
	"BingPaper/internal/service/fetcher"
	"BingPaper/internal/service/image"
	"BingPaper/internal/util"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var GlobalCron *cron.Cron

func InitCron() {
	cfg := config.GetConfig()
	if !cfg.Cron.Enabled {
		util.Logger.Info("Cron is disabled")
		return
	}

	c := cron.New()

	// 每日抓取任务
	_, err := c.AddFunc(cfg.Cron.DailySpec, func() {
		util.Logger.Info("Running scheduled daily fetch")
		f := fetcher.NewFetcher()
		if err := f.Fetch(context.Background(), 1); err != nil {
			util.Logger.Error("Scheduled fetch failed", zap.Error(err))
		}

		// 抓取后顺便清理
		if err := image.CleanupOldImages(context.Background()); err != nil {
			util.Logger.Error("Scheduled cleanup failed", zap.Error(err))
		}
	})

	if err != nil {
		util.Logger.Fatal("Failed to setup cron", zap.Error(err))
	}

	c.Start()
	GlobalCron = c
	util.Logger.Info("Cron service started", zap.String("spec", cfg.Cron.DailySpec))
}

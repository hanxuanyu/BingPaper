package main

import (
	"fmt"

	"BingPaper/internal/bootstrap"
	"BingPaper/internal/config"
	"BingPaper/internal/util"

	"go.uber.org/zap"
)

// @title BingPaper API
// @version 1.0
// @description 必应每日一图抓取、存储、管理与公共 API 服务。
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 1. 初始化
	r := bootstrap.Init()

	// 2. 输出欢迎信息
	bootstrap.LogWelcomeInfo()

	// 3. 启动服务
	cfg := config.GetConfig()
	util.Logger.Info("Server starting", zap.Int("port", cfg.Server.Port))
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		util.Logger.Fatal("Server failed to start", zap.Error(err))
	}
}

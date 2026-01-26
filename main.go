package main

import (
	"embed"
	"flag"
	"fmt"
	"mime"

	"BingPaper/internal/bootstrap"
	"BingPaper/internal/config"
	"BingPaper/internal/util"

	"go.uber.org/zap"
)

//go:embed all:web
var webFS embed.FS

// @title BingPaper API
// @version 1.0
// @description 必应每日一图抓取、存储、管理与公共 API 服务。
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.StringVar(&configPath, "c", "", "配置文件路径 (简写)")
	flag.Parse()

	// 注册常用 MIME 类型，确保嵌入式资源能被正确识别
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".svg", "image/svg+xml")

	// 1. 初始化
	r := bootstrap.Init(webFS, configPath)

	// 2. 输出欢迎信息
	bootstrap.LogWelcomeInfo()

	// 3. 启动服务
	cfg := config.GetConfig()
	util.Logger.Info("Server starting", zap.Int("port", cfg.Server.Port))
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		util.Logger.Fatal("Server failed to start", zap.Error(err))
	}
}

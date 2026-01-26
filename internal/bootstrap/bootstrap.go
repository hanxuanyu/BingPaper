package bootstrap

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"BingPaper/internal/config"
	"BingPaper/internal/cron"
	"BingPaper/internal/http"
	"BingPaper/internal/repo"
	"BingPaper/internal/service/fetcher"
	"BingPaper/internal/storage"
	"BingPaper/internal/storage/local"
	"BingPaper/internal/storage/s3"
	"BingPaper/internal/storage/webdav"
	"BingPaper/internal/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Init 初始化应用各项服务
func Init(webFS embed.FS) *gin.Engine {
	// 0. 确保数据目录存在
	_ = os.MkdirAll("data/picture", 0755)

	// 1. 初始化配置
	if err := config.Init(""); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	cfg := config.GetConfig()

	// 2. 初始化日志
	util.InitLogger(cfg.Log.Level)

	// 3. 初始化数据库
	if err := repo.InitDB(); err != nil {
		util.Logger.Fatal("Failed to initialize database")
	}

	// 4. 初始化存储
	var s storage.Storage
	var err error
	switch cfg.Storage.Type {
	case "s3":
		s, err = s3.NewS3Storage(
			cfg.Storage.S3.Endpoint,
			cfg.Storage.S3.Region,
			cfg.Storage.S3.Bucket,
			cfg.Storage.S3.AccessKey,
			cfg.Storage.S3.SecretKey,
			cfg.Storage.S3.PublicURLPrefix,
			cfg.Storage.S3.ForcePathStyle,
		)
	case "webdav":
		s, err = webdav.NewWebDAVStorage(
			cfg.Storage.WebDAV.URL,
			cfg.Storage.WebDAV.Username,
			cfg.Storage.WebDAV.Password,
			cfg.Storage.WebDAV.PublicURLPrefix,
		)
	default: // local
		s, err = local.NewLocalStorage(cfg.Storage.Local.Root)
	}

	if err != nil {
		util.Logger.Fatal("Failed to initialize storage", zap.Error(err))
	}
	storage.GlobalStorage = s

	// 5. 初始化定时任务
	cron.InitCron()

	// 6. 启动时执行一次抓取 (可选，这里我们默认执行一次以确保有数据)
	go func() {
		f := fetcher.NewFetcher()
		f.Fetch(context.Background(), config.BingFetchN)
	}()

	// 7. 设置路由
	return http.SetupRouter(webFS)
}

// LogWelcomeInfo 输出欢迎信息和快速跳转地址
func LogWelcomeInfo() {
	cfg := config.GetConfig()
	port := cfg.Server.Port
	baseURL := cfg.Server.BaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%d", port)
	}

	fmt.Println("\n---------------------------------------------------------")
	fmt.Println("  BingPaper 服务已启动！")
	fmt.Printf("  - 首页地址:   %s/\n", baseURL)
	fmt.Printf("  - 管理后台:   %s/admin\n", baseURL)
	fmt.Printf("  - API 文档:   %s/swagger/index.html\n", baseURL)
	fmt.Printf("  - 今日图片:   %s/api/v1/image/today\n", baseURL)
	fmt.Println("---------------------------------------------------------")
}

package http

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "BingPaper/docs"
	"BingPaper/internal/config"
	"BingPaper/internal/http/handlers"
	"BingPaper/internal/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(webFS embed.FS) *gin.Engine {
	r := gin.Default()

	// CORS 配置：更宽松的配置以解决 Vue 等前端的预检请求问题
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// 公共接口
		img := api.Group("/image")
		{
			img.GET("/today", handlers.GetToday)
			img.GET("/today/meta", handlers.GetTodayMeta)
			img.GET("/random", handlers.GetRandom)
			img.GET("/random/meta", handlers.GetRandomMeta)
			img.GET("/date/:date", handlers.GetByDate)
			img.GET("/date/:date/meta", handlers.GetByDateMeta)
		}
		api.GET("/images", handlers.ListImages)

		// 管理接口
		admin := api.Group("/admin")
		{
			admin.POST("/login", handlers.AdminLogin)

			// 需要验证的接口
			authorized := admin.Group("/")
			authorized.Use(middleware.AuthMiddleware())
			{
				authorized.GET("/tokens", handlers.ListTokens)
				authorized.POST("/tokens", handlers.CreateToken)
				authorized.PATCH("/tokens/:id", handlers.UpdateToken)
				authorized.DELETE("/tokens/:id", handlers.DeleteToken)

				authorized.POST("/password", handlers.ChangePassword)

				authorized.GET("/config", handlers.GetConfig)
				authorized.PUT("/config", handlers.UpdateConfig)

				authorized.POST("/fetch", handlers.ManualFetch)
				authorized.POST("/cleanup", handlers.ManualCleanup)
			}
		}
	}

	// 静态资源服务与 SPA 路由 (放在最后，确保 API 路由优先)
	webSub, _ := fs.Sub(webFS, "web")

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果请求的是 API 或 Swagger，则不处理静态资源 (让其返回 404)
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/swagger") {
			return
		}

		// 辅助函数：尝试从外部或嵌入服务文件
		serveFile := func(relPath string, allowExternal bool) bool {
			// 1. 优先尝试外部路径
			webPath := config.GetConfig().Web.Path
			if allowExternal && webPath != "" {
				fullPath := filepath.Join(webPath, relPath)
				if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
					c.File(fullPath)
					return true
				}
			}

			// 2. 尝试嵌入式文件
			f, err := webSub.Open(relPath)
			if err == nil {
				defer f.Close()
				stat, err := f.Stat()
				if err == nil && !stat.IsDir() {
					if rs, ok := f.(io.ReadSeeker); ok {
						http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), rs)
						return true
					}
					// 兜底：直接读取并输出
					data, err := io.ReadAll(f)
					if err == nil {
						c.Data(http.StatusOK, "", data)
						return true
					}
				}
			}
			return false
		}

		// 1. 尝试直接请求的文件 (如果是 / 则尝试 index.html)
		requestedPath := strings.TrimPrefix(path, "/")
		if requestedPath == "" {
			requestedPath = "index.html"
		}

		if serveFile(requestedPath, true) {
			return
		}

		// 2. SPA 支持：对于非文件请求（没有后缀或不包含点），尝试返回 index.html
		isAsset := strings.Contains(requestedPath, ".")
		if !isAsset || requestedPath == "index.html" {
			if serveFile("index.html", true) {
				return
			}
		}

		c.Status(http.StatusNotFound)
	})

	return r
}

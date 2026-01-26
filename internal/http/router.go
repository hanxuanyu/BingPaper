package http

import (
	_ "BingPaper/docs"
	"BingPaper/internal/http/handlers"
	"BingPaper/internal/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态文件
	r.Static("/static", "./static")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/admin", "./web/index.html")
	r.StaticFile("/login", "./web/index.html")

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

	return r
}

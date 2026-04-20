package http

import (
	"strings"

	docs "BingPaper/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwaggerUIHandler() gin.HandlerFunc {
	swaggerUI := ginSwagger.WrapHandler(swaggerFiles.Handler)

	return func(c *gin.Context) {
		if c.Param("any") == "/doc.json" {
			SwaggerDocHandler(c)
			return
		}

		swaggerUI(c)
	}
}

func SwaggerDocHandler(c *gin.Context) {
	spec := *docs.SwaggerInfo
	spec.Host = forwardedValue(c.GetHeader("X-Forwarded-Host"))
	if spec.Host == "" {
		spec.Host = c.Request.Host
	}

	spec.Schemes = []string{requestScheme(c)}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, spec.ReadDoc())
}

func requestScheme(c *gin.Context) string {
	if proto := forwardedValue(c.GetHeader("X-Forwarded-Proto")); proto != "" {
		return proto
	}
	if scheme := forwardedValue(c.GetHeader("X-Forwarded-Scheme")); scheme != "" {
		return scheme
	}
	if scheme := forwardedValue(c.GetHeader("X-Url-Scheme")); scheme != "" {
		return scheme
	}
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}

func forwardedValue(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}

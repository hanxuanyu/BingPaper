package middleware

import (
	"BingPaper/internal/service/stat"
	"github.com/gin-gonic/gin"
)

// StatMiddleware 记录 API 调用统计的中间件
func StatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 继续执行后续 Handler
		c.Next()

		// 只统计 GET 请求且状态码为正常响应 (200, 302, 304, 202)
		// 202 是按需抓取启动时的返回码
		status := c.Writer.Status()
		if c.Request.Method != "GET" || (status != 200 && status != 302 && status != 304 && status != 202) {
			return
		}

		endpoint := c.FullPath()
		if endpoint == "" {
			return
		}

		// 提取地区信息
		mkt := c.Query("mkt")

		// 记录统计
		stat.RecordStat(endpoint, mkt)
	}
}

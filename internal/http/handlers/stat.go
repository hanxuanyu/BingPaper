package handlers

import (
	"BingPaper/internal/service/stat"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetStatSummary 获取统计概览
// @Summary 获取统计概览
// @Description 获取总调用量、今日调用量、昨日调用量
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /admin/stats/summary [get]
func GetStatSummary(c *gin.Context) {
	data, err := stat.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetStatTrend 获取趋势
// @Summary 获取调用趋势
// @Description 获取最近 N 天的每日调用量。支持按 endpoint 过滤。
// @Tags admin
// @Security BearerAuth
// @Param days query int false "天数" default(7)
// @Param endpoint query string false "接口路径"
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /admin/stats/trend [get]
func GetStatTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	endpoint := c.Query("endpoint")
	data, err := stat.GetTrend(days, endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetStatEndpoints 获取接口分布
// @Summary 获取接口分布
// @Description 获取各接口的调用量分布
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /admin/stats/endpoints [get]
func GetStatEndpoints(c *gin.Context) {
	data, err := stat.GetEndpointDist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetStatRegions 获取地区分布
// @Summary 获取地区分布
// @Description 获取各地区的调用量分布
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /admin/stats/regions [get]
func GetStatRegions(c *gin.Context) {
	data, err := stat.GetRegionDist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

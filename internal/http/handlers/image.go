package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/service/image"
	"BingPaper/internal/storage"

	"github.com/gin-gonic/gin"
)

// GetToday 获取今日图片
// @Summary 获取今日图片
// @Description 根据参数返回今日必应图片流或重定向
// @Tags image
// @Param variant query string false "分辨率 (UHD, 1920x1080, 1366x768)" default(UHD)
// @Param format query string false "格式 (jpg)" default(jpg)
// @Produce image/jpeg
// @Success 200 {file} binary
// @Router /image/today [get]
func GetToday(c *gin.Context) {
	img, err := image.GetTodayImage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	handleImageResponse(c, img)
}

// GetTodayMeta 获取今日图片元数据
// @Summary 获取今日图片元数据
// @Description 获取今日必应图片的标题、版权等元数据
// @Tags image
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /image/today/meta [get]
func GetTodayMeta(c *gin.Context) {
	img, err := image.GetTodayImage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, formatMeta(img))
}

// GetRandom 获取随机图片
// @Summary 获取随机图片
// @Description 随机返回一张已抓取的图片流或重定向
// @Tags image
// @Param variant query string false "分辨率" default(UHD)
// @Param format query string false "格式" default(jpg)
// @Produce image/jpeg
// @Success 200 {file} binary
// @Router /image/random [get]
func GetRandom(c *gin.Context) {
	img, err := image.GetRandomImage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	handleImageResponse(c, img)
}

// GetRandomMeta 获取随机图片元数据
// @Summary 获取随机图片元数据
// @Description 随机获取一张已抓取图片的元数据
// @Tags image
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /image/random/meta [get]
func GetRandomMeta(c *gin.Context) {
	img, err := image.GetRandomImage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, formatMeta(img))
}

// GetByDate 获取指定日期图片
// @Summary 获取指定日期图片
// @Description 根据日期返回图片流或重定向 (yyyy-mm-dd)
// @Tags image
// @Param date path string true "日期 (yyyy-mm-dd)"
// @Param variant query string false "分辨率" default(UHD)
// @Param format query string false "格式" default(jpg)
// @Produce image/jpeg
// @Success 200 {file} binary
// @Router /image/date/{date} [get]
func GetByDate(c *gin.Context) {
	date := c.Param("date")
	img, err := image.GetImageByDate(date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	handleImageResponse(c, img)
}

// GetByDateMeta 获取指定日期图片元数据
// @Summary 获取指定日期图片元数据
// @Description 根据日期获取图片元数据 (yyyy-mm-dd)
// @Tags image
// @Param date path string true "日期 (yyyy-mm-dd)"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /image/date/{date}/meta [get]
func GetByDateMeta(c *gin.Context) {
	date := c.Param("date")
	img, err := image.GetImageByDate(date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, formatMeta(img))
}

// ListImages 获取图片列表
// @Summary 获取图片列表
// @Description 分页获取已抓取的图片元数据列表
// @Tags image
// @Param limit query int false "限制数量" default(30)
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /images [get]
func ListImages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "30")
	var limit int
	fmt.Sscanf(limitStr, "%d", &limit)

	images, err := image.GetImageList(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := []gin.H{}
	for _, img := range images {
		result = append(result, formatMeta(&img))
	}
	c.JSON(http.StatusOK, result)
}

func handleImageResponse(c *gin.Context, img *model.Image) {
	variant := c.DefaultQuery("variant", "UHD")
	format := c.DefaultQuery("format", "jpg")

	var selected *model.ImageVariant
	for _, v := range img.Variants {
		if v.Variant == variant && v.Format == format {
			selected = &v
			break
		}
	}

	if selected == nil && len(img.Variants) > 0 {
		// 回退逻辑
		selected = &img.Variants[0]
	}

	if selected == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "variant not found"})
		return
	}

	mode := config.GetConfig().API.Mode
	if mode == "redirect" {
		if selected.PublicURL != "" {
			c.Redirect(http.StatusFound, selected.PublicURL)
		} else if img.URLBase != "" {
			// 兜底重定向到原始 Bing
			bingURL := fmt.Sprintf("https://www.bing.com%s_%s.jpg", img.URLBase, selected.Variant)
			c.Redirect(http.StatusFound, bingURL)
		} else {
			serveLocal(c, selected.StorageKey)
		}
	} else {
		serveLocal(c, selected.StorageKey)
	}
}

func serveLocal(c *gin.Context, key string) {
	reader, contentType, err := storage.GlobalStorage.Get(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get image"})
		return
	}
	defer reader.Close()

	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	io.Copy(c.Writer, reader)
}

func formatMeta(img *model.Image) gin.H {
	cfg := config.GetConfig()
	variants := []gin.H{}
	for _, v := range img.Variants {
		url := v.PublicURL
		if url == "" && cfg.API.Mode == "redirect" && img.URLBase != "" {
			url = fmt.Sprintf("https://www.bing.com%s_%s.jpg", img.URLBase, v.Variant)
		} else if cfg.API.Mode == "local" || url == "" {
			url = fmt.Sprintf("%s/api/v1/image/date/%s?variant=%s&format=%s", cfg.Server.BaseURL, img.Date, v.Variant, v.Format)
		}
		variants = append(variants, gin.H{
			"variant":     v.Variant,
			"format":      v.Format,
			"size":        v.Size,
			"url":         url,
			"storage_key": v.StorageKey,
		})
	}

	return gin.H{
		"date":      img.Date,
		"title":     img.Title,
		"copyright": img.Copyright,
		"quiz":      img.Quiz,
		"variants":  variants,
	}
}

package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/service/image"
	"BingPaper/internal/storage"
	"BingPaper/internal/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageVariantResp struct {
	Variant    string `json:"variant"`
	Format     string `json:"format"`
	Size       int64  `json:"size"`
	URL        string `json:"url"`
	StorageKey string `json:"storage_key"`
}

type ImageMetaResp struct {
	Date          string             `json:"date"`
	Title         string             `json:"title"`
	Copyright     string             `json:"copyright"`
	CopyrightLink string             `json:"copyrightlink"`
	Quiz          string             `json:"quiz"`
	StartDate     string             `json:"startdate"`
	FullStartDate string             `json:"fullstartdate"`
	HSH           string             `json:"hsh"`
	Variants      []ImageVariantResp `json:"variants"`
}

// GetToday 获取今日图片
// @Summary 获取今日图片
// @Description 根据参数返回今日必应图片流或重定向
// @Tags image
// @Param variant query string false "分辨率 (UHD, 1920x1080, 1366x768, 1280x720, 1024x768, 800x600, 800x480, 640x480, 640x360, 480x360, 400x240, 320x240)" default(UHD)
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
	handleImageResponse(c, img, 7200) // 2小时
}

// GetTodayMeta 获取今日图片元数据
// @Summary 获取今日图片元数据
// @Description 获取今日必应图片的标题、版权等元数据
// @Tags image
// @Produce json
// @Success 200 {object} ImageMetaResp
// @Router /image/today/meta [get]
func GetTodayMeta(c *gin.Context) {
	img, err := image.GetTodayImage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Header("Cache-Control", "public, max-age=7200") // 2小时
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
	handleImageResponse(c, img, 0) // 禁用缓存
}

// GetRandomMeta 获取随机图片元数据
// @Summary 获取随机图片元数据
// @Description 随机获取一张已抓取图片的元数据
// @Tags image
// @Produce json
// @Success 200 {object} ImageMetaResp
// @Router /image/random/meta [get]
func GetRandomMeta(c *gin.Context) {
	img, err := image.GetRandomImage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
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
	handleImageResponse(c, img, 604800) // 7天
}

// GetByDateMeta 获取指定日期图片元数据
// @Summary 获取指定日期图片元数据
// @Description 根据日期获取图片元数据 (yyyy-mm-dd)
// @Tags image
// @Param date path string true "日期 (yyyy-mm-dd)"
// @Produce json
// @Success 200 {object} ImageMetaResp
// @Router /image/date/{date}/meta [get]
func GetByDateMeta(c *gin.Context) {
	date := c.Param("date")
	img, err := image.GetImageByDate(date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Header("Cache-Control", "public, max-age=604800") // 7天
	c.JSON(http.StatusOK, formatMeta(img))
}

// ListImages 获取图片列表
// @Summary 获取图片列表
// @Description 分页获取已抓取的图片元数据列表。支持分页(page, page_size)、限制数量(limit)和按月份过滤(month, 格式: YYYY-MM)。
// @Tags image
// @Param limit query int false "限制数量 (如果不使用分页)" default(30)
// @Param page query int false "页码 (从1开始)"
// @Param page_size query int false "每页数量"
// @Param month query string false "按月份过滤 (格式: YYYY-MM)"
// @Produce json
// @Success 200 {array} ImageMetaResp
// @Router /images [get]
func ListImages(c *gin.Context) {
	limitStr := c.Query("limit")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	month := c.Query("month")

	// 记录请求参数，便于排查过滤失效问题
	util.Logger.Debug("ListImages parameters",
		zap.String("month", month),
		zap.String("page", pageStr),
		zap.String("page_size", pageSizeStr),
		zap.String("limit", limitStr))

	var limit, offset int

	if pageStr != "" && pageSizeStr != "" {
		page, _ := strconv.Atoi(pageStr)
		pageSize, _ := strconv.Atoi(pageSizeStr)
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 30
		}
		limit = pageSize
		offset = (page - 1) * pageSize
	} else {
		if limitStr == "" {
			limit = 30
		} else {
			limit, _ = strconv.Atoi(limitStr)
		}
		offset = 0
	}

	images, err := image.GetImageList(limit, offset, month)
	if err != nil {
		util.Logger.Error("ListImages service call failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := []gin.H{}
	for _, img := range images {
		result = append(result, formatMeta(&img))
	}
	c.JSON(http.StatusOK, result)
}

func handleImageResponse(c *gin.Context, img *model.Image, maxAge int) {
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
			if maxAge > 0 {
				c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			} else {
				c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			}
			c.Redirect(http.StatusFound, selected.PublicURL)
		} else if img.URLBase != "" {
			// 兜底重定向到原始 Bing
			bingURL := fmt.Sprintf("https://www.bing.com%s_%s.jpg", img.URLBase, selected.Variant)
			if maxAge > 0 {
				c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			} else {
				c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			}
			c.Redirect(http.StatusFound, bingURL)
		} else {
			serveLocal(c, selected.StorageKey, img.Date, maxAge)
		}
	} else {
		serveLocal(c, selected.StorageKey, img.Date, maxAge)
	}
}

func serveLocal(c *gin.Context, key string, etag string, maxAge int) {
	if etag != "" {
		c.Header("ETag", fmt.Sprintf("\"%s\"", etag))
		if c.GetHeader("If-None-Match") == fmt.Sprintf("\"%s\"", etag) {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	}

	reader, contentType, err := storage.GlobalStorage.Get(context.Background(), key)
	if err != nil {
		util.Logger.Error("Failed to get image from storage", zap.String("key", key), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get image"})
		return
	}
	defer reader.Close()

	if contentType != "" {
		c.Header("Content-Type", contentType)
	}

	if maxAge > 0 {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	} else {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
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
		"date":          img.Date,
		"title":         img.Title,
		"copyright":     img.Copyright,
		"copyrightlink": img.CopyrightLink,
		"quiz":          img.Quiz,
		"startdate":     img.StartDate,
		"fullstartdate": img.FullStartDate,
		"hsh":           img.HSH,
		"variants":      variants,
	}
}

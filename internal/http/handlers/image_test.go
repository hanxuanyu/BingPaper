package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleImageResponseRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup config
	err := config.Init("")
	assert.NoError(t, err)
	config.GetConfig().API.Mode = "redirect"

	// Mock Image and Variant
	img := &model.Image{
		Date:    "2026-01-26",
		URLBase: "/th?id=OHR.TestImage",
		Variants: []model.ImageVariant{
			{
				Variant:    "UHD",
				Format:     "jpg",
				PublicURL:  "", // Empty for local storage simulation
				StorageKey: "2026-01-26/2026-01-26_UHD.jpg",
			},
		},
	}

	t.Run("Redirect mode with empty PublicURL should redirect to Bing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/v1/image/today?variant=UHD", nil)

		handleImageResponse(c, img, 0)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "bing.com")
		assert.Contains(t, w.Header().Get("Location"), "UHD")
	})

	t.Run("FormatMeta in redirect mode should return Bing URL if PublicURL is empty", func(t *testing.T) {
		config.GetConfig().API.Mode = "redirect"
		meta := formatMeta(img)

		variants := meta["variants"].([]gin.H)
		assert.Equal(t, 1, len(variants))
		assert.Contains(t, variants[0]["url"].(string), "bing.com")
		assert.Contains(t, variants[0]["url"].(string), "UHD")
	})

	t.Run("FormatMeta in local mode should return API URL", func(t *testing.T) {
		config.GetConfig().API.Mode = "local"
		config.GetConfig().Server.BaseURL = "http://myserver.com"
		meta := formatMeta(img)

		variants := meta["variants"].([]gin.H)
		assert.Equal(t, 1, len(variants))
		assert.Contains(t, variants[0]["url"].(string), "myserver.com")
		assert.Contains(t, variants[0]["url"].(string), "/api/v1/image/date/")
	})

	t.Run("FormatMetaSummary should only return the smallest variant", func(t *testing.T) {
		imgWithMultipleVariants := &model.Image{
			Date: "2026-01-26",
			Variants: []model.ImageVariant{
				{Variant: "UHD", Size: 1000, Format: "jpg"},
				{Variant: "640x480", Size: 200, Format: "jpg"},
				{Variant: "1920x1080", Size: 500, Format: "jpg"},
			},
		}
		meta := formatMetaSummary(imgWithMultipleVariants)
		variants := meta["variants"].([]gin.H)
		assert.Equal(t, 1, len(variants))
		assert.Equal(t, "640x480", variants[0]["variant"])
	})
}

func TestGetRegions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetRegions should respect pinned order", func(t *testing.T) {
		// Setup config with custom pinned regions
		config.Init("")
		config.GetConfig().Fetcher.Regions = []string{"en-US", "ja-JP"}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		GetRegions(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var regions []map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &regions)
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, len(regions), 2)
		assert.Equal(t, "en-US", regions[0]["value"])
		assert.Equal(t, "ja-JP", regions[1]["value"])
	})
}

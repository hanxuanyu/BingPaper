package handlers

import (
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

		handleImageResponse(c, img)

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
}

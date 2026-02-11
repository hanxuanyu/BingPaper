package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type LayoutResponse struct {
	Header string `json:"header"`
	Footer string `json:"footer"`
}

type UpdateLayoutRequest struct {
	Header *string `json:"header"`
	Footer *string `json:"footer"`
}

const (
	layoutDir  = "data/layout"
	headerFile = "header.txt"
	footerFile = "footer.txt"
)

// GetLayout 获取布局内容
func GetLayout(c *gin.Context) {
	header, _ := os.ReadFile(filepath.Join(layoutDir, headerFile))
	footer, _ := os.ReadFile(filepath.Join(layoutDir, footerFile))

	c.JSON(http.StatusOK, LayoutResponse{
		Header: string(header),
		Footer: string(footer),
	})
}

// UpdateLayout 更新布局内容
func UpdateLayout(c *gin.Context) {
	var req UpdateLayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.MkdirAll(layoutDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create layout directory"})
		return
	}

	if req.Header != nil {
		if err := os.WriteFile(filepath.Join(layoutDir, headerFile), []byte(*req.Header), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save header"})
			return
		}
	}

	if req.Footer != nil {
		if err := os.WriteFile(filepath.Join(layoutDir, footerFile), []byte(*req.Footer), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save footer"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layout updated successfully"})
}

package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"BingPaper/internal/config"
	"BingPaper/internal/service/fetcher"
	"BingPaper/internal/service/image"
	"BingPaper/internal/service/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 使用密码登录并获取临时 Token
// @Tags admin
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} model.Token
// @Failure 401 {object} map[string]string
// @Router /admin/login [post]
func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	t, err := token.Login(req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

// ListTokens 获取 Token 列表
// @Summary 获取 Token 列表
// @Description 获取所有已创建的 API Token 列表
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Token
// @Router /admin/tokens [get]
func ListTokens(c *gin.Context) {
	tokens, err := token.ListTokens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

type CreateTokenRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt string `json:"expires_at"` // optional
	ExpiresIn string `json:"expires_in"` // optional, e.g. 168h
}

// CreateToken 创建 Token
// @Summary 创建 Token
// @Description 创建一个新的 API Token
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateTokenRequest true "创建请求"
// @Success 200 {object} model.Token
// @Router /admin/tokens [post]
func CreateToken(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	expiresAt := time.Now().Add(config.GetTokenTTL())
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err == nil {
			expiresAt = t
		}
	} else if req.ExpiresIn != "" {
		d, err := time.ParseDuration(req.ExpiresIn)
		if err == nil {
			expiresAt = time.Now().Add(d)
		}
	}

	t, err := token.CreateToken(req.Name, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

type UpdateTokenRequest struct {
	Disabled bool `json:"disabled"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePassword 修改管理员密码
// @Summary 修改管理员密码
// @Description 验证旧密码并设置新密码，自动更新配置文件
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /admin/password [post]
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	cfg := config.GetConfig()
	// 验证旧密码
	err := bcrypt.CompareHashAndPassword([]byte(cfg.Admin.PasswordBcrypt), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}

	// 生成新密码 Hash
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// 更新配置
	cfg.Admin.PasswordBcrypt = string(hash)
	if err := config.SaveConfig(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "password updated successfully"})
}

// UpdateToken 更新 Token 状态
// @Summary 更新 Token 状态
// @Description 启用或禁用指定的 API Token
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Token ID"
// @Param request body UpdateTokenRequest true "更新请求"
// @Success 200 {object} map[string]string
// @Router /admin/tokens/{id} [patch]
func UpdateToken(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	var req UpdateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := token.UpdateToken(uint(id), req.Disabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteToken 删除 Token
// @Summary 删除 Token
// @Description 永久删除指定的 API Token
// @Tags admin
// @Security BearerAuth
// @Param id path int true "Token ID"
// @Success 200 {object} map[string]string
// @Router /admin/tokens/{id} [delete]
func DeleteToken(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	if err := token.DeleteToken(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetConfig 获取当前配置
// @Summary 获取当前配置
// @Description 获取服务的当前运行配置 (脱敏)
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} config.Config
// @Router /admin/config [get]
func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.GetConfig())
}

// UpdateConfig 更新配置
// @Summary 更新配置
// @Description 在线更新服务配置并保存
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body config.Config true "配置对象"
// @Success 200 {object} config.Config
// @Router /admin/config [put]
func UpdateConfig(c *gin.Context) {
	var cfg config.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := config.SaveConfig(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if c.Query("reload") == "true" {
		// 实际上 viper 会 watch config，但这里可以触发一些重新初始化逻辑
		// 这里暂不实现复杂的 reload
	}

	c.JSON(http.StatusOK, config.GetConfig())
}

type ManualFetchRequest struct {
	N int `json:"n"`
}

// ManualFetch 手动触发抓取
// @Summary 手动触发抓取
// @Description 立即启动抓取 Bing 任务
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ManualFetchRequest false "抓取天数"
// @Success 200 {object} map[string]string
// @Router /admin/fetch [post]
func ManualFetch(c *gin.Context) {
	var req ManualFetchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.N = config.BingFetchN
	}
	if req.N <= 0 {
		req.N = config.BingFetchN
	}

	f := fetcher.NewFetcher()
	go func() {
		f.Fetch(context.Background(), req.N)
	}()

	c.JSON(http.StatusOK, gin.H{"status": "task started"})
}

// ManualCleanup 手动触发清理
// @Summary 手动触发清理
// @Description 立即启动旧图片清理任务
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /admin/cleanup [post]
func ManualCleanup(c *gin.Context) {
	go func() {
		image.CleanupOldImages(context.Background())
	}()
	c.JSON(http.StatusOK, gin.H{"status": "task started"})
}

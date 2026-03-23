package handlers

import (
	"net/http"

	"BingPaper/internal/config"
	"BingPaper/internal/repo"

	"github.com/gin-gonic/gin"
)

type DatabaseConnectionRequest struct {
	Type string `json:"type" binding:"required"`
	DSN  string `json:"dsn" binding:"required"`
}

type DatabaseMigrationRequest struct {
	Type         string `json:"type" binding:"required"`
	DSN          string `json:"dsn" binding:"required"`
	UpdateConfig bool   `json:"update_config"`
}

type DatabaseStatusResponse struct {
	Active         config.DBConfig `json:"active"`
	Configured     config.DBConfig `json:"configured"`
	PendingRestart bool            `json:"pending_restart"`
}

type DatabaseMigrationResponse struct {
	Message         string              `json:"message"`
	Counts          repo.MigrationStats `json:"counts"`
	ConfigUpdated   bool                `json:"config_updated"`
	RestartRequired bool                `json:"restart_required"`
	Target          config.DBConfig     `json:"target"`
}

func sameDBConfig(a, b config.DBConfig) bool {
	return a.Type == b.Type && a.DSN == b.DSN
}

func GetDatabaseStatus(c *gin.Context) {
	active := repo.GetActiveDBConfig()
	configured := config.GetConfig().DB

	c.JSON(http.StatusOK, DatabaseStatusResponse{
		Active:         active,
		Configured:     configured,
		PendingRestart: !sameDBConfig(active, configured),
	})
}

func ValidateDatabaseConnection(c *gin.Context) {
	var req DatabaseConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": "invalid request"})
		return
	}

	target := config.DBConfig{Type: req.Type, DSN: req.DSN}
	if sameDBConfig(repo.GetActiveDBConfig(), target) {
		msg := "目标数据库不能与当前正在使用的数据库相同"
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "message": msg})
		return
	}

	if err := repo.ValidateDBConnection(config.GetConfig(), target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据库连接验证成功",
	})
}

func MigrateDatabase(c *gin.Context) {
	var req DatabaseMigrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": "invalid request"})
		return
	}

	target := config.DBConfig{Type: req.Type, DSN: req.DSN}
	if sameDBConfig(repo.GetActiveDBConfig(), target) {
		msg := "目标数据库不能与当前正在使用的数据库相同"
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "message": msg})
		return
	}

	currentCfg := config.GetConfig()
	stats, err := repo.MigrateDataToNewDB(repo.DB, currentCfg, target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": err.Error()})
		return
	}

	configUpdated := false
	if req.UpdateConfig {
		newCfg := *currentCfg
		newCfg.DB = target
		if err := config.SaveConfig(&newCfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": err.Error()})
			return
		}
		configUpdated = true
	}

	message := "数据库迁移成功，当前服务仍在使用旧库，请手动切换到新库并重启服务"
	if configUpdated {
		message = "数据库迁移成功，数据库配置已更新为新库，重启服务后生效"
	}

	c.JSON(http.StatusOK, DatabaseMigrationResponse{
		Message:         message,
		Counts:          stats,
		ConfigUpdated:   configUpdated,
		RestartRequired: true,
		Target:          target,
	})
}

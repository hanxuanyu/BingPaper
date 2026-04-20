package http

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

// HealthHandler serves the health check endpoint.
// @Summary 健康检查
// @Description 返回服务健康状态，用于部署探针与存活检查
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /healthz [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}

	_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

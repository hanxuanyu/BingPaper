package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSwaggerDocHandlerUsesRequestHostAndScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	req.Host = "example.com:9090"
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	SwaggerDocHandler(ctx)

	body := rec.Body.String()
	if !strings.Contains(body, `"host": "example.com:9090"`) {
		t.Fatalf("expected dynamic host in swagger doc, got %s", body)
	}
	if !strings.Contains(body, `"schemes": ["http"]`) {
		t.Fatalf("expected http scheme in swagger doc, got %s", body)
	}
}

func TestSwaggerDocHandlerPrefersForwardedHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	req.Host = "internal:8080"
	req.Header.Set("X-Forwarded-Host", "api.example.com")
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	SwaggerDocHandler(ctx)

	body := rec.Body.String()
	if !strings.Contains(body, `"host": "api.example.com"`) {
		t.Fatalf("expected forwarded host in swagger doc, got %s", body)
	}
	if !strings.Contains(body, `"schemes": ["https"]`) {
		t.Fatalf("expected https scheme in swagger doc, got %s", body)
	}
}

func TestSwaggerUIHandlerServesDynamicDocJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/swagger/*any", SwaggerUIHandler())

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	req.Host = "example.com:9090"
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"host": "example.com:9090"`) {
		t.Fatalf("expected dynamic host in swagger doc, got %s", body)
	}
}

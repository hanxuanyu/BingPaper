# Go 后端 CORS 配置示例

## 问题说明

开发环境中，前端运行在 `http://localhost:5173`，后端运行在 `http://localhost:8080`。
由于跨域限制，需要在后端配置 CORS 才能正常访问 API。

## 使用 Gin 框架

### 1. 安装 CORS 中间件

```bash
go get github.com/gin-contrib/cors
```

### 2. 配置 CORS

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "time"
)

func main() {
    r := gin.Default()
    
    // 开发环境：配置 CORS
    if gin.Mode() == gin.DebugMode {
        config := cors.Config{
            AllowOrigins:     []string{"http://localhost:5173"},
            AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
            AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
            ExposeHeaders:    []string{"Content-Length"},
            AllowCredentials: true,
            MaxAge:           12 * time.Hour,
        }
        r.Use(cors.New(config))
    }
    
    // API 路由
    api := r.Group("/api/v1")
    {
        // 图片相关
        api.GET("/images", getImages)
        api.GET("/image/today/meta", getTodayImageMeta)
        api.GET("/image/date/:date/meta", getImageMetaByDate)
        api.GET("/image/random/meta", getRandomImageMeta)
        
        // 管理员相关
        api.POST("/admin/login", adminLogin)
        // ... 其他路由
    }
    
    // 生产环境：静态文件服务
    if gin.Mode() == gin.ReleaseMode {
        r.Static("/assets", "./web/assets")
        r.StaticFile("/", "./web/index.html")
        
        // SPA fallback
        r.NoRoute(func(c *gin.Context) {
            c.File("./web/index.html")
        })
    }
    
    r.Run(":8080")
}
```

### 3. 更灵活的 CORS 配置

```go
// 根据环境变量动态配置
func setupCORS(r *gin.Engine) {
    allowOrigins := os.Getenv("ALLOW_ORIGINS")
    if allowOrigins == "" {
        allowOrigins = "http://localhost:5173"
    }
    
    origins := strings.Split(allowOrigins, ",")
    
    config := cors.Config{
        AllowOrigins:     origins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    
    r.Use(cors.New(config))
}

func main() {
    r := gin.Default()
    
    // 只在开发环境启用 CORS
    if gin.Mode() == gin.DebugMode {
        setupCORS(r)
    }
    
    // ... 其他配置
}
```

## 使用标准库

如果不使用 Gin 框架，可以手动实现 CORS：

```go
package main

import (
    "net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 设置 CORS 头
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Max-Age", "43200") // 12 hours
        
        // 处理预检请求
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    
    // 注册路由
    mux.HandleFunc("/api/v1/images", getImages)
    // ... 其他路由
    
    // 包装 CORS 中间件
    handler := corsMiddleware(mux)
    
    http.ListenAndServe(":8080", handler)
}
```

## 配置文件方式

可以通过配置文件管理 CORS 设置：

```yaml
# config.yaml
server:
  port: 8080
  mode: debug  # debug 或 release

cors:
  enabled: true
  allow_origins:
    - http://localhost:5173
    - http://127.0.0.1:5173
  allow_methods:
    - GET
    - POST
    - PUT
    - PATCH
    - DELETE
    - OPTIONS
  allow_headers:
    - Origin
    - Content-Type
    - Accept
    - Authorization
  allow_credentials: true
  max_age: 43200
```

```go
type CORSConfig struct {
    Enabled          bool     `yaml:"enabled"`
    AllowOrigins     []string `yaml:"allow_origins"`
    AllowMethods     []string `yaml:"allow_methods"`
    AllowHeaders     []string `yaml:"allow_headers"`
    AllowCredentials bool     `yaml:"allow_credentials"`
    MaxAge           int      `yaml:"max_age"`
}

func setupCORSFromConfig(r *gin.Engine, cfg *CORSConfig) {
    if !cfg.Enabled {
        return
    }
    
    config := cors.Config{
        AllowOrigins:     cfg.AllowOrigins,
        AllowMethods:     cfg.AllowMethods,
        AllowHeaders:     cfg.AllowHeaders,
        AllowCredentials: cfg.AllowCredentials,
        MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
    }
    
    r.Use(cors.New(config))
}
```

## 安全建议

1. **生产环境不要启用 CORS**：生产环境前后端在同一域名下，不需要 CORS
2. **限制允许的源**：不要使用 `*`，明确指定允许的域名
3. **使用环境变量**：允许的源应该通过环境变量配置，不要硬编码
4. **限制方法和头**：只允许必要的 HTTP 方法和请求头
5. **注意凭证**：如果使用 `AllowCredentials: true`，不能使用通配符源

## 测试 CORS 配置

可以使用 curl 测试 CORS：

```bash
# 测试预检请求
curl -X OPTIONS http://localhost:8080/api/v1/images \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: GET" \
  -v

# 测试实际请求
curl -X GET http://localhost:8080/api/v1/images \
  -H "Origin: http://localhost:5173" \
  -v
```

应该能看到响应头中包含：
```
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization
Access-Control-Allow-Credentials: true
```

## 常见问题

### Q: 为什么要区分开发和生产环境？
A: 生产环境前后端部署在同一域名下，不需要 CORS。只在开发环境需要。

### Q: 可以使用代理吗？
A: 可以在 Vite 中配置代理，但不推荐。直接配置 CORS 更接近生产环境，便于发现问题。

### Q: 如何处理认证？
A: 如果使用 Cookie 或需要发送凭证，必须设置 `AllowCredentials: true` 和明确的源。

### Q: 预检请求是什么？
A: 浏览器在某些跨域请求前会发送 OPTIONS 请求，询问服务器是否允许该跨域请求。
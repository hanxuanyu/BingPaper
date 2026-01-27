# BingPaper WebApp 构建说明

## 构建配置优化

本项目已优化构建配置，支持自定义后端路径和自动输出到上级目录的 `web` 文件夹。

## 环境配置

### 环境变量

项目支持通过环境变量配置后端 API 地址：

- **开发环境** (`.env.development`)：使用完整的后端服务器地址（直连，不使用代理）
- **生产环境** (`.env.production`)：使用相对路径 `/api/v1` 访问后端
- **默认配置** (`.env`)：通用配置

### 开发环境 vs 生产环境

#### 开发环境（npm run dev）
- 直接使用完整的后端 API 地址：`http://localhost:8080/api/v1`
- 不使用代理，前端直接请求后端服务
- 需要确保后端服务器支持 CORS 跨域请求
- 优点：配置简单，调试方便，可以清楚看到实际请求

#### 生产环境（npm run build）
- 使用相对路径：`/api/v1`
- 前后端部署在同一域名下，无跨域问题
- Go 服务器同时提供静态文件和 API 服务

### 自定义后端路径

可以通过修改环境变量 `VITE_API_BASE_URL` 来自定义后端 API 路径：

```bash
# 开发环境 (.env.development)
VITE_API_BASE_URL=http://localhost:8080/api/v1

# 或使用其他端口/域名
VITE_API_BASE_URL=http://192.168.1.100:8080/api/v1
VITE_API_BASE_URL=https://api.example.com/api/v1

# 生产环境 (.env.production)
VITE_API_BASE_URL=/api/v1

# 或自定义路径
VITE_API_BASE_URL=/custom/api/path
```

## 构建命令

### 开发环境

```bash
# 启动开发服务器（使用完整后端 URL）
npm run dev
```

开发环境会直接请求配置的后端服务器地址，无需代理。

**注意**：需要确保后端服务器配置了 CORS，允许来自 `http://localhost:5173` 的请求。

Go 后端 CORS 配置示例（使用 Gin）：
```go
import "github.com/gin-contrib/cors"

r := gin.Default()
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

### 生产环境构建

```bash
# 标准构建（生产模式）
npm run build

# 显式生产环境构建
npm run build:prod

# 开发模式构建（包含 sourcemap）
npm run build:dev
```

### 清理构建

```bash
# 清理输出目录
npm run clean
```

## 输出目录

构建产物会自动输出到项目上级目录的 `web` 文件夹：

```
go-project/
├── webapp/          # Vue 项目源码
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── web/            # 构建输出目录（自动生成）
│   ├── index.html
│   ├── assets/
│   └── ...
└── main.go         # Go 主程序
```

## Go 服务器配置

Go 服务器需要配置静态文件服务来访问 `web` 目录：

```go
// 示例：Gin 框架配置
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    
    // 开发环境：配置 CORS
    if gin.Mode() == gin.DebugMode {
        r.Use(cors.New(cors.Config{
            AllowOrigins:     []string{"http://localhost:5173"},
            AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
            AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
            AllowCredentials: true,
        }))
    }
    
    // API 路由
    api := r.Group("/api/v1")
    {
        // ... 你的 API 路由
    }
    
    // 静态文件服务（生产环境）
    r.Static("/assets", "./web/assets")
    r.StaticFile("/", "./web/index.html")
    
    // SPA fallback
    r.NoRoute(func(c *gin.Context) {
        c.File("./web/index.html")
    })
    
    r.Run(":8080")
}
```

## API 使用示例

项目提供了完整的 TypeScript API 客户端：

```typescript
import { bingPaperApi } from '@/lib/api-service'

// 获取今日图片元数据
const todayMeta = await bingPaperApi.getTodayImageMeta()

// 获取图片列表
const images = await bingPaperApi.getImages({ limit: 10 })

// 管理员登录
const token = await bingPaperApi.login({ password: 'admin123' })
bingPaperApi.setAuthToken(token.token)
```

## 项目结构

```
src/
├── lib/
│   ├── api-config.ts     # API 配置管理
│   ├── api-types.ts      # TypeScript 类型定义
│   ├── api-service.ts    # API 服务封装
│   ├── http-client.ts    # HTTP 客户端
│   └── utils.ts         # 工具函数
├── components/          # Vue 组件
├── views/              # 页面组件
│   ├── Home.vue         # 首页画廊
│   └── ImageView.vue    # 图片查看
├── composables/        # 组合式函数
│   └── useImages.ts    # 图片数据管理
├── router/             # 路由配置
│   └── index.ts
├── App.vue
└── main.ts
```

## 部署注意事项

1. **构建顺序**：确保在 Go 服务启动前完成前端构建
2. **路径配置**：Go 服务器的 API 路径应与前端配置的 `VITE_API_BASE_URL` 一致
3. **静态文件**：Go 服务器需要正确配置静态文件服务路径
4. **路由处理**：对于 SPA 应用，需要配置 fallback 到 `index.html`
5. **CORS 配置**：开发环境需要配置 CORS，生产环境不需要（同域）

## 故障排除

### API 请求失败（开发环境）

1. 检查环境变量配置是否正确
2. 确认 Go 服务器已启动在 `http://localhost:8080`
3. 检查 Go 服务器是否配置了 CORS
4. 在浏览器控制台查看具体错误信息

### CORS 错误

如果看到类似以下错误：
```
Access to fetch at 'http://localhost:8080/api/v1/...' from origin 'http://localhost:5173' 
has been blocked by CORS policy
```

解决方案：
1. 在 Go 服务器添加 CORS 中间件
2. 或者在 vite.config.ts 中添加代理配置（不推荐，因为会隐藏实际的请求路径）

### 构建失败

1. 清理 `node_modules` 并重新安装依赖
2. 检查 TypeScript 类型错误
3. 确认输出目录权限

### 静态资源加载失败

1. 检查 Go 服务器静态文件配置
2. 确认构建产物的路径结构
3. 检查 Vite 构建配置中的 `base` 和 `assetsDir`

## 开发工作流

### 典型的开发流程

1. **启动后端服务**
   ```bash
   cd ..  # 回到 Go 项目根目录
   go run main.go
   ```

2. **启动前端开发服务器**
   ```bash
   cd webapp
   npm run dev
   ```

3. **访问应用**
   ```
   http://localhost:5173/
   ```

4. **开发完成后构建**
   ```bash
   npm run build
   ```

5. **测试生产构建**
   ```bash
   # 停止开发服务器，启动 Go 服务器
   cd ..
   go run main.go
   # 访问 http://localhost:8080/
   ```

## 配置对比

| 配置项 | 开发环境 | 生产环境 |
|--------|----------|----------|
| API Base URL | `http://localhost:8080/api/v1` | `/api/v1` |
| 请求方式 | 直接请求 | 相对路径 |
| CORS | 需要 | 不需要 |
| 服务器 | 前后端分离 | 同一服务器 |
| 端口 | 前端 5173 + 后端 8080 | 8080 |

## 最佳实践

1. **开发环境**使用完整 URL，便于调试和查看实际请求
2. **生产环境**使用相对路径，简化部署
3. 保持 `.env.development` 和 `.env.production` 文件同步更新
4. 在 Go 服务器中使用环境变量区分开发/生产模式
5. 定期测试生产构建，确保配置正确
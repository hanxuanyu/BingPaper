# BingPaper

必应每日一图抓取、存储、多分辨率管理与公共 API 服务。

## 功能特性

- **自动抓取**：每日定时抓取 Bing 每日一图，支持 UHD 探测降级。
- **补抓能力**：支持手动或 API 触发抓取最近 N 天（默认 8 天）的图片。
- **多分辨率管理**：自动生成 UHD, 1920x1080, 1366x768 等分辨率，支持 JPG 格式。
- **灵活存储**：支持本地磁盘、S3 对象存储、WebDAV 存储。
- **数据库支持**：支持 SQLite, MySQL, PostgreSQL。
- **公共 API**：提供今日图片、随机图片、指定日期图片的纯图及元数据接口。
- **管理后台**：内置极简管理后台，支持 Token 管理、任务控制、配置查看。
- **单文件分发**：支持将前端页面嵌入二进制文件，实现无依赖运行。
- **行为模式**：支持 `local`（服务转发）和 `redirect`（302 跳转至公网 URL）两种模式。

## 快速启动

### 1. 配置

复制示例配置文件到 `data` 目录并根据需要修改：

```bash
mkdir -p data
cp config.example.yaml data/config.yaml
```

所有生成的数据（图片、数据库）及配置文件现在统一存放在 `./data` 目录下。

特别注意修改 `admin.password_bcrypt`（默认密码为 `admin123`）。

### 2. 运行

```bash
go run .
# 或者指定配置文件路径
./BingPaper -config /path/to/config.yaml
# 或者使用简写
./BingPaper -c /path/to/config.yaml
```

项目启动后会自动执行一次抓取任务，并根据 `cron.daily_spec` 设置定时任务。

### 3. 访问

- 管理后台：`http://localhost:8080/`
- 今日图片：`http://localhost:8080/api/v1/image/today`
- 今日元数据：`http://localhost:8080/api/v1/image/today/meta`
- API 文档 (Swagger)：`http://localhost:8080/swagger/index.html`

## API 文档 (v1)

### 公共接口 (无需 Token)

- `GET /api/v1/image/today`：返回今日图片
- `GET /api/v1/image/today/meta`：返回今日图片元数据
- `GET /api/v1/image/random`：返回随机图片
- `GET /api/v1/image/date/:yyyy-mm-dd`：返回指定日期图片
- **查询参数**：
  - `variant`：分辨率 (UHD, 1920x1080, 1366x768)，默认 `UHD`
  - `format`：格式 (jpg)，默认 `jpg`

### 管理接口 (需 Bearer Token)

- `POST /api/v1/admin/login`：登录获取 Token
- `GET /api/v1/admin/tokens`：Token 列表
- `POST /api/v1/admin/fetch`：手动触发抓取
- `POST /api/v1/admin/cleanup`：手动触发清理

## 存储模式区别

- **local 模式**：接口直接返回图片的二进制流，图片存储对外部不可见。
- **redirect 模式**：接口返回 302 重定向到图片的 `PublicURL`（通常在 S3 或 WebDAV 配置了 `public_url_prefix` 时使用）。

## 开发与构建

### 本地构建
您可以使用提供的脚本进行多平台构建：

```bash
# Unix/Linux/macOS
./scripts/build.sh [version]

# Windows (CMD)
.\scripts\build.bat [version]

# Windows (PowerShell)
.\scripts\build.ps1 [version]
```

编译后的打包文件将生成在 `output` 目录下。二进制文件已内置默认前端页面，即使不带 `web` 目录也能运行。如果需要自定义页面，可在配置中指定 `web.path`。

### 版本发布 (仅限维护者)
如果您是项目的维护者，可以通过以下步骤发布新版本：

1. 确保在 `master` 分支且代码已提交。
2. 运行标签脚本：`./scripts/tag.sh v1.0.0` (替换为实际版本号)。
3. 脚本会自动推送标签，触发 GitHub Actions 进行构建并发布 Release。

## 贡献指南

我们非常欢迎各种形式的贡献！如果您有任何想法或建议，请遵循以下流程：

1. **Fork** 本仓库到您的 GitHub 账号。
2. **Clone** 您 Fork 的仓库到本地。
3. 创建一个新的 **Feature 分支** (`git checkout -b feature/your-feature`)。
4. **提交** 您的修改 (`git commit -m 'Add some feature'`)。
5. **Push** 分支到 GitHub (`git push origin feature/your-feature`)。
6. 在本仓库提交一个 **Pull Request**。

### Docker 运行

#### 使用 Docker Hub 镜像 (推荐)
```bash
docker run -d \
  --name bingpaper \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  hxuanyu521/bingpaper:latest
```

#### 本地构建镜像
```bash
docker build -t bing-paper .
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data bing-paper
```

### Docker Compose (推荐)
使用项目根目录下的 `docker-compose.yml` 快速启动：

```bash
docker-compose up -d
```

你可以通过修改 `docker-compose.yml` 中的 `environment` 部分来覆盖默认配置，例如：
- `BINGPAPER_SERVER_PORT`: 服务端口
- `BINGPAPER_API_MODE`: API 模式 (`local` 或 `redirect`)
- `BINGPAPER_DB_TYPE`: 数据库类型 (`sqlite`, `mysql`, `postgres`)
- `BINGPAPER_STORAGE_TYPE`: 存储类型 (`local`, `s3`, `webdav`)
- `BINGPAPER_ADMIN_PASSWORD_BCRYPT`: 管理员密码的 Bcrypt 哈希值

## 许可证

MIT

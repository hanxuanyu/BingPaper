# BingPaper

必应每日一图抓取、存储、多分辨率管理与公共 API 服务。

## 功能特性

- **自动抓取**：每日定时抓取 Bing 每日一图，支持 UHD 探测降级。
- **补抓能力**：支持手动或 API 触发抓取最近 N 天（默认 8 天）的图片。
- **多分辨率管理**：自动生成 UHD, 1920x1080, 1366x768 等分辨率，支持 WebP 和 JPG 格式。
- **灵活存储**：支持本地磁盘、S3 对象存储、WebDAV 存储。
- **数据库支持**：支持 SQLite, MySQL, PostgreSQL。
- **公共 API**：提供今日图片、随机图片、指定日期图片的纯图及元数据接口。
- **管理后台**：内置极简管理后台，支持 Token 管理、任务控制、配置查看。
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
  - `format`：格式 (jpg, webp)，默认 `jpg`

### 管理接口 (需 Bearer Token)

- `POST /api/v1/admin/login`：登录获取 Token
- `GET /api/v1/admin/tokens`：Token 列表
- `POST /api/v1/admin/fetch`：手动触发抓取
- `POST /api/v1/admin/cleanup`：手动触发清理

## 存储模式区别

- **local 模式**：接口直接返回图片的二进制流，图片存储对外部不可见。
- **redirect 模式**：接口返回 302 重定向到图片的 `PublicURL`（通常在 S3 或 WebDAV 配置了 `public_url_prefix` 时使用）。

## 开发与构建

```bash
# 构建二进制
go build -o BingPaper .

# 构建 Docker 镜像
docker build -t bing-paper .
```

## 许可证

MIT

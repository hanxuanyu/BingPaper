# 配置指南

BingPaper 支持通过配置文件（YAML）和环境变量进行配置。

## 配置文件

程序启动时默认会查找当前目录下的 `config.yaml` 或 `data/config.yaml`。如果不存在，会自动在 `data/config.yaml` 创建一份带有默认值的配置文件。

你可以通过命令行参数 `-config` 或 `-c` 指定配置文件路径：

```bash
./BingPaper -c my_config.yaml
```

### 完整配置说明

以下是 `config.example.yaml` 的详细说明：

#### server (服务配置)
- `port`: 服务监听端口，默认 `8080`。
- `base_url`: 服务的基础 URL，用于生成某些绝对路径，默认为空。

#### log (日志配置)
- `level`: 业务日志级别，可选 `debug`, `info`, `warn`, `error`，默认 `info`。
- `filename`: 业务日志输出文件路径，默认 `data/logs/app.log`。
- `db_filename`: 数据库日志输出文件路径，默认 `data/logs/db.log`。
- `max_size`: 日志文件切割大小 (MB)，默认 `100`。
- `max_backups`: 保留旧日志文件个数，默认 `3`。
- `max_age`: 保留旧日志文件天数，默认 `7`。
- `compress`: 是否压缩旧日志文件，默认 `true`。
- `log_console`: 是否同时输出到控制台，默认 `true`。
- `show_db_log`: 是否在控制台输出数据库日志（SQL），默认 `false`。
- `db_log_level`: 数据库日志级别，可选 `debug`, `info`, `warn`, `error`, `silent`。`debug`/`info` 会记录所有 SQL。默认 `info`。

#### api (API 模式)
- `mode`: API 行为模式。
    - `local`: (默认) 接口直接返回图片的二进制流，适合图片存储对外部不可见的情况。
    - `redirect`: 接口返回 302 重定向到图片的 `PublicURL`，适合配合 S3 或 WebDAV 的公共访问。

#### cron (定时任务)
- `enabled`: 是否启用定时抓取，默认 `true`。
- `daily_spec`: Cron 表达式，定义每日抓取时间。默认 `"0 10 * * *"` (每日上午 10:00)。

#### retention (数据保留)
- `days`: 图片及元数据保留天数。超过此天数的数据可能会被清理任务处理。设置为 `0` 表示永久保留，不进行自动清理。默认 `0`。

#### db (数据库配置)
- `type`: 数据库类型，可选 `sqlite`, `mysql`, `postgres`。默认 `sqlite`。
- `dsn`: 数据库连接字符串。
    - SQLite: `data/bing_paper.db` (默认)
    - MySQL 示例: `user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
    - Postgres 示例: `host=localhost user=user password=pass dbname=db port=5432 sslmode=disable TimeZone=Asia/Shanghai`

**注意：** BingPaper 支持数据库配置的热更新。如果你在程序运行时修改了 `db.type` 或 `db.dsn`，程序会自动尝试将当前数据库中的所有数据（图片记录、变体信息、Token）迁移到新的数据库中。
- 在迁移开始前，程序会**清空**目标数据库中的相关表以防止数据冲突。
- 迁移过程在事务中执行，确保数据一致性。
- 迁移完成后，程序将无缝切换到新的数据库连接。

#### storage (存储配置)
- `type`: 存储类型，可选 `local`, `s3`, `webdav`。默认 `local`。
- **local (本地存储)**:
    - `root`: 图片存储根目录，默认 `data/picture`。
- **s3 (对象存储)**:
    - `endpoint`: S3 端点（如 `s3.amazonaws.com` 或 MinIO 地址）。
    - `region`: 区域（如 `us-east-1`）。
    - `bucket`: 桶名称。
    - `access_key`: 访问密钥 ID。
    - `secret_key`: 私有访问密钥。
    - `public_url_prefix`: 公网访问前缀，若为空则由 SDK 自动尝试生成。
    - `force_path_style`: 是否强制使用路径样式（MinIO 等通常需要设为 `true`）。
- **webdav (WebDAV 存储)**:
    - `url`: WebDAV 服务器地址。
    - `username`: 用户名。
    - `password`: 密码。
    - `public_url_prefix`: 公网访问前缀。

#### admin (管理配置)
- `password_bcrypt`: 管理员密码的 Bcrypt 哈希值。默认密码为 `admin123`，对应哈希 `$2a$10$fYHPeWHmwObephJvtlyH1O8DIgaLk5TINbi9BOezo2M8cSjmJchka`。
    - **强烈建议修改此项。**

#### token (认证配置)
- `default_ttl`: 管理后台登录 Token 的默认有效期，默认 `168h` (7天)。

#### feature (功能开关)
- `write_daily_files`: 是否在每日目录下写入原始文件（不仅是数据库记录），默认 `true`。

#### web (静态资源)
- `path`: 自定义管理后台前端文件的存放路径，默认 `web`。若指定路径不存在，将尝试使用内置的嵌入页面。

---

## 环境变量配置

所有的配置项都可以通过环境变量进行覆盖。环境变量前缀为 `BINGPAPER_`，层级之间使用下划线 `_` 分隔。

**常用示例：**

- `BINGPAPER_SERVER_PORT=9090`
- `BINGPAPER_DB_TYPE=mysql`
- `BINGPAPER_DB_DSN="user:pass@tcp(127.0.0.1:3306)/bingpaper"`
- `BINGPAPER_STORAGE_TYPE=s3`
- `BINGPAPER_STORAGE_S3_BUCKET=my-images`
- `BINGPAPER_ADMIN_PASSWORD_BCRYPT="$2a$10$..."`
- `HOST_PORT=8080` (仅限 Docker Compose 部署，控制宿主机映射到外部的端口)
- `BINGPAPER_SERVER_PORT=8080` (控制应用监听端口及容器内部端口)

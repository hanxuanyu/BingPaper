/**
 * BingPaper API TypeScript 接口定义
 * 基于 Swagger 文档自动生成
 */

// ===== 通用类型定义 =====

export interface ApiResponse<T = any> {
  data?: T
  message?: string
  success?: boolean
}

export interface PaginationParams {
  limit?: number
  offset?: number
}

// ===== Token 相关 =====

export interface Token {
  id: number
  name: string
  token: string
  disabled: boolean
  created_at: string
  updated_at: string
  expires_at?: string
}

export interface LoginRequest {
  password: string
}

export interface CreateTokenRequest {
  name: string
  expires_at?: string
  expires_in?: string
}

export interface UpdateTokenRequest {
  disabled?: boolean
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// ===== 配置相关 =====

export interface Config {
  Server: ServerConfig
  Log: LogConfig
  API: APIConfig
  Cron: CronConfig
  Retention: RetentionConfig
  DB: DBConfig
  Storage: StorageConfig
  Admin: AdminConfig
  Token: TokenConfig
  Feature: FeatureConfig
  Web: WebConfig
  Fetcher: FetcherConfig
}

export interface FetcherConfig {
  Regions: string[]
}

export interface AdminConfig {
  PasswordBcrypt: string
}

export interface APIConfig {
  Mode: string // 'local' | 'redirect'
  EnableMktFallback: boolean
  EnableOnDemandFetch: boolean
}

export interface CronConfig {
  Enabled: boolean
  DailySpec: string
}

export interface DBConfig {
  Type: string // 'sqlite' | 'mysql' | 'postgres'
  DSN: string
}

export interface FeatureConfig {
  WriteDailyFiles: boolean
}

export interface LogConfig {
  Level: string
  Filename: string
  DBFilename: string
  DBLogLevel: string
  LogConsole: boolean
  ShowDBLog: boolean
  MaxSize: number
  MaxAge: number
  MaxBackups: number
  Compress: boolean
}

export interface RetentionConfig {
  Days: number
}

export interface ServerConfig {
  Port: number
  BaseURL: string
}

export interface StorageConfig {
  Type: string // 'local' | 's3' | 'webdav'
  Local: LocalConfig
  S3: S3Config
  WebDAV: WebDAVConfig
}

export interface LocalConfig {
  Root: string
}

export interface S3Config {
  Endpoint: string
  AccessKey: string
  SecretKey: string
  Bucket: string
  Region: string
  ForcePathStyle: boolean
  PublicURLPrefix: string
}

export interface WebDAVConfig {
  URL: string
  Username: string
  Password: string
  PublicURLPrefix: string
}

export interface TokenConfig {
  DefaultTTL: string
}

export interface WebConfig {
  Path: string
}

// ===== 图片相关 =====

export interface ImageMeta {
  date?: string
  mkt?: string
  title?: string
  copyright?: string
  copyrightlink?: string    // 图片的详细版权链接（指向 Bing 搜索页面）
  quiz?: string             // 旧字段，保留向后兼容
  startdate?: string        // 图片的发布开始日期（格式：YYYYMMDD）
  fullstartdate?: string    // 图片的完整发布时间（格式：YYYYMMDDHHMM）
  hsh?: string              // 图片的唯一哈希值
  url?: string
  variant?: string
  format?: string
  variants?: ImageVariantResp[]  // 图片变体列表
  [key: string]: any
}

export interface ImageVariantResp {
  variant: string      // 分辨率变体 (UHD, 1920x1080, 等)
  format: string       // 格式 (jpg)
  url: string          // 访问 URL
  storage_key: string  // 存储键
  size: number         // 文件大小（字节）
}

export interface ImageListParams extends PaginationParams {
  page?: number        // 页码（从1开始）
  page_size?: number   // 每页数量
  month?: string       // 按月份过滤（格式：YYYY-MM）
  mkt?: string         // 地区编码
}

export interface Region {
  value: string
  label: string
}

export interface ManualFetchRequest {
  n?: number // 抓取天数
}

// ===== API 端点类型定义 =====

export type ImageVariant = 'UHD' | '1920x1080' | '1366x768' | '1280x720' | '1024x768' | '800x600' | '800x480' | '640x480' | '640x360' | '480x360' | '400x240' | '320x240'
export type ImageFormat = 'jpg'
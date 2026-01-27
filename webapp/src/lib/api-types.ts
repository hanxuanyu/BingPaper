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
  admin: AdminConfig
  api: APIConfig
  cron: CronConfig
  db: DBConfig
  feature: FeatureConfig
  log: LogConfig
  retention: RetentionConfig
  server: ServerConfig
  storage: StorageConfig
  token: TokenConfig
  web: WebConfig
}

export interface AdminConfig {
  passwordBcrypt: string
}

export interface APIConfig {
  mode: string // 'local' | 'redirect'
}

export interface CronConfig {
  enabled: boolean
  dailySpec: string
}

export interface DBConfig {
  type: string // 'sqlite' | 'mysql' | 'postgres'
  dsn: string
}

export interface FeatureConfig {
  writeDailyFiles: boolean
}

export interface LogConfig {
  level: string
  filename: string
  dbfilename: string
  dblogLevel: string
  logConsole: boolean
  showDBLog: boolean
  maxSize: number
  maxAge: number
  maxBackups: number
  compress: boolean
}

export interface RetentionConfig {
  days: number
}

export interface ServerConfig {
  port: number
  baseURL: string
}

export interface StorageConfig {
  type: string // 'local' | 's3' | 'webdav'
  local: LocalConfig
  s3: S3Config
  webDAV: WebDAVConfig
}

export interface LocalConfig {
  root: string
}

export interface S3Config {
  endpoint: string
  accessKey: string
  secretKey: string
  bucket: string
  region: string
  forcePathStyle: boolean
  publicURLPrefix: string
}

export interface WebDAVConfig {
  url: string
  username: string
  password: string
  publicURLPrefix: string
}

export interface TokenConfig {
  defaultTTL: string
}

export interface WebConfig {
  path: string
}

// ===== 图片相关 =====

export interface ImageMeta {
  date?: string
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
  [key: string]: any
}

export interface ImageListParams extends PaginationParams {
  page?: number        // 页码（从1开始）
  page_size?: number   // 每页数量
  month?: string       // 按月份过滤（格式：YYYY-MM）
}

export interface ManualFetchRequest {
  n?: number // 抓取天数
}

// ===== API 端点类型定义 =====

export type ImageVariant = 'UHD' | '1920x1080' | '1366x768' | '1280x720' | '1024x768' | '800x600' | '800x480' | '640x480' | '640x360' | '480x360' | '400x240' | '320x240'
export type ImageFormat = 'jpg'
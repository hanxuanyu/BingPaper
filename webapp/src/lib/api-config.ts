/**
 * API 配置管理
 * 用于管理后端 API 的基础配置
 */

// 获取环境变量中的 API 基础 URL
const getApiBaseUrl = (): string => {
  // 在构建时，Vite 会替换这个变量
  const baseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'
  
  // 确保以 / 开头但不以 / 结尾
  return baseUrl.replace(/\/$/, '')
}

// API 基础 URL
export const API_BASE_URL = getApiBaseUrl()

// API 配置
export const apiConfig = {
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 秒超时
  headers: {
    'Content-Type': 'application/json'
  }
}

/**
 * 构建完整的 API 端点 URL
 * @param endpoint - API 端点路径（如：'/images'）
 * @returns 完整的 API URL
 */
export const buildApiUrl = (endpoint: string): string => {
  // 确保端点以 / 开头
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`
  return `${API_BASE_URL}${normalizedEndpoint}`
}

/**
 * 标准化图片 URL
 * 当后端返回相对路径且配置了绝对 API 基础地址时，自动拼接完整域名
 */
export const normalizeImageUrl = (url: string | undefined): string => {
  if (!url) return ''
  if (url.startsWith('http')) return url

  // 处理相对路径问题：如果配置了绝对 API 基础地址，则拼接 Origin
  if (API_BASE_URL.startsWith('http')) {
    try {
      const origin = new URL(API_BASE_URL).origin
      return url.startsWith('/') ? origin + url : origin + '/' + url
    } catch (e) {
      // 解析失败则返回原样
    }
  }
  return url
}

/**
 * HTTP 状态码枚举
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500
} as const

export type HttpStatus = typeof HTTP_STATUS[keyof typeof HTTP_STATUS]
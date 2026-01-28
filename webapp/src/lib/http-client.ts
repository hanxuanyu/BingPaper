import { apiConfig, buildApiUrl } from './api-config'

/**
 * API 错误类
 */
export class ApiError extends Error {
  status: number
  response?: any

  constructor(message: string, status: number, response?: any) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.response = response
  }
}

/**
 * HTTP 请求选项
 */
export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  headers?: Record<string, string>
  body?: any
  timeout?: number
}

/**
 * 简单的 HTTP 客户端
 */
export class ApiClient {
  private defaultHeaders: Record<string, string>
  private defaultTimeout: number

  constructor() {
    this.defaultHeaders = { ...apiConfig.headers }
    this.defaultTimeout = apiConfig.timeout
  }

  /**
   * 设置默认请求头
   */
  setHeader(key: string, value: string) {
    this.defaultHeaders[key] = value
  }

  /**
   * 移除请求头
   */
  removeHeader(key: string) {
    delete this.defaultHeaders[key]
  }

  /**
   * 设置认证 Token
   */
  setAuthToken(token: string) {
    this.setHeader('Authorization', `Bearer ${token}`)
  }

  /**
   * 清除认证 Token
   */
  clearAuthToken() {
    this.removeHeader('Authorization')
  }

  /**
   * 发送 HTTP 请求
   */
  async request<T = any>(
    endpoint: string, 
    options: RequestOptions = {}
  ): Promise<T> {
    const url = buildApiUrl(endpoint)
    const {
      method = 'GET',
      headers = {},
      body,
      timeout = this.defaultTimeout
    } = options

    // 合并请求头
    const requestHeaders = {
      ...this.defaultHeaders,
      ...headers
    }

    // 构建请求配置
    const requestConfig: RequestInit = {
      method,
      headers: requestHeaders,
      signal: AbortSignal.timeout(timeout)
    }

    // 处理请求体
    if (body && method !== 'GET') {
      if (typeof body === 'object') {
        requestConfig.body = JSON.stringify(body)
      } else {
        requestConfig.body = body
      }
    }

    try {
      const response = await fetch(url, requestConfig)
      
      // 检查响应状态
      if (!response.ok) {
        const errorData = await this.parseResponse(response)
        const apiError = new ApiError(
          errorData?.message || `HTTP ${response.status}: ${response.statusText}`,
          response.status,
          errorData
        )
        
        // 401 未授权错误，自动跳转到登录页
        if (response.status === 401) {
          this.handle401Error()
        }
        
        throw apiError
      }

      return await this.parseResponse(response)
    } catch (error) {
      if (error instanceof ApiError) {
        throw error
      }
      
      // 处理网络错误或其他异常
      if (error instanceof Error) {
        throw new ApiError(error.message, 0)
      }
      
      throw new ApiError('Unknown error occurred', 0)
    }
  }

  /**
   * 处理 401 错误
   */
  private handle401Error() {
    // 清除本地存储的 token
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_token_expires')
    this.clearAuthToken()
    
    // 只有在管理页面时才跳转到登录页
    if (typeof window !== 'undefined' && window.location.pathname.startsWith('/admin')) {
      // 避免重复跳转
      if (!window.location.pathname.includes('/admin/login')) {
        window.location.href = '/admin/login'
      }
    }
  }

  /**
   * 解析响应数据
   */
  private async parseResponse(response: Response): Promise<any> {
    const contentType = response.headers.get('content-type')
    
    if (contentType?.includes('application/json')) {
      return await response.json()
    }
    
    if (contentType?.includes('text/')) {
      return await response.text()
    }
    
    // 对于图片等二进制数据
    return response
  }

  /**
   * GET 请求
   */
  async get<T = any>(endpoint: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET', headers })
  }

  /**
   * POST 请求
   */
  async post<T = any>(endpoint: string, body?: any, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: 'POST', body, headers })
  }

  /**
   * PUT 请求
   */
  async put<T = any>(endpoint: string, body?: any, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: 'PUT', body, headers })
  }

  /**
   * DELETE 请求
   */
  async delete<T = any>(endpoint: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE', headers })
  }

  /**
   * PATCH 请求
   */
  async patch<T = any>(endpoint: string, body?: any, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(endpoint, { method: 'PATCH', body, headers })
  }
}

// 导出默认实例
export const apiClient = new ApiClient()
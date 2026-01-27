import { apiClient } from './http-client'
import { apiConfig } from './api-config'
import type {
  Token,
  LoginRequest,
  CreateTokenRequest,
  UpdateTokenRequest,
  ChangePasswordRequest,
  Config,
  ImageMeta,
  ImageListParams,
  ManualFetchRequest,
  ImageVariant,
  ImageFormat
} from './api-types'

/**
 * BingPaper API 服务类
 */
export class BingPaperApiService {
  
  // ===== 认证相关 =====
  
  /**
   * 管理员登录
   */
  async login(request: LoginRequest): Promise<Token> {
    return apiClient.post<Token>('/admin/login', request)
  }

  /**
   * 修改管理员密码
   */
  async changePassword(request: ChangePasswordRequest): Promise<{ message: string }> {
    return apiClient.post('/admin/password', request)
  }

  // ===== Token 管理 =====

  /**
   * 获取 Token 列表
   */
  async getTokens(): Promise<Token[]> {
    return apiClient.get<Token[]>('/admin/tokens')
  }

  /**
   * 创建 Token
   */
  async createToken(request: CreateTokenRequest): Promise<Token> {
    return apiClient.post<Token>('/admin/tokens', request)
  }

  /**
   * 更新 Token 状态
   */
  async updateToken(id: number, request: UpdateTokenRequest): Promise<{ message: string }> {
    return apiClient.patch(`/admin/tokens/${id}`, request)
  }

  /**
   * 删除 Token
   */
  async deleteToken(id: number): Promise<{ message: string }> {
    return apiClient.delete(`/admin/tokens/${id}`)
  }

  // ===== 配置管理 =====

  /**
   * 获取当前配置
   */
  async getConfig(): Promise<Config> {
    return apiClient.get<Config>('/admin/config')
  }

  /**
   * 更新配置
   */
  async updateConfig(config: Config): Promise<Config> {
    return apiClient.put<Config>('/admin/config', config)
  }

  // ===== 系统管理 =====

  /**
   * 手动触发抓取
   */
  async manualFetch(request?: ManualFetchRequest): Promise<{ message: string }> {
    return apiClient.post('/admin/fetch', request)
  }

  /**
   * 手动触发清理
   */
  async manualCleanup(): Promise<{ message: string }> {
    return apiClient.post('/admin/cleanup')
  }

  // ===== 图片相关 =====

  /**
   * 获取图片列表
   */
  async getImages(params?: ImageListParams): Promise<ImageMeta[]> {
    const searchParams = new URLSearchParams()
    if (params?.limit) searchParams.set('limit', params.limit.toString())
    if (params?.offset) searchParams.set('offset', params.offset.toString())
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.page_size) searchParams.set('page_size', params.page_size.toString())
    if (params?.month) searchParams.set('month', params.month)
    
    const queryString = searchParams.toString()
    const endpoint = queryString ? `/images?${queryString}` : '/images'
    
    return apiClient.get<ImageMeta[]>(endpoint)
  }

  /**
   * 获取今日图片元数据
   */
  async getTodayImageMeta(): Promise<ImageMeta> {
    return apiClient.get<ImageMeta>('/image/today/meta')
  }

  /**
   * 获取指定日期图片元数据
   */
  async getImageMetaByDate(date: string): Promise<ImageMeta> {
    return apiClient.get<ImageMeta>(`/image/date/${date}/meta`)
  }

  /**
   * 获取随机图片元数据
   */
  async getRandomImageMeta(): Promise<ImageMeta> {
    return apiClient.get<ImageMeta>('/image/random/meta')
  }

  /**
   * 构建图片 URL
   */
  getTodayImageUrl(variant: ImageVariant = 'UHD', format: ImageFormat = 'jpg'): string {
    const params = new URLSearchParams({ variant, format })
    return `${apiConfig.baseURL}/image/today?${params.toString()}`
  }

  /**
   * 构建指定日期图片 URL
   */
  getImageUrlByDate(date: string, variant: ImageVariant = 'UHD', format: ImageFormat = 'jpg'): string {
    const params = new URLSearchParams({ variant, format })
    return `${apiConfig.baseURL}/image/date/${date}?${params.toString()}`
  }

  /**
   * 构建随机图片 URL
   */
  getRandomImageUrl(variant: ImageVariant = 'UHD', format: ImageFormat = 'jpg'): string {
    const params = new URLSearchParams({ variant, format })
    return `${apiConfig.baseURL}/image/random?${params.toString()}`
  }

  // ===== 认证状态管理 =====

  /**
   * 设置认证 Token
   */
  setAuthToken(token: string): void {
    apiClient.setAuthToken(token)
  }

  /**
   * 清除认证 Token
   */
  clearAuthToken(): void {
    apiClient.clearAuthToken()
  }
}

// 导出默认实例
export const bingPaperApi = new BingPaperApiService()

// 导出便捷方法
export const {
  login,
  changePassword,
  getTokens,
  createToken,
  updateToken,
  deleteToken,
  getConfig,
  updateConfig,
  manualFetch,
  manualCleanup,
  getImages,
  getTodayImageMeta,
  getImageMetaByDate,
  getRandomImageMeta,
  getTodayImageUrl,
  getImageUrlByDate,
  getRandomImageUrl,
  setAuthToken,
  clearAuthToken
} = bingPaperApi

const MKT_STORAGE_KEY = 'bing_paper_selected_mkt'
const DEFAULT_MKT = 'zh-CN'

/**
 * 默认地区列表 (兜底用)
 */
export const DEFAULT_REGIONS = [
  { value: 'zh-CN', label: '中国' },
  { value: 'en-US', label: '美国' },
  { value: 'ja-JP', label: '日本' },
  { value: 'en-AU', label: '澳大利亚' },
  { value: 'en-GB', label: '英国' },
  { value: 'de-DE', label: '德国' },
  { value: 'en-NZ', label: '新西兰' },
  { value: 'en-CA', label: '加拿大' },
  { value: 'fr-FR', label: '法国' },
  { value: 'it-IT', label: '意大利' },
  { value: 'es-ES', label: '西班牙' },
  { value: 'pt-BR', label: '巴西' },
  { value: 'ko-KR', label: '韩国' },
  { value: 'en-IN', label: '印度' },
  { value: 'ru-RU', label: '俄罗斯' },
  { value: 'zh-HK', label: '中国香港' },
  { value: 'zh-TW', label: '中国台湾' },
]

/**
 * 支持的地区列表 (优先使用后端提供的)
 */
export let SUPPORTED_REGIONS = [...DEFAULT_REGIONS]

/**
 * 更新支持的地区列表
 */
export function setSupportedRegions(regions: typeof DEFAULT_REGIONS): void {
  SUPPORTED_REGIONS = regions
}

/**
 * 获取浏览器首选地区
 */
export function getBrowserMkt(): string {
  const lang = navigator.language || (navigator as any).userLanguage
  if (!lang) return DEFAULT_MKT

  // 尝试精确匹配
  const exactMatch = SUPPORTED_REGIONS.find(r => r.value.toLowerCase() === lang.toLowerCase())
  if (exactMatch) return exactMatch.value

  // 尝试模糊匹配 (前两个字符，如 en-GB 匹配 en-US)
  const prefix = lang.split('-')[0].toLowerCase()
  const prefixMatch = SUPPORTED_REGIONS.find(r => r.value.split('-')[0].toLowerCase() === prefix)
  if (prefixMatch) return prefixMatch.value

  return DEFAULT_MKT
}

/**
 * 获取当前选择的地区 (优先从 localStorage 获取，其次从浏览器获取)
 */
export function getDefaultMkt(): string {
  const saved = localStorage.getItem(MKT_STORAGE_KEY)
  if (saved && SUPPORTED_REGIONS.some(r => r.value === saved)) {
    return saved
  }
  return getBrowserMkt()
}

/**
 * 保存选择的地区
 */
export function setSavedMkt(mkt: string): void {
  localStorage.setItem(MKT_STORAGE_KEY, mkt)
}

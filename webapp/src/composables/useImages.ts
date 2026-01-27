import { ref, onMounted } from 'vue'
import { bingPaperApi } from '@/lib/api-service'
import type { ImageMeta } from '@/lib/api-types'

/**
 * 获取今日图片
 */
export function useTodayImage() {
  const image = ref<ImageMeta | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const fetchImage = async () => {
    loading.value = true
    error.value = null
    try {
      image.value = await bingPaperApi.getTodayImageMeta()
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch today image:', e)
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchImage()
  })

  return {
    image,
    loading,
    error,
    refetch: fetchImage
  }
}

/**
 * 获取图片列表（支持分页）
 */
export function useImageList(initialLimit = 30) {
  const images = ref<ImageMeta[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  const hasMore = ref(true)

  const fetchImages = async (limit = initialLimit) => {
    if (loading.value) return
    
    loading.value = true
    error.value = null
    try {
      const newImages = await bingPaperApi.getImages({ limit })
      
      if (newImages.length < limit) {
        hasMore.value = false
      }
      
      images.value = [...images.value, ...newImages]
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch images:', e)
    } finally {
      loading.value = false
    }
  }

  const loadMore = () => {
    if (!loading.value && hasMore.value) {
      fetchImages()
    }
  }

  onMounted(() => {
    fetchImages()
  })

  return {
    images,
    loading,
    error,
    hasMore,
    loadMore,
    refetch: () => {
      images.value = []
      hasMore.value = true
      fetchImages()
    }
  }
}

/**
 * 获取指定日期的图片
 */
export function useImageByDate(date: string) {
  const image = ref<ImageMeta | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const fetchImage = async () => {
    loading.value = true
    error.value = null
    try {
      image.value = await bingPaperApi.getImageMetaByDate(date)
    } catch (e) {
      error.value = e as Error
      console.error(`Failed to fetch image for date ${date}:`, e)
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchImage()
  })

  return {
    image,
    loading,
    error,
    refetch: fetchImage
  }
}

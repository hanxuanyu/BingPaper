import { ref, onMounted, watch } from 'vue'
import type { Ref } from 'vue'
import { bingPaperApi } from '@/lib/api-service'
import type { ImageMeta } from '@/lib/api-types'
import { getDefaultMkt } from '@/lib/mkt-utils'

/**
 * 获取今日图片
 */
export function useTodayImage(mkt?: string) {
  const image = ref<ImageMeta | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const fetchImage = async () => {
    loading.value = true
    error.value = null
    try {
      image.value = await bingPaperApi.getTodayImageMeta(mkt || getDefaultMkt())
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
 * 获取图片列表（支持分页和月份筛选）
 */
export function useImageList(pageSize = 30) {
  const images = ref<ImageMeta[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  const hasMore = ref(true)
  const currentPage = ref(1)
  const currentMonth = ref<string | undefined>(undefined)
  const currentMkt = ref<string | undefined>(getDefaultMkt())

  const fetchImages = async (page = 1, month?: string, mkt?: string) => {
    if (loading.value) return
    
    loading.value = true
    error.value = null
    try {
      const params: any = {
        page,
        page_size: pageSize,
        mkt: mkt || currentMkt.value || getDefaultMkt()
      }
      if (month) {
        params.month = month
      }
      
      const newImages = await bingPaperApi.getImages(params)
      
      if (page === 1) {
        // 首次加载或重新筛选
        images.value = newImages
      } else {
        // 加载更多
        images.value = [...images.value, ...newImages]
      }
      
      // 判断是否还有更多数据
      hasMore.value = newImages.length === pageSize
      currentPage.value = page
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch images:', e)
    } finally {
      loading.value = false
    }
  }

  const loadMore = () => {
    if (!loading.value && hasMore.value) {
      fetchImages(currentPage.value + 1, currentMonth.value, currentMkt.value)
    }
  }

  const filterByMonth = (month?: string) => {
    currentMonth.value = month
    currentPage.value = 1
    hasMore.value = true
    fetchImages(1, month, currentMkt.value)
  }

  const filterByMkt = (mkt?: string) => {
    currentMkt.value = mkt
    currentPage.value = 1
    hasMore.value = true
    fetchImages(1, currentMonth.value, mkt)
  }

  onMounted(() => {
    fetchImages(1)
  })

  return {
    images,
    loading,
    error,
    hasMore,
    loadMore,
    filterByMonth,
    filterByMkt,
    refetch: () => {
      currentPage.value = 1
      hasMore.value = true
      fetchImages(1, currentMonth.value, currentMkt.value)
    }
  }
}

/**
 * 获取指定日期的图片
 */
export function useImageByDate(dateRef: Ref<string>, mktRef?: Ref<string | undefined>) {
  const image = ref<ImageMeta | null>(null)
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const fetchImage = async () => {
    loading.value = true
    error.value = null
    try {
      image.value = await bingPaperApi.getImageMetaByDate(dateRef.value, mktRef?.value || getDefaultMkt())
    } catch (e) {
      error.value = e as Error
      console.error(`Failed to fetch image for date ${dateRef.value}:`, e)
    } finally {
      loading.value = false
    }
  }

  // 监听日期和地区变化，自动重新获取
  if (mktRef) {
    watch([dateRef, mktRef], () => {
      fetchImage()
    }, { immediate: true })
  } else {
    watch(dateRef, () => {
      fetchImage()
    }, { immediate: true })
  }

  return {
    image,
    loading,
    error,
    refetch: fetchImage
  }
}

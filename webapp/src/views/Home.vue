<template>
  <div class="min-h-screen bg-gradient-to-b from-gray-900 via-gray-800 to-gray-900">
    <!-- Hero Section - 今日图片 -->
    <section class="relative h-screen w-full overflow-hidden">
      <div v-if="todayLoading" class="absolute inset-0 flex items-center justify-center">
        <div class="w-12 h-12 border-4 border-white/20 border-t-white rounded-full animate-spin"></div>
      </div>
      
      <div v-else-if="latestImage" class="relative h-full w-full group">
        <!-- 背景图片 -->
        <div class="absolute inset-0">
          <img 
            :src="getLatestImageUrl()" 
            :alt="latestImage.title || 'Latest Bing Image'"
            class="w-full h-full object-cover"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>
        </div>

        <!-- 更新提示（仅在非今日图片时显示） - 右上角简约徽章 -->
        <div v-if="!isToday" class="absolute top-4 right-4 md:top-8 md:right-8 z-20">
          <div class="flex items-center gap-1.5 px-3 py-1.5 bg-black/30 backdrop-blur-md rounded-full border border-white/10 text-white/70 text-xs">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span>下次更新 {{ nextUpdateTime }}</span>
          </div>
        </div>

        <!-- 内容叠加层 -->
        <div class="relative h-full flex flex-col justify-end p-8 md:p-16 z-10">
          <div class="max-w-4xl space-y-4 transform transition-transform duration-500 group-hover:translate-y-[-10px]">
            <div class="inline-block px-4 py-2 bg-white/10 backdrop-blur-md rounded-full text-white/90 text-sm font-medium">
              {{ isToday ? '今日精选' : '最新图片' }} · {{ formatDate(latestImage.date) }}
            </div>
            
            <h1 class="text-4xl md:text-6xl font-bold text-white leading-tight drop-shadow-2xl">
              {{ latestImage.title || '必应每日一图' }}
            </h1>
            
            <p v-if="latestImage.copyright" class="text-lg md:text-xl text-white/80 max-w-2xl">
              {{ latestImage.copyright }}
            </p>

            <div class="flex gap-4 pt-4">
              <button 
                @click="viewImage(latestImage.date!)"
                class="px-6 py-3 bg-white text-gray-900 rounded-lg font-semibold hover:bg-white/90 transition-all transform hover:scale-105 shadow-xl"
              >
                查看大图
              </button>
              <button 
                v-if="latestImage.copyrightlink"
                @click="openCopyrightLink(latestImage.copyrightlink)"
                class="px-6 py-3 bg-white/10 backdrop-blur-md text-white rounded-lg font-semibold hover:bg-white/20 transition-all border border-white/30"
              >
                了解更多
              </button>
            </div>
          </div>
        </div>

        <!-- 滚动提示 -->
        <div class="absolute bottom-8 left-1/2 transform -translate-x-1/2 animate-bounce">
          <svg class="w-6 h-6 text-white/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
          </svg>
        </div>
      </div>
    </section>

    <!-- Gallery Section - 历史图片 -->
    <section class="py-16 px-4 md:px-8 lg:px-16">
      <div class="max-w-7xl mx-auto">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
          <h2 class="text-3xl md:text-4xl font-bold text-white">
            历史精选
          </h2>
          
          <!-- 筛选器 -->
          <div class="flex flex-wrap items-center gap-3">
            <!-- 年份选择 -->
            <Select v-model="selectedYear" @update:model-value="onYearChange">
              <SelectTrigger class="w-[180px] bg-white/10 backdrop-blur-md text-white border-white/20 hover:bg-white/15 hover:border-white/30 focus:ring-white/50 shadow-lg">
                <SelectValue placeholder="选择年份" />
              </SelectTrigger>
              <SelectContent class="bg-gray-900/95 backdrop-blur-xl border-white/20 text-white">
                <SelectItem 
                  v-for="year in availableYears" 
                  :key="year" 
                  :value="String(year)"
                  class="focus:bg-white/10 focus:text-white cursor-pointer"
                >
                  {{ year }} 年
                </SelectItem>
              </SelectContent>
            </Select>
            
            <!-- 月份选择 -->
            <Select v-model="selectedMonth" @update:model-value="onFilterChange" :disabled="!selectedYear">
              <SelectTrigger 
                class="w-[180px] bg-white/10 backdrop-blur-md text-white border-white/20 hover:bg-white/15 hover:border-white/30 focus:ring-white/50 shadow-lg disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <SelectValue placeholder="选择月份" />
              </SelectTrigger>
              <SelectContent class="bg-gray-900/95 backdrop-blur-xl border-white/20 text-white">
                <SelectItem 
                  v-for="month in 12" 
                  :key="month" 
                  :value="String(month)"
                  class="focus:bg-white/10 focus:text-white cursor-pointer"
                >
                  {{ month }} 月
                </SelectItem>
              </SelectContent>
            </Select>
            
            <!-- 重置按钮 -->
            <button 
              v-if="selectedYear && selectedMonth"
              @click="resetFilters"
              class="px-4 py-2.5 bg-gradient-to-r from-red-500/20 to-pink-500/20 backdrop-blur-md text-white rounded-lg hover:from-red-500/30 hover:to-pink-500/30 border border-red-400/30 hover:border-red-400/50 transition-all flex items-center gap-2 text-sm font-medium shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
              清除筛选
            </button>
          </div>
        </div>

        <!-- 筛选结果提示 -->
        <div v-if="selectedYear && selectedMonth" class="mb-6 flex items-center gap-2 text-white/60 text-sm">
          <span>当前显示：</span>
          <span class="text-white font-medium">
            {{ selectedYear }} 年 {{ selectedMonth }} 月
          </span>
          <span>的图片（共 {{ images.length }} 张）</span>
        </div>

        <!-- 图片网格 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div 
            v-for="(image, index) in images" 
            :key="image.date || index"
            :ref="el => setImageRef(el, index)"
            class="group relative aspect-video rounded-xl overflow-hidden cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-2xl"
            @click="viewImage(image.date!)"
          >
            <!-- 图片（懒加载） -->
            <div v-if="!imageVisibility[index]" class="w-full h-full bg-white/5 flex items-center justify-center">
              <div class="w-8 h-8 border-2 border-white/20 border-t-white/60 rounded-full animate-spin"></div>
            </div>
            <img 
              v-else
              :src="getImageUrl(image.date!)" 
              :alt="image.title || 'Bing Image'"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
              loading="lazy"
            />
            
            <!-- 悬浮信息层 -->
            <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/50 to-transparent md:opacity-0 md:group-hover:opacity-100 transition-opacity duration-300">
              <div class="absolute bottom-0 left-0 right-0 p-4 md:p-6 transform md:translate-y-4 md:group-hover:translate-y-0 transition-transform duration-300">
                <div class="text-xs text-white/70 mb-1">
                  {{ formatDate(image.date) }}
                </div>
                <h3 class="text-base md:text-lg font-semibold text-white mb-1 md:mb-2 line-clamp-2">
                  {{ image.title || '未命名' }}
                </h3>
                <p v-if="image.copyright" class="text-xs md:text-sm text-white/80 line-clamp-2">
                  {{ image.copyright }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- 加载更多 -->
        <div class="mt-12 text-center">
          <div v-if="loading" class="inline-flex items-center gap-2 text-white/60">
            <div class="w-5 h-5 border-2 border-white/20 border-t-white/60 rounded-full animate-spin"></div>
            <span>加载中...</span>
          </div>
          
          <div 
            v-else-if="hasMore"
            ref="loadMoreTrigger"
            class="inline-block"
          >
            <button 
              @click="loadMore"
              class="px-8 py-3 bg-white/10 backdrop-blur-md text-white rounded-lg font-semibold hover:bg-white/20 transition-all border border-white/30"
            >
              加载更多
            </button>
          </div>

          <p v-else class="text-white/40">
            已加载全部图片
          </p>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="py-12 px-4 border-t border-white/10">
      <div class="max-w-7xl mx-auto">
        <div class="flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="text-white/60 text-sm">
            <p>数据来源于必应每日一图 API</p>
            <p class="mt-1 text-white/40">BingPaper © 2026</p>
          </div>
          
          <div class="flex gap-6">
            <router-link 
              to="/api-docs"
              class="text-white/60 hover:text-white transition-colors text-sm flex items-center gap-2 group"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
              </svg>
              <span>API 文档</span>
              <svg class="w-4 h-4 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </router-link>
            
            <a 
              href="https://github.com/hanxuanyu/BingPaper" 
              target="_blank"
              class="text-white/60 hover:text-white transition-colors text-sm flex items-center gap-2"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd"></path>
              </svg>
              <span>GitHub</span>
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useImageList } from '@/composables/useImages'
import { bingPaperApi } from '@/lib/api-service'
import { useRouter } from 'vue-router'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const router = useRouter()

// 顶部最新图片（独立加载，不受筛选影响）
const latestImage = ref<any>(null)
const todayLoading = ref(false)

// 历史图片列表（使用服务端分页和筛选，每页15张）
const { images, loading, hasMore, loadMore, filterByMonth } = useImageList(15)

// 加载顶部最新图片
const loadLatestImage = async () => {
  todayLoading.value = true
  try {
    const params = { page: 1, page_size: 1 }
    const result = await bingPaperApi.getImages(params)
    if (result.length > 0) {
      latestImage.value = result[0]
    }
  } catch (error) {
    console.error('Failed to load latest image:', error)
  } finally {
    todayLoading.value = false
  }
}

// 初始化加载顶部图片
onMounted(() => {
  loadLatestImage()
})

// 判断最新图片是否为今天的图片
const isToday = computed(() => {
  if (!latestImage.value?.date) return false
  const imageDate = new Date(latestImage.value.date).toDateString()
  const today = new Date().toDateString()
  return imageDate === today
})

// 计算下次更新时间提示
const nextUpdateTime = computed(() => {
  const now = new Date()
  const hours = now.getHours()
  
  // 更新时间点：8:20, 12:20, 16:20, 20:20, 0:20, 4:20
  const updateHours = [0, 4, 8, 12, 16, 20]
  const updateMinute = 20
  
  // 找到下一个更新时间点
  for (const hour of updateHours) {
    if (hours < hour || (hours === hour && now.getMinutes() < updateMinute)) {
      return `${String(hour).padStart(2, '0')}:${String(updateMinute).padStart(2, '0')}`
    }
  }
  
  // 如果今天没有下一个更新点，返回明天的第一个更新点
  return `次日 00:20`
})

// 筛选相关状态
const selectedYear = ref('')
const selectedMonth = ref('')

// 懒加载相关
const imageRefs = ref<(HTMLElement | null)[]>([])
const imageVisibility = ref<boolean[]>([])
let observer: IntersectionObserver | null = null

// 无限滚动加载
const loadMoreTrigger = ref<HTMLElement | null>(null)
let loadMoreObserver: IntersectionObserver | null = null

// 计算可用的年份列表（基于当前日期生成，计算前20年）
const availableYears = computed(() => {
  const currentYear = new Date().getFullYear()
  const years: number[] = []
  for (let year = currentYear; year >= currentYear - 20; year--) {
    years.push(year)
  }
  return years
})

// 年份选择变化时的处理
const onYearChange = () => {
  if (!selectedYear.value) {
    // 清空年份时，重置所有筛选
    selectedMonth.value = ''
    filterByMonth(undefined)
    imageVisibility.value = []
    setTimeout(() => {
      setupObserver()
    }, 100)
  } else if (selectedMonth.value) {
    // 如果已经有月份选择，立即触发筛选
    onFilterChange()
  }
}

// 筛选变化时调用服务端筛选
const onFilterChange = () => {
  // 只有同时选择年份和月份时才触发筛选
  if (selectedYear.value && selectedMonth.value) {
    const yearStr = selectedYear.value
    const monthStr = String(selectedMonth.value).padStart(2, '0')
    const monthParam = `${yearStr}-${monthStr}`
    
    // 调用服务端筛选
    filterByMonth(monthParam)
    
    // 重置懒加载状态
    imageVisibility.value = []
    setTimeout(() => {
      setupObserver()
    }, 100)
  }
}

// 重置筛选
const resetFilters = () => {
  selectedYear.value = ''
  selectedMonth.value = ''
  
  // 重置为加载默认数据
  filterByMonth(undefined)
  
  // 重置懒加载状态
  imageVisibility.value = []
  setTimeout(() => {
    setupObserver()
  }, 100)
}

// 设置图片ref
const setImageRef = (el: any, index: number) => {
  if (el && el instanceof HTMLElement) {
    imageRefs.value[index] = el
  }
}

// 设置 Intersection Observer
const setupObserver = () => {
  // 清理旧的 observer
  if (observer) {
    observer.disconnect()
  }
  
  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const index = imageRefs.value.findIndex(ref => ref === entry.target)
          if (index !== -1) {
            imageVisibility.value[index] = true
            // 加载后取消观察
            observer?.unobserve(entry.target)
          }
        }
      })
    },
    {
      root: null,
      rootMargin: '200px', // 提前 200px 开始加载
      threshold: 0.01
    }
  )
  
  // 观察所有图片元素
  imageRefs.value.forEach((ref, index) => {
    if (ref && !imageVisibility.value[index]) {
      observer?.observe(ref)
    }
  })
}


// 监听 images 变化，动态更新 imageVisibility
watch(() => images.value.length, (newLength, oldLength) => {
  if (newLength > oldLength) {
    // 图片列表增长时，扩展 imageVisibility 数组
    const newItems = new Array(newLength - oldLength).fill(false)
    imageVisibility.value = [...imageVisibility.value, ...newItems]
    
    // 重新设置 observer
    setTimeout(() => {
      setupObserver()
    }, 100)
  } else if (newLength < oldLength) {
    // 图片列表减少时（如筛选），截断数组
    imageVisibility.value = imageVisibility.value.slice(0, newLength)
  } else if (newLength === 0) {
    // 如果是首次加载
    imageVisibility.value = []
  }
  
  // 如果从 0 变为有数据，初始化并设置 observer
  if (oldLength === 0 && newLength > 0) {
    imageVisibility.value = new Array(newLength).fill(false)
    setTimeout(() => {
      setupObserver()
    }, 100)
  }
})

// 设置无限滚动 Observer
const setupLoadMoreObserver = () => {
  if (loadMoreObserver) {
    loadMoreObserver.disconnect()
  }
  
  loadMoreObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting && !loading.value && hasMore.value) {
          loadMore()
        }
      })
    },
    {
      root: null,
      rootMargin: '100px',
      threshold: 0.1
    }
  )
  
  if (loadMoreTrigger.value) {
    loadMoreObserver.observe(loadMoreTrigger.value)
  }
}

// 初始化
onMounted(() => {
  loadLatestImage()
  
  if (images.value.length > 0) {
    imageVisibility.value = new Array(images.value.length).fill(false)
    setTimeout(() => {
      setupObserver()
      setupLoadMoreObserver()
    }, 100)
  }
})

// 清理
onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
  if (loadMoreObserver) {
    loadMoreObserver.disconnect()
  }
})

// 格式化日期
const formatDate = (dateStr?: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

// 获取最新图片 URL（顶部大图使用UHD高清）
const getLatestImageUrl = () => {
  if (!latestImage.value?.date) return ''
  return bingPaperApi.getImageUrlByDate(latestImage.value.date, 'UHD', 'jpg')
}

// 获取图片 URL（缩略图 - 使用较小分辨率节省流量）
const getImageUrl = (date: string) => {
  return bingPaperApi.getImageUrlByDate(date, '640x480', 'jpg')
}

// 查看图片详情
const viewImage = (date: string) => {
  router.push(`/image/${date}`)
}

// 打开版权详情链接
const openCopyrightLink = (link: string) => {
  // copyrightlink 是完整的 URL，直接打开
  window.open(link, '_blank')
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  overflow: hidden;
}
</style>

<style>
/* 隐藏滚动条但保持滚动功能 */
body {
  overflow-y: scroll;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

body::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

html {
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

html::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}
</style>

<template>
  <div class="fixed inset-0 bg-black z-50 overflow-hidden">
    <!-- 加载状态 -->
    <div v-if="loading" class="absolute inset-0 flex items-center justify-center">
      <div class="w-16 h-16 border-4 border-white/20 border-t-white rounded-full animate-spin"></div>
    </div>

    <!-- 主要内容 -->
    <div v-else-if="image" class="relative h-full w-full">
      <!-- 全屏图片 -->
      <div class="absolute inset-0 flex items-center justify-center">
        <img 
          :src="getFullImageUrl()" 
          :alt="image.title || 'Bing Image'"
          class="max-w-full max-h-full object-contain"
        />
      </div>

      <!-- 顶部工具栏 -->
      <div class="absolute top-0 left-0 right-0 bg-gradient-to-b from-black/80 to-transparent p-6 z-10">
        <div class="flex items-center justify-between max-w-7xl mx-auto">
          <button 
            @click="goBack"
            class="flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-md text-white rounded-lg hover:bg-white/20 transition-all"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
            <span>返回</span>
          </button>

          <div class="text-white/80 text-sm">
            {{ formatDate(image.date) }}
          </div>
        </div>
      </div>

      <!-- 信息悬浮层（类似 Windows 聚焦） -->
      <div 
        v-if="showInfo"
        class="absolute bottom-24 left-8 right-8 md:left-16 md:right-auto md:max-w-md bg-black/60 backdrop-blur-xl rounded-2xl p-6 transform transition-all duration-500 z-10"
        :class="{ 'translate-y-0 opacity-100': showInfo, 'translate-y-4 opacity-0': !showInfo }"
      >
        <h2 class="text-2xl font-bold text-white mb-3">
          {{ image.title || '未命名' }}
        </h2>
        
        <p v-if="image.copyright" class="text-white/80 text-sm mb-4 leading-relaxed">
          {{ image.copyright }}
        </p>

        <!-- 版权详情链接 -->
        <a 
          v-if="image.copyrightlink"
          :href="image.copyrightlink"
          target="_blank"
          class="inline-flex items-center gap-2 px-4 py-2 bg-white/20 hover:bg-white/30 text-white rounded-lg text-sm font-medium transition-all group"
        >
          <span>了解更多信息</span>
          <svg class="w-4 h-4 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
          </svg>
        </a>

        <!-- 切换信息显示按钮 -->
        <button 
          @click="showInfo = false"
          class="absolute top-4 right-4 p-2 hover:bg-white/10 rounded-lg transition-all"
        >
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- 底部控制栏 -->
      <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-6 z-10">
        <div class="flex items-center justify-between max-w-7xl mx-auto">
          <!-- 日期切换按钮 -->
          <div class="flex items-center gap-4">
            <button 
              @click="previousDay"
              :disabled="navigating"
              class="flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-md text-white rounded-lg hover:bg-white/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
              <span class="hidden sm:inline">前一天</span>
            </button>

            <button 
              @click="nextDay"
              :disabled="navigating || isToday"
              class="flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-md text-white rounded-lg hover:bg-white/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span class="hidden sm:inline">后一天</span>
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </button>
          </div>

          <!-- 信息按钮 -->
          <button 
            v-if="!showInfo"
            @click="showInfo = true"
            class="flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-md text-white rounded-lg hover:bg-white/20 transition-all"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span class="hidden sm:inline">显示信息</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="absolute inset-0 flex items-center justify-center">
      <div class="text-center">
        <p class="text-white/60 mb-4">加载失败</p>
        <button 
          @click="goBack"
          class="px-6 py-3 bg-white/10 backdrop-blur-md text-white rounded-lg hover:bg-white/20 transition-all"
        >
          返回首页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useImageByDate } from '@/composables/useImages'
import { bingPaperApi } from '@/lib/api-service'

const route = useRoute()
const router = useRouter()

const currentDate = ref(route.params.date as string)
const showInfo = ref(true)
const navigating = ref(false)

// 使用 composable 获取图片数据
const { image, loading, error, refetch } = useImageByDate(currentDate.value)

// 监听日期变化
watch(currentDate, () => {
  refetch()
})

// 格式化日期
const formatDate = (dateStr?: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric',
    weekday: 'long'
  })
}

// 判断是否是今天
const isToday = computed(() => {
  const today = new Date().toISOString().split('T')[0]
  return currentDate.value === today
})

// 获取完整图片 URL
const getFullImageUrl = () => {
  return bingPaperApi.getImageUrlByDate(currentDate.value, 'UHD', 'jpg')
}

// copyrightlink 现在是完整的 URL，无需额外处理

// 返回首页
const goBack = () => {
  router.push('/')
}

// 前一天
const previousDay = () => {
  if (navigating.value) return
  
  navigating.value = true
  const date = new Date(currentDate.value)
  date.setDate(date.getDate() - 1)
  const newDate = date.toISOString().split('T')[0]
  
  currentDate.value = newDate
  router.replace(`/image/${newDate}`)
  
  setTimeout(() => {
    navigating.value = false
  }, 500)
}

// 后一天
const nextDay = () => {
  if (navigating.value || isToday.value) return
  
  navigating.value = true
  const date = new Date(currentDate.value)
  date.setDate(date.getDate() + 1)
  const newDate = date.toISOString().split('T')[0]
  
  currentDate.value = newDate
  router.replace(`/image/${newDate}`)
  
  setTimeout(() => {
    navigating.value = false
  }, 500)
}

// 键盘导航
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'ArrowLeft') {
    previousDay()
  } else if (e.key === 'ArrowRight' && !isToday.value) {
    nextDay()
  } else if (e.key === 'Escape') {
    goBack()
  } else if (e.key === 'i' || e.key === 'I') {
    showInfo.value = !showInfo.value
  }
}

// 添加键盘事件监听
if (typeof window !== 'undefined') {
  window.addEventListener('keydown', handleKeydown)
}

// 清理
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', handleKeydown)
  }
})
</script>

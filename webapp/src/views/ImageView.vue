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
        ref="infoPanel"
        class="fixed w-[90%] max-w-md bg-black/40 backdrop-blur-lg rounded-xl p-4 transform transition-opacity duration-300 z-10 select-none"
        :style="{ left: infoPanelPos.x + 'px', top: infoPanelPos.y + 'px' }"
        :class="{ 'opacity-100': showInfo, 'opacity-0': !showInfo }"
      >
        <!-- 拖动手柄 -->
        <div 
          @mousedown="startDrag"
          @touchstart="startDrag"
          class="absolute top-2 left-1/2 -translate-x-1/2 w-12 h-1 bg-white/30 rounded-full cursor-move hover:bg-white/50 transition-colors touch-none"
        ></div>

        <h2 class="text-lg font-bold text-white mb-2 mt-2">
          {{ image.title || '未命名' }}
        </h2>
        
        <p v-if="image.copyright" class="text-white/80 text-xs mb-3 leading-relaxed">
          {{ image.copyright }}
        </p>

        <!-- 版权详情链接 -->
        <a 
          v-if="image.copyrightlink"
          :href="image.copyrightlink"
          target="_blank"
          class="inline-flex items-center gap-2 px-3 py-1.5 bg-white/15 hover:bg-white/25 text-white rounded-lg text-xs font-medium transition-all group"
        >
          <span>了解更多</span>
          <svg class="w-3 h-3 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
          </svg>
        </a>

        <!-- 切换信息显示按钮 -->
        <button 
          @click="showInfo = false"
          class="absolute top-3 right-3 p-1.5 hover:bg-white/10 rounded-lg transition-all"
        >
          <svg class="w-4 h-4 text-white/80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
              :disabled="navigating || !hasPreviousDay"
              class="flex items-center gap-2 px-4 py-2 bg-white/10 backdrop-blur-md text-white rounded-lg hover:bg-white/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
              <span class="hidden sm:inline">前一天</span>
            </button>

            <button 
              @click="nextDay"
              :disabled="navigating || !hasNextDay"
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
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useImageByDate } from '@/composables/useImages'
import { bingPaperApi } from '@/lib/api-service'

const route = useRoute()
const router = useRouter()

const currentDate = ref(route.params.date as string)
const showInfo = ref(true)
const navigating = ref(false)

// 前后日期可用性
const hasPreviousDay = ref(true)
const hasNextDay = ref(true)
const checkingDates = ref(false)

// 拖动相关状态
const infoPanel = ref<HTMLElement | null>(null)
const infoPanelPos = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })

// 初始化浮窗位置（居中偏下）
const initPanelPosition = () => {
  if (typeof window !== 'undefined') {
    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight
    const panelWidth = Math.min(windowWidth * 0.9, 448) // max-w-md = 448px
    infoPanelPos.value = {
      x: (windowWidth - panelWidth) / 2,
      y: windowHeight - 200 // 距底部200px
    }
  }
}

// 开始拖动
const startDrag = (e: MouseEvent | TouchEvent) => {
  e.preventDefault()
  isDragging.value = true
  
  const clientX = e instanceof MouseEvent ? e.clientX : e.touches[0].clientX
  const clientY = e instanceof MouseEvent ? e.clientY : e.touches[0].clientY
  
  dragStart.value = {
    x: clientX - infoPanelPos.value.x,
    y: clientY - infoPanelPos.value.y
  }
  
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
  document.addEventListener('touchmove', onDrag, { passive: false })
  document.addEventListener('touchend', stopDrag)
}

// 拖动中
const onDrag = (e: MouseEvent | TouchEvent) => {
  if (!isDragging.value) return
  
  if (e instanceof TouchEvent) {
    e.preventDefault()
  }
  
  const clientX = e instanceof MouseEvent ? e.clientX : e.touches[0].clientX
  const clientY = e instanceof MouseEvent ? e.clientY : e.touches[0].clientY
  
  const newX = clientX - dragStart.value.x
  const newY = clientY - dragStart.value.y
  
  // 限制在视口内
  if (infoPanel.value) {
    const rect = infoPanel.value.getBoundingClientRect()
    const maxX = window.innerWidth - rect.width
    const maxY = window.innerHeight - rect.height
    
    infoPanelPos.value = {
      x: Math.max(0, Math.min(newX, maxX)),
      y: Math.max(0, Math.min(newY, maxY))
    }
  }
}

// 停止拖动
const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)
}

// 使用 composable 获取图片数据（传递 ref，自动响应日期变化）
const { image, loading, error } = useImageByDate(currentDate)

// 检测指定日期是否有数据
const checkDateAvailability = async (dateStr: string): Promise<boolean> => {
  try {
    await bingPaperApi.getImageMetaByDate(dateStr)
    return true
  } catch (e) {
    return false
  }
}

// 检测前后日期可用性
const checkAdjacentDates = async () => {
  if (checkingDates.value) return
  
  checkingDates.value = true
  const date = new Date(currentDate.value)
  
  // 检测前一天
  const prevDate = new Date(date)
  prevDate.setDate(prevDate.getDate() - 1)
  hasPreviousDay.value = await checkDateAvailability(prevDate.toISOString().split('T')[0])
  
  // 检测后一天（不能超过今天）
  const nextDate = new Date(date)
  nextDate.setDate(nextDate.getDate() + 1)
  const today = new Date().toISOString().split('T')[0]
  if (nextDate.toISOString().split('T')[0] > today) {
    hasNextDay.value = false
  } else {
    hasNextDay.value = await checkDateAvailability(nextDate.toISOString().split('T')[0])
  }
  
  checkingDates.value = false
}

// 初始化位置
initPanelPosition()

// 监听日期变化，检测前后日期可用性
import { watch } from 'vue'
watch(currentDate, () => {
  checkAdjacentDates()
}, { immediate: true })

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
  if (navigating.value || !hasPreviousDay.value) return
  
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
  if (navigating.value || !hasNextDay.value) return
  
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
  if (e.key === 'ArrowLeft' && hasPreviousDay.value) {
    previousDay()
  } else if (e.key === 'ArrowRight' && hasNextDay.value) {
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
  window.addEventListener('resize', initPanelPosition)
}

// 清理
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', handleKeydown)
    window.removeEventListener('resize', initPanelPosition)
    document.removeEventListener('mousemove', onDrag)
    document.removeEventListener('mouseup', stopDrag)
    document.removeEventListener('touchmove', onDrag)
    document.removeEventListener('touchend', stopDrag)
  }
})
</script>

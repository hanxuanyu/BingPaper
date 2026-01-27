<template>
  <div class="min-h-screen bg-gradient-to-b from-gray-900 via-gray-800 to-gray-900">
    <!-- Hero Section - 今日图片 -->
    <section class="relative h-screen w-full overflow-hidden">
      <div v-if="todayLoading" class="absolute inset-0 flex items-center justify-center">
        <div class="w-12 h-12 border-4 border-white/20 border-t-white rounded-full animate-spin"></div>
      </div>
      
      <div v-else-if="todayImage" class="relative h-full w-full group">
        <!-- 背景图片 -->
        <div class="absolute inset-0">
          <img 
            :src="getTodayImageUrl()" 
            :alt="todayImage.title || 'Today\'s Bing Image'"
            class="w-full h-full object-cover"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>
        </div>

        <!-- 内容叠加层 -->
        <div class="relative h-full flex flex-col justify-end p-8 md:p-16 z-10">
          <div class="max-w-4xl space-y-4 transform transition-transform duration-500 group-hover:translate-y-[-10px]">
            <div class="inline-block px-4 py-2 bg-white/10 backdrop-blur-md rounded-full text-white/90 text-sm font-medium">
              今日精选 · {{ formatDate(todayImage.date) }}
            </div>
            
            <h1 class="text-4xl md:text-6xl font-bold text-white leading-tight drop-shadow-2xl">
              {{ todayImage.title || '必应每日一图' }}
            </h1>
            
            <p v-if="todayImage.copyright" class="text-lg md:text-xl text-white/80 max-w-2xl">
              {{ todayImage.copyright }}
            </p>

            <div class="flex gap-4 pt-4">
              <button 
                @click="viewImage(todayImage.date!)"
                class="px-6 py-3 bg-white text-gray-900 rounded-lg font-semibold hover:bg-white/90 transition-all transform hover:scale-105 shadow-xl"
              >
                查看大图
              </button>
              <button 
                v-if="todayImage.quiz"
                @click="openQuiz(todayImage.quiz)"
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
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-8">
          历史精选
        </h2>

        <!-- 图片网格 -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div 
            v-for="(image, index) in images" 
            :key="image.date || index"
            class="group relative aspect-video rounded-xl overflow-hidden cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-2xl"
            @click="viewImage(image.date!)"
          >
            <!-- 图片 -->
            <img 
              :src="getImageUrl(image.date!)" 
              :alt="image.title || 'Bing Image'"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
              loading="lazy"
            />
            
            <!-- 悬浮信息层 -->
            <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/50 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300">
              <div class="absolute bottom-0 left-0 right-0 p-6 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300">
                <div class="text-xs text-white/70 mb-2">
                  {{ formatDate(image.date) }}
                </div>
                <h3 class="text-lg font-semibold text-white mb-2 line-clamp-2">
                  {{ image.title || '未命名' }}
                </h3>
                <p v-if="image.copyright" class="text-sm text-white/80 line-clamp-2">
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
          
          <button 
            v-else-if="hasMore"
            @click="loadMore"
            class="px-8 py-3 bg-white/10 backdrop-blur-md text-white rounded-lg font-semibold hover:bg-white/20 transition-all border border-white/30"
          >
            加载更多
          </button>

          <p v-else class="text-white/40">
            已加载全部图片
          </p>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="py-8 text-center text-white/40 border-t border-white/10">
      <p>数据来源于必应每日一图 API</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useTodayImage, useImageList } from '@/composables/useImages'
import { bingPaperApi } from '@/lib/api-service'
import { useRouter } from 'vue-router'

const router = useRouter()

// 获取今日图片
const { image: todayImage, loading: todayLoading } = useTodayImage()

// 获取图片列表
const { images, loading, hasMore, loadMore } = useImageList(30)

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

// 获取今日图片 URL
const getTodayImageUrl = () => {
  return bingPaperApi.getTodayImageUrl('UHD', 'jpg')
}

// 获取图片 URL
const getImageUrl = (date: string) => {
  return bingPaperApi.getImageUrlByDate(date, '1920x1080', 'jpg')
}

// 查看图片详情
const viewImage = (date: string) => {
  router.push(`/image/${date}`)
}

// 打开必应 quiz 链接
const openQuiz = (quiz: string) => {
  // 拼接完整的必应地址
  const bingUrl = `https://www.bing.com${quiz}`
  window.open(bingUrl, '_blank')
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

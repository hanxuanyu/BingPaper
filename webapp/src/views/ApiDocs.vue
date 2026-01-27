<template>
  <div class="min-h-screen bg-gradient-to-b from-gray-900 via-gray-800 to-gray-900">
    <!-- Header -->
    <header class="sticky top-0 bg-gray-900/80 backdrop-blur-lg border-b border-white/10 z-40">
      <div class="max-w-7xl mx-auto px-4 md:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <router-link 
              to="/"
              class="flex items-center gap-2 text-white/60 hover:text-white transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
              </svg>
              <span>返回首页</span>
            </router-link>
            <span class="text-white/20">|</span>
            <h1 class="text-xl font-bold text-white">API 文档</h1>
          </div>
          
          <div class="text-sm text-white/40">
            v1.0
          </div>
        </div>
      </div>
    </header>

    <!-- Content -->
    <main class="max-w-7xl mx-auto px-4 md:px-8 py-12">
      <!-- Intro Section -->
      <section class="mb-12">
        <div class="bg-white/5 backdrop-blur-sm rounded-2xl p-8 border border-white/10">
          <h2 class="text-3xl font-bold text-white mb-4">必应每日一图 API</h2>
          <p class="text-white/70 text-lg leading-relaxed">
            提供必应每日一图的公共 API 接口，支持获取今日图片、历史图片、随机图片等功能。
            所有接口均为 RESTful 风格，返回 JSON 格式数据或图片流。
          </p>
          <div class="mt-6 flex flex-wrap gap-3">
            <div class="px-4 py-2 bg-green-500/20 text-green-300 rounded-lg text-sm font-medium border border-green-500/30">
              ✓ 无需认证
            </div>
            <div class="px-4 py-2 bg-blue-500/20 text-blue-300 rounded-lg text-sm font-medium border border-blue-500/30">
              RESTful API
            </div>
            <div class="px-4 py-2 bg-purple-500/20 text-purple-300 rounded-lg text-sm font-medium border border-purple-500/30">
              JSON / 图片流
            </div>
          </div>
        </div>
      </section>

      <!-- Base URL -->
      <section class="mb-12">
        <h2 class="text-2xl font-bold text-white mb-6">基础地址</h2>
        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
          <div class="flex items-center gap-3 mb-2">
            <code class="text-green-400 font-mono">{{ baseURL }}</code>
            <button 
              @click="copyToClipboard(baseURL)"
              class="p-2 hover:bg-white/10 rounded-lg transition-colors"
              title="复制"
            >
              <svg class="w-4 h-4 text-white/60" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
              </svg>
            </button>
          </div>
          <p class="text-white/50 text-sm">所有 API 请求都基于此地址</p>
        </div>
      </section>

      <!-- Image APIs -->
      <section class="mb-12">
        <h2 class="text-2xl font-bold text-white mb-6">图片 API</h2>
        
        <!-- Get Today's Image -->
        <div class="mb-8">
          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-start justify-between mb-4">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <span class="px-3 py-1 bg-green-500/20 text-green-300 rounded-lg text-sm font-semibold">GET</span>
                  <h3 class="text-xl font-semibold text-white">获取今日图片</h3>
                </div>
                <code class="text-blue-400 font-mono text-sm">/image/today</code>
              </div>
            </div>
            
            <p class="text-white/70 mb-4">返回今日必应图片，支持不同分辨率和格式</p>
            
            <!-- Parameters -->
            <div class="mb-4">
              <h4 class="text-white/80 font-semibold mb-2">查询参数</h4>
              <div class="space-y-2">
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">variant</code>
                  <div class="flex-1">
                    <span class="text-white/50 block mb-1">分辨率 (默认: UHD)</span>
                    <span class="text-white/40 text-xs">可选值: UHD, 1920x1080, 1366x768, 1280x720, 1024x768, 800x600, 800x480, 640x480, 640x360, 480x360, 400x240, 320x240</span>
                  </div>
                </div>
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">format</code>
                  <span class="text-white/50">格式: jpg (默认: jpg)</span>
                </div>
              </div>
            </div>

            <!-- Example -->
            <div class="bg-black/30 rounded-lg p-4 mb-4">
              <div class="flex items-center justify-between mb-2">
                <span class="text-white/50 text-sm">示例</span>
                <button 
                  @click="copyToClipboard(getTodayImageExample())"
                  class="text-white/60 hover:text-white text-sm flex items-center gap-1"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                  </svg>
                  复制
                </button>
              </div>
              <code class="text-green-400 font-mono text-sm">{{ getTodayImageExample() }}</code>
            </div>

            <!-- Try It -->
            <div class="flex gap-3">
              <a 
                :href="getTodayImageExample()"
                target="_blank"
                class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
                </svg>
                在新窗口打开
              </a>
              <button 
                @click="showImagePreview(getTodayImageExample())"
                class="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg text-sm font-medium transition-colors"
              >
                预览图片
              </button>
            </div>
          </div>
        </div>

        <!-- Get Image by Date -->
        <div class="mb-8">
          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-start justify-between mb-4">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <span class="px-3 py-1 bg-green-500/20 text-green-300 rounded-lg text-sm font-semibold">GET</span>
                  <h3 class="text-xl font-semibold text-white">获取指定日期图片</h3>
                </div>
                <code class="text-blue-400 font-mono text-sm">/image/date/:date</code>
              </div>
            </div>
            
            <p class="text-white/70 mb-4">根据日期返回对应的必应图片</p>
            
            <!-- Parameters -->
            <div class="mb-4">
              <h4 class="text-white/80 font-semibold mb-2">路径参数</h4>
              <div class="space-y-2 mb-3">
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">date</code>
                  <span class="text-white/50">日期 (格式: YYYY-MM-DD)</span>
                </div>
              </div>
              <h4 class="text-white/80 font-semibold mb-2">查询参数</h4>
              <div class="space-y-2">
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">variant</code>
                  <div class="flex-1">
                    <span class="text-white/50 block mb-1">分辨率 (默认: UHD)</span>
                    <span class="text-white/40 text-xs">可选值: UHD, 1920x1080, 1366x768, 1280x720, 1024x768, 800x600, 800x480, 640x480, 640x360, 480x360, 400x240, 320x240</span>
                  </div>
                </div>
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">format</code>
                  <span class="text-white/50">格式 (默认: jpg)</span>
                </div>
              </div>
            </div>

            <!-- Example -->
            <div class="bg-black/30 rounded-lg p-4 mb-4">
              <div class="flex items-center justify-between mb-2">
                <span class="text-white/50 text-sm">示例</span>
                <button 
                  @click="copyToClipboard(getDateImageExample())"
                  class="text-white/60 hover:text-white text-sm flex items-center gap-1"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                  </svg>
                  复制
                </button>
              </div>
              <code class="text-green-400 font-mono text-sm">{{ getDateImageExample() }}</code>
            </div>

            <!-- Try It -->
            <div class="flex gap-3">
              <a 
                :href="getDateImageExample()"
                target="_blank"
                class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
                </svg>
                在新窗口打开
              </a>
              <button 
                @click="showImagePreview(getDateImageExample())"
                class="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg text-sm font-medium transition-colors"
              >
                预览图片
              </button>
            </div>
          </div>
        </div>

        <!-- Get Random Image -->
        <div class="mb-8">
          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-start justify-between mb-4">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <span class="px-3 py-1 bg-green-500/20 text-green-300 rounded-lg text-sm font-semibold">GET</span>
                  <h3 class="text-xl font-semibold text-white">获取随机图片</h3>
                </div>
                <code class="text-blue-400 font-mono text-sm">/image/random</code>
              </div>
            </div>
            
            <p class="text-white/70 mb-4">随机返回一张历史图片</p>
            
            <!-- Parameters -->
            <div class="mb-4">
              <h4 class="text-white/80 font-semibold mb-2">查询参数</h4>
              <div class="space-y-2">
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">variant</code>
                  <div class="flex-1">
                    <span class="text-white/50 block mb-1">分辨率 (默认: UHD)</span>
                    <span class="text-white/40 text-xs">可选值: UHD, 1920x1080, 1366x768, 1280x720, 1024x768, 800x600, 800x480, 640x480, 640x360, 480x360, 400x240, 320x240</span>
                  </div>
                </div>
                <div class="flex gap-4 text-sm">
                  <code class="text-yellow-400 min-w-24">format</code>
                  <span class="text-white/50">格式 (默认: jpg)</span>
                </div>
              </div>
            </div>

            <!-- Example -->
            <div class="bg-black/30 rounded-lg p-4 mb-4">
              <div class="flex items-center justify-between mb-2">
                <span class="text-white/50 text-sm">示例</span>
                <button 
                  @click="copyToClipboard(getRandomImageExample())"
                  class="text-white/60 hover:text-white text-sm flex items-center gap-1"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                  </svg>
                  复制
                </button>
              </div>
              <code class="text-green-400 font-mono text-sm">{{ getRandomImageExample() }}</code>
            </div>

            <!-- Try It -->
            <div class="flex gap-3">
              <a 
                :href="getRandomImageExample()"
                target="_blank"
                class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
                </svg>
                在新窗口打开
              </a>
              <button 
                @click="showImagePreview(getRandomImageExample())"
                class="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg text-sm font-medium transition-colors"
              >
                预览图片
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- Metadata APIs -->
      <section class="mb-12">
        <h2 class="text-2xl font-bold text-white mb-6">元数据 API</h2>
        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
          <p class="text-white/70 mb-4">获取图片的元数据信息（标题、版权、日期等），只返回 JSON 数据不返回图片</p>
          
          <div class="space-y-3 mb-6">
            <div class="flex items-center gap-3 text-sm">
              <code class="text-blue-400 font-mono">/image/today/meta</code>
              <span class="text-white/50">-</span>
              <span class="text-white/60">今日图片元数据</span>
            </div>
            <div class="flex items-center gap-3 text-sm">
              <code class="text-blue-400 font-mono">/image/date/:date/meta</code>
              <span class="text-white/50">-</span>
              <span class="text-white/60">指定日期图片元数据</span>
            </div>
            <div class="flex items-center gap-3 text-sm">
              <code class="text-blue-400 font-mono">/image/random/meta</code>
              <span class="text-white/50">-</span>
              <span class="text-white/60">随机图片元数据</span>
            </div>
            <div class="flex items-center gap-3 text-sm">
              <code class="text-blue-400 font-mono">/images?limit=30</code>
              <span class="text-white/50">-</span>
              <span class="text-white/60">图片列表（支持分页）</span>
            </div>
          </div>

          <!-- 元数据字段说明 -->
          <div class="mt-6 pt-6 border-t border-white/10">
            <h4 class="text-white/80 font-semibold mb-4">响应字段说明</h4>
            <div class="space-y-3 text-sm">
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">date</code>
                <span class="text-white/60">图片日期（格式：YYYY-MM-DD）</span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">title</code>
                <span class="text-white/60">图片标题</span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">copyright</code>
                <span class="text-white/60">版权信息</span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">copyrightlink</code>
                <span class="text-white/60">版权详情链接（指向 Bing 搜索页面）</span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">startdate</code>
                <span class="text-white/60">发布开始日期（格式：YYYYMMDD）</span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">fullstartdate</code>
                <span class="text-white/60">完整发布时间（格式：YYYYMMDDHHMM）</span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">hsh</code>
                <span class="text-white/60">图片唯一哈希值 </span>
              </div>
              <div class="flex gap-4">
                <code class="text-yellow-400 min-w-32">quiz</code>
                <span class="text-white/60">必应 quiz 链接（已废弃，建议使用 copyrightlink）</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Usage Tips -->
      <section class="mb-12">
        <h2 class="text-2xl font-bold text-white mb-6">使用提示</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-10 h-10 bg-blue-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-white">直接使用</h3>
            </div>
            <p class="text-white/60 text-sm">
              所有图片 API 都可以直接在 HTML <code class="text-yellow-400">&lt;img&gt;</code> 标签中使用，无需认证。
            </p>
          </div>

          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-10 h-10 bg-green-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-white">CORS 支持</h3>
            </div>
            <p class="text-white/60 text-sm">
              API 支持跨域请求，可以从任何网站调用，适合用作壁纸服务。
            </p>
          </div>

          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-white">多种分辨率</h3>
            </div>
            <p class="text-white/60 text-sm">
              支持 UHD、1920x1080、1366x768 等多种分辨率，适配不同设备。
            </p>
          </div>

          <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-10 h-10 bg-orange-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-white">每日更新</h3>
            </div>
            <p class="text-white/60 text-sm">
              图片数据每日自动更新，与必应官方保持同步。
            </p>
          </div>
        </div>
      </section>
    </main>

    <!-- Image Preview Modal -->
    <div 
      v-if="previewImage"
      class="fixed inset-0 bg-black/90 z-50 flex items-center justify-center p-4"
      @click="previewImage = null"
    >
      <div class="relative max-w-6xl w-full">
        <button 
          @click="previewImage = null"
          class="absolute top-4 right-4 p-2 bg-white/10 hover:bg-white/20 rounded-lg text-white transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
        <img 
          :src="previewImage" 
          alt="Preview"
          class="w-full h-auto rounded-lg shadow-2xl"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { API_BASE_URL } from '@/lib/api-config'

const baseURL = ref(API_BASE_URL)
const previewImage = ref<string | null>(null)

// 获取今日图片示例
const getTodayImageExample = () => {
  return `${baseURL.value}/image/today?variant=UHD&format=jpg`
}

// 获取指定日期图片示例
const getDateImageExample = () => {
  const today = new Date().toISOString().split('T')[0]
  return `${baseURL.value}/image/date/${today}?variant=1920x1080&format=jpg`
}

// 获取随机图片示例
const getRandomImageExample = () => {
  return `${baseURL.value}/image/random?variant=UHD&format=jpg`
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    // 可以添加一个 toast 提示
    console.log('已复制到剪贴板')
  } catch (err) {
    console.error('复制失败:', err)
  }
}

// 显示图片预览
const showImagePreview = (url: string) => {
  previewImage.value = url
}
</script>

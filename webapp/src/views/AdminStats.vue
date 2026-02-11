<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-semibold">
        调用统计 Dashboard
        <span v-if="selectedEndpoint" class="text-sm font-normal text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full ml-2">
          接口: {{ selectedEndpoint }}
        </span>
      </h3>
      <div class="flex gap-2">
        <Button v-if="selectedEndpoint" variant="ghost" size="sm" @click="resetFilter">
          清除筛选
        </Button>
        <Button variant="outline" size="sm" @click="fetchData" :disabled="loading">
          <RefreshCw class="w-4 h-4 mr-2" :class="{ 'animate-spin': loading }" />
          刷新数据
        </Button>
      </div>
    </div>

    <!-- 概览卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium text-gray-500">今日调用</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ summary.today }}</div>
          <p class="text-xs text-gray-400 mt-1">
            较昨日: {{ summary.today - summary.yesterday >= 0 ? '+' : '' }}{{ summary.today - summary.yesterday }}
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium text-gray-500">昨日调用</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ summary.yesterday }}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="pb-2">
          <CardTitle class="text-sm font-medium text-gray-500">累计调用</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ summary.total }}</div>
        </CardContent>
      </Card>
    </div>

    <!-- 趋势图 (简单实现) -->
    <Card>
      <CardHeader class="flex flex-row items-center justify-between">
        <CardTitle>{{ selectedEndpoint ? '接口详细趋势' : '全局调用趋势' }} (14天)</CardTitle>
        <div v-if="selectedEndpoint" class="text-xs text-gray-500 font-mono">
          {{ selectedEndpoint }}
        </div>
      </CardHeader>
      <CardContent>
        <div class="h-64 flex items-end gap-1 px-2 border-b border-l pt-8 relative">
          <!-- 刻度线 (简单模拟) -->
          <div class="absolute left-0 right-0 top-8 border-t border-gray-100 border-dashed"></div>
          <div class="absolute left-0 right-0 top-1/2 border-t border-gray-100 border-dashed"></div>
          
          <div v-for="item in trend" :key="item.date" class="flex-1 flex flex-col items-center group relative">
            <div 
              class="w-full bg-blue-500 rounded-t transition-all duration-300 group-hover:bg-blue-600 relative"
              :style="{ height: `${maxTrendCount > 0 ? (item.count / maxTrendCount) * 100 : 0}%` }"
            >
              <!-- Tooltip 效果 -->
              <div class="absolute -top-8 left-1/2 -translate-x-1/2 bg-gray-800 text-white text-[10px] py-1 px-2 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10">
                {{ item.date }}: {{ item.count }}
              </div>
            </div>
            <div class="text-[9px] mt-2 text-gray-500 truncate w-full text-center">{{ item.date.split('-').slice(1).join('/') }}</div>
          </div>
        </div>
      </CardContent>
    </Card>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 接口分布 -->
      <Card>
        <CardHeader>
          <CardTitle>接口调用排行</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-3">
            <div 
              v-for="item in endpoints" 
              :key="item.endpoint" 
              class="space-y-1 cursor-pointer hover:bg-gray-50 p-2 rounded-lg transition-colors group"
              :class="{ 'bg-blue-50 ring-1 ring-blue-100': selectedEndpoint === item.endpoint }"
              @click="selectEndpoint(item.endpoint)"
            >
              <div class="flex justify-between text-sm items-center">
                <span class="font-mono text-[11px] truncate mr-2" :class="{ 'text-blue-600 font-bold': selectedEndpoint === item.endpoint }" :title="item.endpoint">
                  {{ item.endpoint }}
                </span>
                <div class="flex items-center gap-2">
                  <span class="font-bold text-xs">{{ item.count }}</span>
                  <ChevronRight class="w-3 h-3 text-gray-300 group-hover:text-blue-400" />
                </div>
              </div>
              <div class="w-full bg-gray-100 rounded-full h-1.5 mt-1">
                <div 
                  class="bg-blue-400 h-1.5 rounded-full transition-all" 
                  :class="{ 'bg-blue-600': selectedEndpoint === item.endpoint }"
                  :style="{ width: `${maxEndpointCount > 0 ? (item.count / maxEndpointCount) * 100 : 0}%` }"
                ></div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 地区分布 -->
      <Card>
        <CardHeader>
          <CardTitle>地区调用排行</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div v-for="item in regions" :key="item.mkt" class="space-y-1">
              <div class="flex justify-between text-sm">
                <span class="text-xs">{{ getRegionLabel(item.mkt) }}</span>
                <span class="font-bold text-xs">{{ item.count }}</span>
              </div>
              <div class="w-full bg-gray-100 rounded-full h-1.5">
                <div 
                  class="bg-orange-500 h-1.5 rounded-full" 
                  :style="{ width: `${maxRegionCount > 0 ? (item.count / maxRegionCount) * 100 : 0}%` }"
                ></div>
              </div>
            </div>
            <div v-if="regions.length === 0" class="text-center py-4 text-gray-400 text-sm">
              暂无数据
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RefreshCw, ChevronRight } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { apiService } from '@/lib/api-service'
import type { StatSummary, StatTrendItem, StatEndpointItem, StatRegionItem, Region } from '@/lib/api-types'

const loading = ref(false)
const selectedEndpoint = ref<string | null>(null)
const summary = ref<StatSummary>({ total: 0, today: 0, yesterday: 0 })
const trend = ref<StatTrendItem[]>([])
const endpoints = ref<StatEndpointItem[]>([])
const regions = ref<StatRegionItem[]>([])
const regionList = ref<Region[]>([])

const maxTrendCount = computed(() => Math.max(...trend.value.map(i => i.count), 0))
const maxEndpointCount = computed(() => Math.max(...endpoints.value.map(i => i.count), 0))
const maxRegionCount = computed(() => Math.max(...regions.value.map(i => i.count), 0))

const getRegionLabel = (mkt: string) => {
  if (!mkt || mkt === 'default') return '默认/其他'
  const region = regionList.value.find(r => r.value === mkt)
  return region ? region.label : mkt
}

const selectEndpoint = async (endpoint: string) => {
  if (selectedEndpoint.value === endpoint) {
    selectedEndpoint.value = null
  } else {
    selectedEndpoint.value = endpoint
  }
  await fetchTrend()
}

const resetFilter = async () => {
  selectedEndpoint.value = null
  await fetchTrend()
}

const fetchTrend = async () => {
  try {
    trend.value = await apiService.getStatTrend(14, selectedEndpoint.value || undefined)
  } catch (err) {
    console.error('Failed to fetch trend:', err)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const [s, e, r, rl] = await Promise.all([
      apiService.getStatSummary(),
      apiService.getStatEndpoints(),
      apiService.getStatRegions(),
      apiService.getRegions()
    ])
    summary.value = s
    endpoints.value = e
    regions.value = r
    regionList.value = rl
    await fetchTrend()
  } catch (err) {
    console.error('Failed to fetch stats:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

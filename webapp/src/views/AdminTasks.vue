<template>
  <div class="space-y-6">
    <div>
      <h3 class="text-lg font-semibold mb-4">定时任务管理</h3>
      <p class="text-sm text-gray-600 mb-4">
        手动触发图片抓取和清理任务
      </p>
    </div>

    <!-- 手动抓取 -->
    <Card>
      <CardHeader>
        <CardTitle>手动抓取图片</CardTitle>
        <CardDescription>
          立即从 Bing 抓取最新的图片
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <Label for="fetch-days">抓取天数</Label>
          <Input
            id="fetch-days"
            v-model.number="fetchDays"
            type="number"
            min="1"
            max="30"
            placeholder="输入要抓取的天数，默认 1 天"
          />
          <p class="text-xs text-gray-500">
            指定要抓取的天数（包括今天），最多 30 天
          </p>
        </div>
        <Button 
          @click="handleManualFetch" 
          :disabled="fetchLoading"
          class="w-full sm:w-auto"
        >
          {{ fetchLoading ? '抓取中...' : '开始抓取' }}
        </Button>
      </CardContent>
    </Card>

    <!-- 手动清理 -->
    <Card>
      <CardHeader>
        <CardTitle>手动清理旧图片</CardTitle>
        <CardDescription>
          清理超过保留期限的旧图片
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <p class="text-sm text-gray-600">
          根据系统配置的保留天数，清理过期的图片文件和数据库记录
        </p>
        <Button 
          @click="handleManualCleanup" 
          :disabled="cleanupLoading"
          variant="destructive"
          class="w-full sm:w-auto"
        >
          {{ cleanupLoading ? '清理中...' : '开始清理' }}
        </Button>
      </CardContent>
    </Card>

    <!-- 任务历史记录 -->
    <Card>
      <CardHeader>
        <CardTitle>任务执行历史</CardTitle>
        <CardDescription>
          最近的任务执行记录
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="taskHistory.length === 0" class="text-center py-8 text-gray-500">
          暂无执行记录
        </div>
        <div v-else class="space-y-2">
          <div
            v-for="(task, index) in taskHistory"
            :key="index"
            class="flex items-center justify-between p-3 border rounded-md"
          >
            <div class="flex-1">
              <div class="flex items-center gap-2">
                <Badge :variant="task.success ? 'default' : 'destructive'">
                  {{ task.type }}
                </Badge>
                <span class="text-sm">{{ task.message }}</span>
              </div>
              <div class="text-xs text-gray-500 mt-1">
                {{ task.timestamp }}
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { apiService } from '@/lib/api-service'

interface TaskRecord {
  type: string
  success: boolean
  message: string
  timestamp: string
}

const fetchDays = ref<number>(1)
const fetchLoading = ref(false)

const cleanupLoading = ref(false)

const taskHistory = ref<TaskRecord[]>([])

const handleManualFetch = async () => {
  fetchLoading.value = true
  
  try {
    const response = await apiService.manualFetch({ n: fetchDays.value })
    toast.success(response.message || '抓取任务已启动')
    
    // 添加到历史记录
    taskHistory.value.unshift({
      type: '图片抓取',
      success: true,
      message: `抓取 ${fetchDays.value} 天的图片`,
      timestamp: new Date().toLocaleString('zh-CN')
    })
    
    // 只保留最近 10 条记录
    if (taskHistory.value.length > 10) {
      taskHistory.value = taskHistory.value.slice(0, 10)
    }
  } catch (err: any) {
    toast.error(err.message || '抓取失败')
    
    taskHistory.value.unshift({
      type: '图片抓取',
      success: false,
      message: err.message || '抓取失败',
      timestamp: new Date().toLocaleString('zh-CN')
    })
  } finally {
    fetchLoading.value = false
  }
}

const handleManualCleanup = async () => {
  cleanupLoading.value = true
  
  try {
    const response = await apiService.manualCleanup()
    toast.success(response.message || '清理任务已完成')
    
    taskHistory.value.unshift({
      type: '清理任务',
      success: true,
      message: '清理旧图片',
      timestamp: new Date().toLocaleString('zh-CN')
    })
    
    if (taskHistory.value.length > 10) {
      taskHistory.value = taskHistory.value.slice(0, 10)
    }
  } catch (err: any) {
    toast.error(err.message || '清理失败')
    
    taskHistory.value.unshift({
      type: '清理任务',
      success: false,
      message: err.message || '清理失败',
      timestamp: new Date().toLocaleString('zh-CN')
    })
  } finally {
    cleanupLoading.value = false
  }
}
</script>

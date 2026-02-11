<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-semibold">布局扩展管理</h3>
      <Button @click="handleSaveLayout" :disabled="saveLoading">
        {{ saveLoading ? '保存中...' : '保存更改' }}
      </Button>
    </div>

    <div v-if="loading" class="text-center py-8">
      <p class="text-gray-500">加载布局信息中...</p>
    </div>

    <div v-else-if="loadError" class="text-red-600 bg-red-50 p-4 rounded-md">
      {{ loadError }}
    </div>

    <div v-else class="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>自定义 Header 代码</CardTitle>
          <CardDescription>
            添加的代码将插入到页面 &lt;head&gt; 标签末尾。可用于添加自定义 CSS 样式或 Meta 标签。
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Textarea
            v-model="layout.header"
            class="font-mono text-sm min-h-[200px]"
            placeholder="例如: <style> body { background: #f0f0f0; } </style>"
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>自定义 Footer 代码</CardTitle>
          <CardDescription>
            添加的代码将插入到页面 &lt;body&gt; 标签末尾。可用于添加统计代码、自定义脚本等。
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Textarea
            v-model="layout.footer"
            class="font-mono text-sm min-h-[200px]"
            placeholder="例如: <script> console.log('BingPaper Loaded'); </script>"
          />
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { apiService } from '@/lib/api-service'
import type { LayoutResponse } from '@/lib/api-types'

const loading = ref(false)
const loadError = ref('')
const saveLoading = ref(false)

const layout = ref<LayoutResponse>({
  header: '',
  footer: ''
})

const fetchLayout = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const data = await apiService.getAdminLayout()
    layout.value = data
  } catch (err: any) {
    loadError.value = err.message || '获取布局信息失败'
    console.error('获取布局信息失败:', err)
  } finally {
    loading.value = false
  }
}

const handleSaveLayout = async () => {
  saveLoading.value = true
  try {
    await apiService.updateLayout(layout.value)
    toast.success('布局已保存')
  } catch (err: any) {
    toast.error(err.message || '保存失败')
    console.error('保存布局失败:', err)
  } finally {
    saveLoading.value = false
  }
}

onMounted(() => {
  fetchLayout()
})
</script>

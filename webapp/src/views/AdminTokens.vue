<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-semibold">Token 管理</h3>
      <Button @click="showCreateDialog = true">
        <span>创建 Token</span>
      </Button>
    </div>

    <div v-if="loading" class="text-center py-8">
      <p class="text-gray-500">加载中...</p>
    </div>

    <div v-else-if="error" class="text-red-600 bg-red-50 p-4 rounded-md">
      {{ error }}
    </div>

    <div v-else-if="tokens.length === 0" class="text-center py-8 text-gray-500">
      暂无 Token，点击上方按钮创建
    </div>

    <div v-else class="border rounded-lg overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>名称</TableHead>
            <TableHead>Token</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>过期时间</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead class="text-right">操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="token in tokens" :key="token.id">
            <TableCell>{{ token.id }}</TableCell>
            <TableCell>{{ token.name }}</TableCell>
            <TableCell>
              <code class="text-xs bg-gray-100 px-2 py-1 rounded">
                {{ token.token.substring(0, 20) }}...
              </code>
            </TableCell>
            <TableCell>
              <Badge :variant="token.disabled ? 'destructive' : 'default'">
                {{ token.disabled ? '已禁用' : '启用' }}
              </Badge>
            </TableCell>
            <TableCell>{{ formatDate(token.expires_at) }}</TableCell>
            <TableCell>{{ formatDate(token.created_at) }}</TableCell>
            <TableCell class="text-right space-x-2">
              <Button
                size="sm"
                variant="outline"
                @click="toggleTokenStatus(token)"
              >
                {{ token.disabled ? '启用' : '禁用' }}
              </Button>
              <Button
                size="sm"
                variant="destructive"
                @click="handleDeleteToken(token)"
              >
                删除
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- 创建 Token 对话框 -->
    <Dialog v-model:open="showCreateDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>创建 Token</DialogTitle>
          <DialogDescription>
            创建新的 API Token 用于访问接口
          </DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleCreateToken" class="space-y-4">
          <div class="space-y-2">
            <Label for="name">名称</Label>
            <Input
              id="name"
              v-model="createForm.name"
              placeholder="输入 Token 名称"
              required
            />
          </div>
          <div class="space-y-2">
            <Label for="expires_in">过期时间</Label>
            <Input
              id="expires_in"
              v-model="createForm.expires_in"
              placeholder="例如: 168h (7天), 720h (30天)"
            />
            <p class="text-xs text-gray-500">留空表示永不过期</p>
          </div>
          <div v-if="createError" class="text-sm text-red-600">
            {{ createError }}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="showCreateDialog = false">
              取消
            </Button>
            <Button type="submit" :disabled="createLoading">
              {{ createLoading ? '创建中...' : '创建' }}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- 删除确认对话框 -->
    <AlertDialog v-model:open="showDeleteDialog">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除 Token "{{ deleteTarget?.name }}" 吗？此操作无法撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>取消</AlertDialogCancel>
          <AlertDialogAction @click="confirmDelete" class="bg-red-600 hover:bg-red-700">
            删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { apiService } from '@/lib/api-service'
import type { Token } from '@/lib/api-types'

const tokens = ref<Token[]>([])
const loading = ref(false)
const error = ref('')

const showCreateDialog = ref(false)
const createForm = ref({
  name: '',
  expires_in: ''
})
const createLoading = ref(false)
const createError = ref('')

const showDeleteDialog = ref(false)
const deleteTarget = ref<Token | null>(null)

const fetchTokens = async () => {
  loading.value = true
  error.value = ''
  try {
    tokens.value = await apiService.getTokens()
  } catch (err: any) {
    error.value = err.message || '获取 Token 列表失败'
    console.error('获取 Token 失败:', err)
  } finally {
    loading.value = false
  }
}

const handleCreateToken = async () => {
  createLoading.value = true
  createError.value = ''
  try {
    await apiService.createToken(createForm.value)
    showCreateDialog.value = false
    createForm.value = { name: '', expires_in: '' }
    toast.success('Token 创建成功')
    await fetchTokens()
  } catch (err: any) {
    createError.value = err.message || '创建 Token 失败'
    console.error('创建 Token 失败:', err)
  } finally {
    createLoading.value = false
  }
}

const toggleTokenStatus = async (token: Token) => {
  try {
    await apiService.updateToken(token.id, { disabled: !token.disabled })
    toast.success(`Token 已${token.disabled ? '启用' : '禁用'}`)
    await fetchTokens()
  } catch (err: any) {
    console.error('更新 Token 状态失败:', err)
    toast.error(err.message || '更新失败')
  }
}

const handleDeleteToken = (token: Token) => {
  deleteTarget.value = token
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  
  try {
    await apiService.deleteToken(deleteTarget.value.id)
    showDeleteDialog.value = false
    deleteTarget.value = null
    toast.success('Token 删除成功')
    await fetchTokens()
  } catch (err: any) {
    console.error('删除 Token 失败:', err)
    toast.error(err.message || '删除失败')
  }
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('zh-CN')
  } catch {
    return dateStr
  }
}

onMounted(() => {
  fetchTokens()
})
</script>

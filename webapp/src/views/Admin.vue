<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <header class="bg-white border-b">
      <div class="container mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <h1 class="text-xl font-bold">BingPaper 管理后台</h1>
          </div>
          <div class="flex items-center gap-4">
            <Button variant="outline" size="sm" @click="showPasswordDialog = true">
              修改密码
            </Button>
            <Button variant="destructive" size="sm" @click="handleLogout">
              退出登录
            </Button>
          </div>
        </div>
      </div>
    </header>

    <!-- 主内容区 -->
    <div class="container mx-auto px-4 py-6">
      <Tabs v-model="activeTab" class="space-y-4">
        <TabsList class="grid w-full grid-cols-3 lg:w-[400px]">
          <TabsTrigger value="tokens">Token 管理</TabsTrigger>
          <TabsTrigger value="tasks">定时任务</TabsTrigger>
          <TabsTrigger value="config">系统配置</TabsTrigger>
        </TabsList>

        <TabsContent value="tokens" class="space-y-4">
          <AdminTokens />
        </TabsContent>

        <TabsContent value="tasks" class="space-y-4">
          <AdminTasks />
        </TabsContent>

        <TabsContent value="config" class="space-y-4">
          <AdminConfig />
        </TabsContent>
      </Tabs>
    </div>

    <!-- 修改密码对话框 -->
    <Dialog v-model:open="showPasswordDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>修改管理员密码</DialogTitle>
          <DialogDescription>
            请输入旧密码和新密码
          </DialogDescription>
        </DialogHeader>
        <form @submit.prevent="handleChangePassword" class="space-y-4">
          <div class="space-y-2">
            <Label for="old-password">旧密码</Label>
            <Input
              id="old-password"
              v-model="passwordForm.oldPassword"
              type="password"
              required
            />
          </div>
          <div class="space-y-2">
            <Label for="new-password">新密码</Label>
            <Input
              id="new-password"
              v-model="passwordForm.newPassword"
              type="password"
              required
            />
          </div>
          <div class="space-y-2">
            <Label for="confirm-password">确认新密码</Label>
            <Input
              id="confirm-password"
              v-model="passwordForm.confirmPassword"
              type="password"
              required
            />
          </div>
          <div v-if="passwordError" class="text-sm text-red-600">
            {{ passwordError }}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" @click="showPasswordDialog = false">
              取消
            </Button>
            <Button type="submit" :disabled="passwordLoading">
              {{ passwordLoading ? '提交中...' : '确认修改' }}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { apiService } from '@/lib/api-service'
import { apiClient } from '@/lib/http-client'
import AdminTokens from './AdminTokens.vue'
import AdminTasks from './AdminTasks.vue'
import AdminConfig from './AdminConfig.vue'

const router = useRouter()
const activeTab = ref('tokens')

const showPasswordDialog = ref(false)
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const passwordLoading = ref(false)
const passwordError = ref('')

// 检查认证状态
const checkAuth = () => {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    router.push('/admin/login')
    return false
  }
  
  // 设置认证头
  apiClient.setAuthToken(token)
  
  // 检查是否过期
  const expiresAt = localStorage.getItem('admin_token_expires')
  if (expiresAt) {
    const expireDate = new Date(expiresAt)
    if (expireDate < new Date()) {
      toast.warning('登录已过期，请重新登录')
      handleLogout()
      return false
    }
  }
  
  return true
}

const handleLogout = () => {
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_token_expires')
  apiClient.clearAuthToken()
  router.push('/admin/login')
}

const handleChangePassword = async () => {
  passwordError.value = ''
  
  // 验证新密码
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    passwordError.value = '两次输入的新密码不一致'
    return
  }
  
  if (passwordForm.value.newPassword.length < 6) {
    passwordError.value = '新密码长度至少为 6 位'
    return
  }
  
  passwordLoading.value = true
  
  try {
    await apiService.changePassword({
      old_password: passwordForm.value.oldPassword,
      new_password: passwordForm.value.newPassword
    })
    
    toast.success('密码修改成功，请重新登录')
    showPasswordDialog.value = false
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    handleLogout()
  } catch (err: any) {
    passwordError.value = err.message || '密码修改失败'
    console.error('修改密码失败:', err)
  } finally {
    passwordLoading.value = false
  }
}

onMounted(() => {
  checkAuth()
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
    <Card class="w-full max-w-md">
      <CardHeader class="space-y-1">
        <CardTitle class="text-2xl font-bold text-center">管理员登录</CardTitle>
        <CardDescription class="text-center">
          输入管理员密码以访问后台管理系统
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="space-y-2">
            <Label for="password">密码</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              placeholder="请输入管理员密码"
              required
              :disabled="loading"
            />
          </div>
          
          <div v-if="error" class="text-sm text-red-600 bg-red-50 p-3 rounded-md">
            {{ error }}
          </div>

          <Button type="submit" class="w-full" :disabled="loading">
            <span v-if="loading">登录中...</span>
            <span v-else>登录</span>
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { apiService } from '@/lib/api-service'
import { apiClient } from '@/lib/http-client'

const router = useRouter()
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    const response = await apiService.login({ password: password.value })
    
    // 保存 token 到 localStorage
    localStorage.setItem('admin_token', response.token)
    localStorage.setItem('admin_token_expires', response.expires_at || '')
    
    // 设置 HTTP 客户端的认证头
    apiClient.setAuthToken(response.token)
    
    toast.success('登录成功')
    
    // 跳转到管理后台
    router.push('/admin')
  } catch (err: any) {
    console.error('登录失败:', err)
    error.value = err.message || '登录失败，请检查密码'
  } finally {
    loading.value = false
  }
}
</script>

import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import ImageView from '@/views/ImageView.vue'
import ApiDocs from '@/views/ApiDocs.vue'
import AdminLogin from '@/views/AdminLogin.vue'
import Admin from '@/views/Admin.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home,
      meta: {
        title: '必应每日一图'
      }
    },
    {
      path: '/image/:date?',
      name: 'ImageView',
      component: ImageView,
      meta: {
        title: '图片详情'
      },
      beforeEnter: (to, _from, next) => {
        // 如果没有提供日期参数，重定向到今天的日期
        if (!to.params.date) {
          const today = new Date().toISOString().split('T')[0]
          next({ path: `/image/${today}`, replace: true })
        } else {
          next()
        }
      }
    },
    {
      path: '/api-docs',
      name: 'ApiDocs',
      component: ApiDocs,
      meta: {
        title: 'API 文档'
      }
    },
    {
      path: '/admin/login',
      name: 'AdminLogin',
      component: AdminLogin,
      meta: {
        title: '管理员登录'
      }
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      meta: {
        title: '管理后台',
        requiresAuth: true
      }
    }
  ]
})

// 路由守卫 - 更新页面标题和认证检查
router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || '必应每日一图'
  
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('admin_token')
    if (!token) {
      // 未登录，重定向到登录页
      next('/admin/login')
      return
    }
    
    // 检查 token 是否过期
    const expiresAt = localStorage.getItem('admin_token_expires')
    if (expiresAt) {
      const expireDate = new Date(expiresAt)
      if (expireDate < new Date()) {
        // token 已过期，清除并重定向到登录页
        localStorage.removeItem('admin_token')
        localStorage.removeItem('admin_token_expires')
        next('/admin/login')
        return
      }
    }
  }
  
  next()
})

export default router

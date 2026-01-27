import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import ImageView from '@/views/ImageView.vue'

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
      path: '/image/:date',
      name: 'ImageView',
      component: ImageView,
      meta: {
        title: '图片详情'
      }
    }
  ]
})

// 路由守卫 - 更新页面标题
router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || '必应每日一图'
  next()
})

export default router

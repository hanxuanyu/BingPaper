<script setup lang="ts">
  import { onMounted } from 'vue'
  import 'vue-sonner/style.css'
  import { Toaster } from '@/components/ui/sonner'
  import { apiService } from '@/lib/api-service'

  onMounted(async () => {
    try {
      const layout = await apiService.getLayout()
      
      const injectHTML = (html: string, target: HTMLElement) => {
        if (!html) return
        const range = document.createRange()
        range.selectNode(target)
        const fragment = range.createContextualFragment(html)
        target.appendChild(fragment)
      }

      injectHTML(layout.header, document.head)
      injectHTML(layout.footer, document.body)
    } catch (err) {
      console.error('Failed to load custom layout:', err)
    }
  })
</script>

<template>
  <div id="app">
    <RouterView />
    <Toaster />
  </div>
</template>

<style>
#app {
  font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

html, body {
  margin: 0;
  padding: 0;
  overflow-x: hidden;
}
</style>

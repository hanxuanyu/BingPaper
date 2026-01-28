<script setup lang="ts">
  import 'vue-sonner/style.css'
  import { Toaster } from '@/components/ui/sonner'
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

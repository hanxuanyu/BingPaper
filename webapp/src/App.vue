<script setup lang="ts">
  import 'vue-sonner/style.css'
</script>

<template>
  <div id="app">
    <RouterView />
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

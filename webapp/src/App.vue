<script setup lang="ts">
  import 'vue-sonner/style.css'
  import ComponentShowcase from './views/ComponentShowcase.vue'
</script>

<template>
  <div id="app">
    <ComponentShowcase />
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
}
</style>

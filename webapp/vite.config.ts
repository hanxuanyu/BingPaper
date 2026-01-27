import path from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [vue(), tailwindcss()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    // 构建配置
    build: {
      // 输出到上级目录的web文件夹
      outDir: path.resolve(__dirname, '../web'),
      // 清空输出目录
      emptyOutDir: true,
      // 静态资源处理
      assetsDir: 'assets',
      // 生成 sourcemap 用于生产环境调试
      sourcemap: mode === 'development',
      // 静态资源内联阈值（小于此大小的资源会被内联为 base64）
      assetsInlineLimit: 4096,
      // chunk 分割策略
      rollupOptions: {
        output: {
          // 静态资源文件名（包含内容哈希，利于长期缓存）
          assetFileNames: 'assets/[name]-[hash][extname]',
          // JS chunk 文件名
          chunkFileNames: 'assets/[name]-[hash].js',
          // 入口文件名
          entryFileNames: 'assets/[name]-[hash].js',
          // 手动分割代码
          manualChunks: {
            // 将 Vue 相关代码单独打包
            'vue-vendor': ['vue', 'vue-router'],
            // 将 UI 组件库单独打包（如果有的话）
            // 'ui-vendor': ['其他UI库']
          }
        }
      }
    },
    // 开发服务器配置
    server: {
      port: 5173,
      strictPort: false,
      open: false
    },
    // 环境变量配置 - 开发环境使用完整URL，生产环境使用相对路径
    define: {
      __API_BASE_URL__: JSON.stringify(env.VITE_API_BASE_URL || '/api/v1')
    }
  }
})

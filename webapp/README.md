# BingPaper WebApp

BingPaper 的前端 Web 应用，使用 Vue 3 + TypeScript + Vite 构建。

## 特性

- ✨ Vue 3 组合式 API
- 🎨 Tailwind CSS + shadcn-vue 组件库
- 📦 TypeScript 类型支持
- 🔧 完整的 API 客户端封装
- 🚀 优化的构建配置
- 🌐 支持自定义后端路径
- 📁 自动输出到上级目录的 web 文件夹

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发环境

```bash
npm run dev
```

开发服务器会在 `http://localhost:5173` 启动，并自动代理 `/api` 请求到后端服务器。

### 构建生产版本

```bash
npm run build
```

构建产物会自动输出到项目上级目录的 `web` 文件夹，供 Go 服务器使用。

### 预览构建结果

```bash
npm run preview
```

## 构建配置

### 环境变量

项目支持通过环境变量配置后端 API 路径：

- `.env` - 默认配置
- `.env.development` - 开发环境配置
- `.env.production` - 生产环境配置

修改 `VITE_API_BASE_URL` 来自定义后端 API 路径。

### 输出目录

构建产物输出到 `../web/`，目录结构：

```
web/
├── index.html
├── assets/
│   ├── index-[hash].js
│   └── index-[hash].css
└── vite.svg
```

## API 使用

项目提供了完整的 TypeScript API 客户端：

```typescript
import { bingPaperApi } from '@/lib/api-service'

// 获取今日图片
const meta = await bingPaperApi.getTodayImageMeta()

// 获取图片列表
const images = await bingPaperApi.getImages({ limit: 10 })
```

详细的 API 使用示例请参阅 [API_EXAMPLES.md](./API_EXAMPLES.md)

## 构建说明

详细的构建配置和部署说明请参阅 [BUILD.md](./BUILD.md)

## 项目结构

```
src/
├── lib/              # 核心库
│   ├── api-config.ts    # API 配置
│   ├── api-types.ts     # TypeScript 类型定义
│   ├── api-service.ts   # API 服务封装
│   ├── http-client.ts   # HTTP 客户端
│   └── utils.ts        # 工具函数
├── components/       # Vue 组件
│   └── ui/          # UI 组件库
├── views/           # 页面视图
├── assets/          # 静态资源
├── App.vue          # 根组件
└── main.ts          # 入口文件
```

## 技术栈

- [Vue 3](https://vuejs.org/) - 渐进式 JavaScript 框架
- [TypeScript](https://www.typescriptlang.org/) - 类型安全的 JavaScript
- [Vite](https://vitejs.dev/) - 下一代前端构建工具
- [Tailwind CSS](https://tailwindcss.com/) - 实用优先的 CSS 框架
- [shadcn-vue](https://www.shadcn-vue.com/) - 高质量的 Vue 组件库

## IDE 支持

推荐使用 [VS Code](https://code.visualstudio.com/) + [Vue - Official](https://marketplace.visualstudio.com/items?itemName=Vue.volar) 扩展。

## License

MIT

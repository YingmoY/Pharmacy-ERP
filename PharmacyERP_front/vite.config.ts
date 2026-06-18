import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    // 按需引入 Ant Design Vue 组件
    Components({
      resolvers: [
        AntDesignVueResolver({
          importStyle: false,
        }),
      ],
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    allowedHosts: ['localhost', '127.0.0.1', 'pharmacy.yingmoy.com'],
    proxy: {
      // 代理 ERP 主系统 API
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        proxyTimeout: 660000,
        timeout: 660000,
      },
      // 代理 AI 子模块 API（仅匹配 /ai/api/，路径原样转发给 Python 服务）
      // Python FastAPI 路由前缀本身就是 /ai/api/v1，无需 rewrite
      '/ai/api/': {
        target: 'http://localhost:9080',
        changeOrigin: true,
        ws: true,
      },
    },
  },
})

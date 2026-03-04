import path from 'node:path'
import { defineConfig, type Plugin, type ViteDevServer } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'
import tsconfigPaths from 'vite-tsconfig-paths'

const coiHeaders: Plugin = {
  name: 'coi-headers',
  configureServer(server: ViteDevServer) {
    server.middlewares.use((req, res, next) => {
      res.setHeader('Cross-Origin-Opener-Policy', 'same-origin')
      res.setHeader('Cross-Origin-Embedder-Policy', 'require-corp')
      res.setHeader('Cross-Origin-Resource-Policy', 'cross-origin') // für CDN okay
      next()
    })
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(), 
    tsconfigPaths(),
    coiHeaders
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  // Required for SharedArrayBuffer used by @ffmpeg/core-mt (multi-threaded WebAssembly).
  // For production, configure your web server (Nginx/Caddy) with the same headers:
  //   Cross-Origin-Opener-Policy: same-origin
  //   Cross-Origin-Embedder-Policy: require-corp
  server: {
    headers: {
      'Cross-Origin-Opener-Policy': 'same-origin',
      'Cross-Origin-Embedder-Policy': 'require-corp',
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'),
      },
    },
  },
  preview: {
    headers: {
      'Cross-Origin-Opener-Policy': 'same-origin',
      'Cross-Origin-Embedder-Policy': 'require-corp',
    },
  },
  optimizeDeps: {
    // FFmpeg WASM modules must be served as raw ES modules without Vite
    // pre-bundling; they use dynamic imports + SharedArrayBuffer internally.
    exclude: ['@ffmpeg/ffmpeg', '@ffmpeg/util'],
  },
  build: {
    // esnext is required for top-level await and BigInt64Array used by the
    // FFmpeg WASM core. All modern browsers (Chrome 89+, Firefox 89+,
    // Safari 15+) support this target.
    target: 'esnext',
  },
})

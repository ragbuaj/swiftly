import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  optimizeDeps: {
    include: [
      'reka-ui',
      'lucide-vue-next',
      'clsx',
      'tailwind-merge',
      'axios',
      'pinia'
    ]
  },
  server: {
    host: '0.0.0.0', // Allow access from outside the container
    port: 5173,
    headers: {
      'Cross-Origin-Opener-Policy': 'same-origin-allow-popups',
      'Cross-Origin-Embedder-Policy': 'unsafe-none',
    },
    watch: {
      usePolling: true, // Needed for Windows/WSL2 file systems
    },
  },
})

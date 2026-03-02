import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // Allow access from outside the container
    port: 5173,
    watch: {
      usePolling: true, // Needed for Windows/WSL2 file systems
    },
  },
})

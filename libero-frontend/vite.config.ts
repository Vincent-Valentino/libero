import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'; // Use node:url for path resolution

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: { // Add server configuration
    proxy: {
      // Proxy requests starting with /api to the backend server
      '/api': {
        target: 'http://localhost:8080', // Your Go backend address
        changeOrigin: true, // Needed for virtual hosted sites
        // secure: false, // Uncomment if backend uses self-signed SSL cert
        // rewrite: (path) => path.replace(/^\/api/, ''), // Uncomment if backend doesn't expect /api prefix
      },
    },
  },
  resolve: {
    alias: {
      // Use import.meta.url and URL constructor for ES Modules compatibility
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
})

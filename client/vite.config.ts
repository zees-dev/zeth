import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 3000,
  },
  base: './',
  resolve: {
    alias: {
      '@': __dirname + '/src',
    }
  },
  build: { // dev build
    minify: false
  },
})

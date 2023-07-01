import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    minify: false,
    outDir: '../public',
    emptyOutDir: true,
    rollupOptions: {
      input: {
        main: './index.html',
        support: './support.html',
      },
    },
  },
});

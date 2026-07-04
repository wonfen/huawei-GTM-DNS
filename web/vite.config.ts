import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  ssr: {
    noExternal: ['@tanstack/query-core'],
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
      '/help.html': 'http://localhost:8080',
    },
  },
});

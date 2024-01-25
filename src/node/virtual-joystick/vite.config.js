import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import sveltePreprocess from 'svelte-preprocess';

export default defineConfig({
    plugins: [svelte({
        preprocess: sveltePreprocess({
            typescript: true,
        }),
    })],
    server: {
        host: '0.0.0.0' // Listen on all network interfaces
    }
});

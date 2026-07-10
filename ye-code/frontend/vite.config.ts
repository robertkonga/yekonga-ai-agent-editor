import path from 'path'
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';
import wails from "@wailsio/runtime/plugins/vite";

// https://vite.dev/config/
export default defineConfig({
    
    server: {
        host: "127.0.0.1",
        port: Number(process.env.WAILS_VITE_PORT) || 9245,
        strictPort: true,
    },
    plugins: [
        vue(),
        tailwindcss(),
        wails("./bindings"),
    ],
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.vue', '.json'],
        alias: {
            '@': path.resolve(__dirname, 'src'),
            '@wails': path.resolve(__dirname, "bindings"),
        },
    },
    build: {
        emptyOutDir: true,  // already default, but forces it explicitly
    }
})

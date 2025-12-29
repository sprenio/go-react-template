import {defineConfig, loadEnv} from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

const envPath = path.resolve(__dirname, '..')

export default defineConfig(({mode}) => {
    const env = loadEnv(mode, envPath)
    return {
        plugins: [
            react(),
            tailwindcss(),
        ],
        server: {
            host: '0.0.0.0',
            port: Number(env.VITE_PORT) || 5173,
        },
        envDir: envPath,
        resolve: {
            alias: {
                '@': path.resolve(__dirname, './src'),
            }
        },
    }
})

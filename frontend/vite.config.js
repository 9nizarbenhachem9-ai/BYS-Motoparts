import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

// Konfigurasi Vite untuk frontend vanilla (tanpa framework)
export default defineConfig({
  // Aplikasi multi-halaman: index, login, dan register
  build: {
    rollupOptions: {
      input: {
        main: fileURLToPath(new URL("./index.html", import.meta.url)),
        login: fileURLToPath(new URL("./login.html", import.meta.url)),
        register: fileURLToPath(new URL("./register.html", import.meta.url)),
      },
    },
  },

  server: {
    port: 5173,
    // Mode development: teruskan panggilan /api ke backend Go (localhost:8080)
    // sehingga browser selalu same-origin dan tidak masalah CORS.
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});

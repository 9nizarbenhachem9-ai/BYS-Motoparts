package main

import (
	"log"

	"sparepart-backend/config"
	"sparepart-backend/database"
	"sparepart-backend/handlers"
	"sparepart-backend/middleware"
	"sparepart-backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[INFO] Memulai aplikasi Web Manajemen Sparepart...")

	// 1. Memuat konfigurasi dari environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[FATAL] Gagal memuat konfigurasi: %v", err)
	}
	log.Printf("[INFO] Konfigurasi berhasil dimuat (Port: %s, Origin: %s)", cfg.AppPort, cfg.AllowedOrigin)

	// 2. Menguji koneksi database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[FATAL] Koneksi database gagal: %v", err)
	}
	log.Println("[INFO] Berhasil terhubung ke database PostgreSQL!")

	// 3. Menjalankan migrasi skema tabel
	if err := database.Migrate(db); err != nil {
		log.Fatalf("[FATAL] Migrasi database gagal: %v", err)
	}
	log.Println("[INFO] Migrasi database berhasil: tabel users, spareparts, dan stock_transactions siap.")

	// 4. Inisialisasi Service dan Handler
	authService := services.NewAuthService(db, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	// 5. Inisialisasi Gin Engine & Middleware
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigin))

	// Rute API
	api := router.Group("/api")
	{
		// Rute Publik (Autentikasi)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Rute Terproteksi (Wajib Bearer Token JWT)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/auth/me", authHandler.Me)
		}
	}

	// 6. Menjalankan HTTP Server
	log.Printf("[INFO] Server web berjalan pada port %s (http://localhost:%s)...", cfg.AppPort, cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan HTTP server: %v", err)
	}
}

package main

import (
	"log"

	"sparepart-backend/config"
	"sparepart-backend/database"
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
}

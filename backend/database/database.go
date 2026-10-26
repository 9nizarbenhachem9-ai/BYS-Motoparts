package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah instance global GORM yang digunakan oleh seluruh aplikasi
var DB *gorm.DB

// Connect menginisialisasi koneksi ke PostgreSQL menggunakan connection string dari DATABASE_URL.
func Connect(databaseURL string) (*gorm.DB, error) {
	// Konfigurasi koneksi GORM
	// Menggunakan logger default level Silent atau Warn agar kredensial tidak bocor ke log
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(postgres.Open(databaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	// Validasi koneksi dasar dengan ping ke database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan objek generic database: %w", err)
	}

	// Atur connection pool sederhana
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("gagal menghubungi PostgreSQL (ping): %w", err)
	}

	DB = db
	return DB, nil
}

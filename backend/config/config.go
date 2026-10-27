package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menampung semua konfigurasi aplikasi yang dibaca dari environment variables.
type Config struct {
	AppPort            string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiresHours    int
	AllowedOrigin      string
	RegisterInviteCode string
}

// LoadConfig memuat variabel lingkungan dari file .env (jika ada) dan dari sistem OS.
// Jika variabel wajib tidak ada, fungsi mengembalikan error dengan menyebut nama variabel yang hilang.
func LoadConfig() (*Config, error) {
	// Coba muat file .env dari folder saat ini atau folder root proyek
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = os.Getenv("PORT") // Mendukung environment PaaS (Hostless, Railway, Render, dll)
	}
	if appPort == "" {
		appPort = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("konfigurasi wajib hilang: DATABASE_URL belum diatur")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("konfigurasi wajib hilang: JWT_SECRET belum diatur")
	}

	jwtExpiresHoursStr := os.Getenv("JWT_EXPIRES_HOURS")
	jwtExpiresHours := 24
	if jwtExpiresHoursStr != "" {
		parsedHours, err := strconv.Atoi(jwtExpiresHoursStr)
		if err != nil || parsedHours <= 0 {
			return nil, fmt.Errorf("nilai JWT_EXPIRES_HOURS tidak valid: harus bilangan bulat positif")
		}
		jwtExpiresHours = parsedHours
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5500"
	}

	registerInviteCode := os.Getenv("REGISTER_INVITE_CODE")

	return &Config{
		AppPort:            appPort,
		DatabaseURL:        databaseURL,
		JWTSecret:          jwtSecret,
		JWTExpiresHours:    jwtExpiresHours,
		AllowedOrigin:      allowedOrigin,
		RegisterInviteCode: registerInviteCode,
	}, nil
}


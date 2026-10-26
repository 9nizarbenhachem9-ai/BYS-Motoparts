package config_test

import (
	"os"
	"strings"
	"testing"

	"sparepart-backend/config"
)

func TestLoadConfig_MissingDatabaseURL(t *testing.T) {
	os.Clearenv()

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("Diharapkan error karena DATABASE_URL tidak ada, tetapi tidak terjadi error")
	}

	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("Pesan error harus menyebut DATABASE_URL, didapat: %v", err)
	}
}

func TestLoadConfig_MissingJWTSecret(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("Diharapkan error karena JWT_SECRET tidak ada, tetapi tidak terjadi error")
	}

	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Errorf("Pesan error harus menyebut JWT_SECRET, didapat: %v", err)
	}
}

func TestLoadConfig_Success(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("JWT_SECRET", "supersecret123")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	if cfg.AppPort != "8080" {
		t.Errorf("Port default harus 8080, didapat: %s", cfg.AppPort)
	}
	if cfg.JWTExpiresHours != 24 {
		t.Errorf("JWTExpiresHours default harus 24, didapat: %d", cfg.JWTExpiresHours)
	}
	if cfg.AllowedOrigin != "http://localhost:5500" {
		t.Errorf("AllowedOrigin default harus http://localhost:5500, didapat: %s", cfg.AllowedOrigin)
	}
}

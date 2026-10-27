package services_test

import (
	"testing"

	"sparepart-backend/config"
	"sparepart-backend/models"
	"sparepart-backend/services"
)

func TestAuthService_PasswordHashAndVerify(t *testing.T) {
	svc := services.NewAuthService(nil, nil)

	password := "rahasia123"
	hash, err := svc.HashPassword(password)
	if err != nil {
		t.Fatalf("Gagal melakukan hash password: %v", err)
	}

	if !svc.CheckPassword(hash, password) {
		t.Errorf("Verifikasi password yang benar harus bernilai true")
	}

	if svc.CheckPassword(hash, "passwordsalah") {
		t.Errorf("Verifikasi password yang salah harus bernilai false")
	}
}

func TestAuthService_JWTGenerateAndValidate(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:       "kunci-rahasia-test-1234567890",
		JWTExpiresHours: 1,
	}
	svc := services.NewAuthService(nil, cfg)

	user := &models.User{
		ID:       10,
		Username: "admin_test",
	}

	token, err := svc.GenerateToken(user)
	if err != nil {
		t.Fatalf("Gagal menghasilkan token JWT: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("Gagal memvalidasi token JWT valid: %v", err)
	}

	if claims.UserID != 10 {
		t.Errorf("UserID harus 10, didapat: %d", claims.UserID)
	}
	if claims.Username != "admin_test" {
		t.Errorf("Username harus 'admin_test', didapat: %s", claims.Username)
	}

	// Uji coba token dengan secret yang salah
	svcWrongKey := services.NewAuthService(nil, &config.Config{
		JWTSecret: "kunci-salah-999999",
	})
	_, err = svcWrongKey.ValidateToken(token)
	if err == nil {
		t.Errorf("Diharapkan error saat memvalidasi token dengan secret berbeda")
	}
}

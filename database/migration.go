package database

import (
	"fmt"

	"sparepart-backend/models"

	"gorm.io/gorm"
)

// Migrate menjalankan migrasi otomatis skema tabel menggunakan GORM AutoMigrate
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Sparepart{},
		&models.StockTransaction{},
	)
	if err != nil {
		return fmt.Errorf("migrasi database gagal: %w", err)
	}

	return nil
}


package services

import (
	"errors"
	"fmt"
	"strings"

	"sparepart-backend/models"

	"gorm.io/gorm"
)

// TransactionService mengelola pencatatan barang masuk/keluar dan riwayat transaksi
type TransactionService struct {
	db *gorm.DB
}

// NewTransactionService membuat instance baru TransactionService
func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

// RecordTransactionRequest payload pencatatan barang masuk / keluar
type RecordTransactionRequest struct {
	Tipe    string `json:"tipe"` // 'masuk' atau 'keluar'
	Jumlah  int    `json:"jumlah"`
	Catatan string `json:"catatan"`
}

// RecordTransaction mencatat transaksi stok dalam satu database transaction atomic
func (s *TransactionService) RecordTransaction(userID uint, sparepartID uint, req RecordTransactionRequest) (*models.StockTransaction, *models.Sparepart, error) {
	tipe := strings.ToLower(strings.TrimSpace(req.Tipe))
	if tipe != "masuk" && tipe != "keluar" {
		return nil, nil, errors.New("tipe transaksi harus 'masuk' atau 'keluar'")
	}

	if req.Jumlah <= 0 {
		return nil, nil, errors.New("jumlah transaksi harus lebih besar dari 0")
	}

	var transaction models.StockTransaction
	var sparepart models.Sparepart

	// Jalankan dalam satu Database Transaction atomic
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Ambil data sparepart dan kunci baris untuk update (FOR UPDATE)
		if err := tx.First(&sparepart, sparepartID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("sparepart tidak ditemukan")
			}
			return errors.New("gagal memuat data sparepart")
		}

		// 2. Logika stok berdasarkan tipe transaksi
		if tipe == "keluar" {
			if sparepart.Stok < req.Jumlah {
				return fmt.Errorf("stok tidak mencukupi (stok saat ini: %d, diminta: %d)", sparepart.Stok, req.Jumlah)
			}
			sparepart.Stok -= req.Jumlah
		} else {
			sparepart.Stok += req.Jumlah
		}

		// 3. Simpan perubahan stok pada sparepart
		if err := tx.Model(&sparepart).Update("stok", sparepart.Stok).Error; err != nil {
			return errors.New("gagal memperbarui stok sparepart")
		}

		// 4. Catat ke tabel riwayat transaksi
		transaction = models.StockTransaction{
			SparepartID: sparepartID,
			Tipe:        tipe,
			Jumlah:      req.Jumlah,
			Catatan:     strings.TrimSpace(req.Catatan),
			UserID:      userID,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return errors.New("gagal mencatat transaksi stok")
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &transaction, &sparepart, nil
}

// GetTransactionsBySparepartID mengambil daftar riwayat transaksi untuk suatu sparepart
func (s *TransactionService) GetTransactionsBySparepartID(sparepartID uint) ([]models.StockTransaction, error) {
	// Pastikan sparepart ada
	var count int64
	s.db.Model(&models.Sparepart{}).Where("id = ?", sparepartID).Count(&count)
	if count == 0 {
		return nil, errors.New("sparepart tidak ditemukan")
	}

	var history []models.StockTransaction
	err := s.db.Where("sparepart_id = ?", sparepartID).
		Preload("User").
		Order("created_at DESC").
		Find(&history).Error

	if err != nil {
		return nil, errors.New("gagal mengambil riwayat transaksi")
	}

	return history, nil
}


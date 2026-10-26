package services_test

import (
	"testing"

	"sparepart-backend/services"
)

func TestTransactionService_Validation(t *testing.T) {
	svc := services.NewTransactionService(nil)

	// Uji tipe transaksi tidak valid
	_, _, err := svc.RecordTransaction(1, 1, services.RecordTransactionRequest{
		Tipe:   "rusak",
		Jumlah: 10,
	})
	if err == nil || err.Error() != "tipe transaksi harus 'masuk' atau 'keluar'" {
		t.Errorf("Diharapkan error 'tipe transaksi harus masuk atau keluar', didapat: %v", err)
	}

	// Uji jumlah transaksi 0
	_, _, err = svc.RecordTransaction(1, 1, services.RecordTransactionRequest{
		Tipe:   "masuk",
		Jumlah: 0,
	})
	if err == nil || err.Error() != "jumlah transaksi harus lebih besar dari 0" {
		t.Errorf("Diharapkan error 'jumlah transaksi harus lebih besar dari 0', didapat: %v", err)
	}

	// Uji jumlah transaksi negatif
	_, _, err = svc.RecordTransaction(1, 1, services.RecordTransactionRequest{
		Tipe:   "keluar",
		Jumlah: -5,
	})
	if err == nil || err.Error() != "jumlah transaksi harus lebih besar dari 0" {
		t.Errorf("Diharapkan error 'jumlah transaksi harus lebih besar dari 0', didapat: %v", err)
	}
}


package services_test

import (
	"testing"

	"sparepart-backend/services"
)

func TestSparepartService_Validation(t *testing.T) {
	svc := services.NewSparepartService(nil)

	// Uji validasi kode kosong
	_, err := svc.Create(services.CreateSparepartRequest{
		Kode: "",
		Nama: "Busi Racing",
	})
	if err == nil || err.Error() != "kode sparepart wajib diisi" {
		t.Errorf("Diharapkan error 'kode sparepart wajib diisi', didapat: %v", err)
	}

	// Uji validasi nama kosong
	_, err = svc.Create(services.CreateSparepartRequest{
		Kode: "SP-001",
		Nama: "",
	})
	if err == nil || err.Error() != "nama sparepart wajib diisi" {
		t.Errorf("Diharapkan error 'nama sparepart wajib diisi', didapat: %v", err)
	}

	// Uji validasi harga negatif
	_, err = svc.Create(services.CreateSparepartRequest{
		Kode:  "SP-001",
		Nama:  "Busi Racing",
		Harga: -1000,
	})
	if err == nil || err.Error() != "harga tidak boleh negatif" {
		t.Errorf("Diharapkan error 'harga tidak boleh negatif', didapat: %v", err)
	}

	// Uji validasi stok minimum negatif
	_, err = svc.Create(services.CreateSparepartRequest{
		Kode:        "SP-001",
		Nama:        "Busi Racing",
		Harga:       15000,
		StokMinimum: -5,
	})
	if err == nil || err.Error() != "stok minimum tidak boleh negatif" {
		t.Errorf("Diharapkan error 'stok minimum tidak boleh negatif', didapat: %v", err)
	}
}


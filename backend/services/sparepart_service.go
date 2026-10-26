package services

import (
	"errors"
	"strings"

	"sparepart-backend/models"

	"gorm.io/gorm"
)

// SparepartService mengelola logika bisnis untuk CRUD sparepart
type SparepartService struct {
	db *gorm.DB
}

// NewSparepartService membuat instance baru dari SparepartService
func NewSparepartService(db *gorm.DB) *SparepartService {
	return &SparepartService{db: db}
}

// CreateSparepartRequest payload penambahan sparepart baru
type CreateSparepartRequest struct {
	Kode        string `json:"kode"`
	Nama        string `json:"nama"`
	Kategori    string `json:"kategori"`
	Merek       string `json:"merek"`
	Harga       int64  `json:"harga"`
	StokMinimum int    `json:"stok_minimum"`
	LokasiRak   string `json:"lokasi_rak"`
}

// UpdateSparepartRequest payload pengubahan data sparepart (tanpa mengubah stok langsung)
type UpdateSparepartRequest struct {
	Kode        string `json:"kode"`
	Nama        string `json:"nama"`
	Kategori    string `json:"kategori"`
	Merek       string `json:"merek"`
	Harga       int64  `json:"harga"`
	StokMinimum int    `json:"stok_minimum"`
	LokasiRak   string `json:"lokasi_rak"`
}

// PaginationResult menampung data daftar sparepart dan metadata paginasi
type PaginationResult struct {
	Items      []models.Sparepart `json:"items"`
	TotalItems int64              `json:"total_items"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// Create menambahkan sparepart baru
func (s *SparepartService) Create(req CreateSparepartRequest) (*models.Sparepart, error) {
	kode := strings.TrimSpace(req.Kode)
	nama := strings.TrimSpace(req.Nama)

	if kode == "" {
		return nil, errors.New("kode sparepart wajib diisi")
	}
	if nama == "" {
		return nil, errors.New("nama sparepart wajib diisi")
	}
	if req.Harga < 0 {
		return nil, errors.New("harga tidak boleh negatif")
	}
	if req.StokMinimum < 0 {
		return nil, errors.New("stok minimum tidak boleh negatif")
	}

	// Cek keunikan kode sparepart (abaikan yang sudah soft deleted)
	var count int64
	s.db.Model(&models.Sparepart{}).Where("kode = ?", kode).Count(&count)
	if count > 0 {
		return nil, errors.New("kode sparepart sudah digunakan")
	}

	item := models.Sparepart{
		Kode:        kode,
		Nama:        nama,
		Kategori:    strings.TrimSpace(req.Kategori),
		Merek:       strings.TrimSpace(req.Merek),
		Harga:       req.Harga,
		Stok:        0, // Stok awal selalu 0, hanya berubah lewat transaksi barang masuk
		StokMinimum: req.StokMinimum,
		LokasiRak:   strings.TrimSpace(req.LokasiRak),
	}

	if err := s.db.Create(&item).Error; err != nil {
		return nil, errors.New("gagal menyimpan data sparepart")
	}

	return &item, nil
}

// GetAll mengambil daftar sparepart dengan pencarian, filter kategori, dan paginasi
func (s *SparepartService) GetAll(search, kategori string, page, limit int) (*PaginationResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	query := s.db.Model(&models.Sparepart{})

	search = strings.TrimSpace(search)
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(nama) LIKE ? OR LOWER(kode) LIKE ? OR LOWER(merek) LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	kategori = strings.TrimSpace(kategori)
	if kategori != "" {
		query = query.Where("LOWER(kategori) = ?", strings.ToLower(kategori))
	}

	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, errors.New("gagal menghitung total sparepart")
	}

	var items []models.Sparepart
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, errors.New("gagal mengambil daftar sparepart")
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationResult{
		Items:      items,
		TotalItems: totalItems,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetByID mengambil detail satu sparepart berdasarkan ID
func (s *SparepartService) GetByID(id uint) (*models.Sparepart, error) {
	var item models.Sparepart
	if err := s.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sparepart tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil data sparepart")
	}
	return &item, nil
}

// Update memperbarui data sparepart tanpa mengubah stok langsung
func (s *SparepartService) Update(id uint, req UpdateSparepartRequest) (*models.Sparepart, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	kode := strings.TrimSpace(req.Kode)
	nama := strings.TrimSpace(req.Nama)

	if kode == "" {
		return nil, errors.New("kode sparepart wajib diisi")
	}
	if nama == "" {
		return nil, errors.New("nama sparepart wajib diisi")
	}
	if req.Harga < 0 {
		return nil, errors.New("harga tidak boleh negatif")
	}
	if req.StokMinimum < 0 {
		return nil, errors.New("stok minimum tidak boleh negatif")
	}

	// Cek apakah kode sudah dipakai oleh sparepart lain
	var count int64
	s.db.Model(&models.Sparepart{}).Where("kode = ? AND id != ?", kode, id).Count(&count)
	if count > 0 {
		return nil, errors.New("kode sparepart sudah digunakan")
	}

	// Perbarui field (stok sengaja TIDAK diperbarui di sini sesuai aturan bisnis 2.5)
	item.Kode = kode
	item.Nama = nama
	item.Kategori = strings.TrimSpace(req.Kategori)
	item.Merek = strings.TrimSpace(req.Merek)
	item.Harga = req.Harga
	item.StokMinimum = req.StokMinimum
	item.LokasiRak = strings.TrimSpace(req.LokasiRak)

	if err := s.db.Save(item).Error; err != nil {
		return nil, errors.New("gagal memperbarui sparepart")
	}

	return item, nil
}

// Delete menghapus sparepart menggunakan soft delete
func (s *SparepartService) Delete(id uint) error {
	item, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Soft delete menggunakan fitur bawaan GORM DeletedAt
	if err := s.db.Delete(item).Error; err != nil {
		return errors.New("gagal menghapus sparepart")
	}

	return nil
}

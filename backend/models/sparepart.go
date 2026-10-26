package models

import (
	"time"

	"gorm.io/gorm"
)

// Sparepart merepresentasikan data barang sparepart di gudang/bengkel
type Sparepart struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Kode        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"kode"`
	Nama        string         `gorm:"type:varchar(100);not null" json:"nama"`
	Kategori    string         `gorm:"type:varchar(50)" json:"kategori"`
	Merek       string         `gorm:"type:varchar(50)" json:"merek"`
	Harga       int64          `gorm:"default:0" json:"harga"`
	Stok        int            `gorm:"default:0" json:"stok"`
	StokMinimum int            `gorm:"default:0" json:"stok_minimum"`
	LokasiRak   string         `gorm:"type:varchar(50)" json:"lokasi_rak"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

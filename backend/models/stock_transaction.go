package models

import (
	"time"
)

// StockTransaction mencatat riwayat perubahan stok barang (masuk / keluar)
type StockTransaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SparepartID uint      `gorm:"not null;index" json:"sparepart_id"`
	Sparepart   Sparepart `gorm:"foreignKey:SparepartID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"sparepart,omitempty"`
	Tipe        string    `gorm:"type:varchar(10);not null" json:"tipe"` // 'masuk' atau 'keluar'
	Jumlah      int       `gorm:"not null" json:"jumlah"`
	Catatan     string    `gorm:"type:text" json:"catatan"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

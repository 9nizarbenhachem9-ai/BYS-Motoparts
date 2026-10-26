package models

import (
	"time"
)

// User merepresentasikan data akun administrator
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // Jangan tampilkan password di JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

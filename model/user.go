package model

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"user_id"`
	NamaLengkap  string    `gorm:"size:255;not null" json:"nama_lengkap"`
	Email        string    `gorm:"size:255;unique;not null" json:"email"`
	NomorTelepon string    `gorm:"size:20;unique" json:"nomor_telepon"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	StatusAkun   string    `gorm:"size:20;not null;default:'aktif'" json:"status_akun"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
} 
package model

import "time"

type Terminal struct {
	TerminalID    int       `gorm:"primaryKey;autoIncrement" json:"terminal_id"`
	NamaTerminal  string    `gorm:"size:100;unique;not null" json:"nama_terminal"`
	Lokasi        string    `gorm:"type:text" json:"lokasi"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type CreateTerminalRequest struct {
	NamaTerminal string `json:"nama_terminal" binding:"required"`
	Lokasi       string `json:"lokasi" binding:"required"`
}

type TerminalResponse struct {
	TerminalID    int       `json:"terminal_id"`
	NamaTerminal  string    `json:"nama_terminal"`
	Lokasi        string    `json:"lokasi"`
	CreatedAt     time.Time `json:"created_at"`
} 
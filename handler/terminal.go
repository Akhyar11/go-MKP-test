package handler

import (
	"encoding/json"
	"net/http"

	"tester/config"
	"tester/model"
)

func CreateTerminalHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req model.CreateTerminalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if req.NamaTerminal == "" || req.Lokasi == "" {
		http.Error(w, "Nama terminal dan lokasi harus diisi", http.StatusBadRequest)
		return
	}

	// Create terminal baru
	terminal := model.Terminal{
		NamaTerminal: req.NamaTerminal,
		Lokasi:      req.Lokasi,
	}

	// Simpan ke database
	if err := config.DB.Create(&terminal).Error; err != nil {
		http.Error(w, "Gagal menyimpan terminal", http.StatusInternalServerError)
		return
	}

	// Siapkan response
	response := model.TerminalResponse{
		TerminalID:   terminal.TerminalID,
		NamaTerminal: terminal.NamaTerminal,
		Lokasi:      terminal.Lokasi,
		CreatedAt:    terminal.CreatedAt,
	}

	// Kirim response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
} 
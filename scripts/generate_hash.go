package main

import (
	"fmt"
	"log"

	"tester/config"
	"tester/util"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init DB
	config.InitDB()

	// Generate hash untuk password "password123"
	password := "password123"
	hash, err := util.HashPassword(password)
	if err != nil {
		log.Fatal("Error generating hash:", err)
	}

	// Update user yang ada dengan hash baru
	result := config.DB.Exec("UPDATE users SET password_hash = ? WHERE email = ?", hash, "admin@example.com")
	if result.Error != nil {
		log.Fatal("Error updating user:", result.Error)
	}

	if result.RowsAffected == 0 {
		// Jika user belum ada, buat baru
		result = config.DB.Exec(`
			INSERT INTO users (nama_lengkap, email, nomor_telepon, password_hash, status_akun)
			VALUES (?, ?, ?, ?, ?)
		`, "Admin User", "admin@example.com", "081234567890", hash, "aktif")
		if result.Error != nil {
			log.Fatal("Error creating user:", result.Error)
		}
	}

	fmt.Println("Password hash updated successfully!")
	fmt.Println("You can now login with:")
	fmt.Println("Email: admin@example.com")
	fmt.Println("Password: password123")
} 
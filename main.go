package main

import (
	"log"
	"net/http"

	"tester/config"
	"tester/handler"
	"tester/middleware"

	"github.com/gorilla/mux"
)

func main() {
	config.InitDB()

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/api/login", handler.LoginHandler).Methods("POST")

	// Protected routes
	r.HandleFunc("/api/terminals", middleware.AuthMiddleware(handler.CreateTerminalHandler)).Methods("POST")

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}

// LoginHandler hanya kerangka, implementasi di handler/login.go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"message": "Not implemented"}`))
} 
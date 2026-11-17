package main

import (
	"log"
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/database"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/users"
	"github.com/hllqbl1/DriveRent/backend/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Memuat environment variable
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env: ", err)
	}

	// Menghubungkan database
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal("Koneksi database gagal: ", err)
	}
	defer db.Close()

	log.Println("Database berhasil terhubung!")

	// Membuat repository
	userRepo := users.NewUserRepository(db)

	// Membuat service
	userService := users.NewUserService(userRepo)

	// Membuat handler
	userHandler := users.NewUserHandler(userService)

	// Mendaftarkan routes
	mux := routes.SetupRoutes(userHandler)

	// Menjalankan HTTP server
	log.Println("Server berjalan di http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server gagal dijalankan: ", err)
	}
}

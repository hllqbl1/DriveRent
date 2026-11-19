package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/hllqbl1/DriveRent/backend/internal/database"
	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/auth"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/payments"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/rentals"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/reservations"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/returns"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/users"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/vehicles"
	"github.com/hllqbl1/DriveRent/backend/internal/routes"
)

func main() {
	// =========================
	// LOAD ENVIRONMENT
	// =========================

	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env: ", err)
	}

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET belum dikonfigurasi di file .env")
	}

	// JWT Expiration
	expirationHours := 24

	if value := os.Getenv("JWT_EXPIRE_HOURS"); value != "" {
		hours, err := strconv.Atoi(value)
		if err != nil || hours <= 0 {
			log.Fatal("JWT_EXPIRE_HOURS harus berupa angka positif")
		}
		expirationHours = hours
	}

	tokenConfig := auth.TokenConfig{
		Secret:     jwtSecret,
		Expiration: time.Duration(expirationHours) * time.Hour,
	}

	// =========================
	// DATABASE CONNECTION
	// =========================

	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal("Koneksi database gagal: ", err)
	}
	defer db.Close()

	log.Println("Database berhasil terhubung!")

	// =========================
	// USERS MODULE
	// =========================

	userRepo := users.NewUserRepository(db)
	userService := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService)

	// =========================
	// AUTH MODULE
	// =========================

	authService := auth.NewAuthService(
		userService,
		tokenConfig,
	)
	authHandler := auth.NewAuthHandler(authService)

	// =========================
	// VEHICLES MODULE
	// =========================

	vehicleRepo := vehicles.NewVehicleRepository(db)
	vehicleService := vehicles.NewVehicleService(vehicleRepo)
	vehicleHandler := vehicles.NewVehicleHandler(vehicleService)

	// =========================
	// RESERVATIONS MODULE
	// =========================

	reservationRepo := reservations.NewReservationRepository(db)

	reservationService := reservations.NewReservationService(
		reservationRepo,
		vehicleRepo,
	)

	reservationHandler := reservations.NewReservationHandler(
		reservationService,
	)

	// =========================
	// PAYMENTS MODULE
	// =========================

	paymentRepo := payments.NewPaymentRepository(db)
	paymentService := payments.NewPaymentService(paymentRepo)
	paymentHandler := payments.NewPaymentHandler(paymentService)

	// =========================
	// RENTALS MODULE
	// =========================

	rentalRepo := rentals.NewRentalRepository(db)
	rentalService := rentals.NewRentalService(rentalRepo)
	rentalHandler := rentals.NewRentalHandler(rentalService)

	// =========================
	// RETURNS MODULE
	// =========================

	returnRepo := returns.NewReturnRepository(db)
	returnService := returns.NewReturnService(returnRepo)
	returnHandler := returns.NewReturnHandler(returnService)

	// =========================
	// REGISTER ROUTES
	// =========================

	mux := routes.SetupRoutes(
		authHandler,
		userHandler,
		vehicleHandler,
		reservationHandler,
		paymentHandler,
		rentalHandler,
		returnHandler,
		jwtSecret,
	)

	// =========================
	// STATIC FILES (VEHICLE IMAGES)
	// =========================

	uploadDir, err := filepath.Abs("./uploads")
	if err != nil {
		log.Fatal("Gagal menentukan folder upload: ", err)
	}

	if err := os.MkdirAll(
		filepath.Join(uploadDir, "vehicles"),
		0755,
	); err != nil {
		log.Fatal("Gagal membuat folder upload: ", err)
	}

	mux.Handle(
		"GET /uploads/",
		http.StripPrefix(
			"/uploads/",
			http.FileServer(http.Dir(uploadDir)),
		),
	)

	// =========================
	// START HTTP SERVER
	// =========================

	log.Println("Server berjalan di port 8080")

	handler := middleware.CORSMiddleware(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

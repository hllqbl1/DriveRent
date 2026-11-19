package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/rentals"
)

func RegisterRentalRoutes(
	mux *http.ServeMux,
	handler rentals.RentalHandler,
	jwtSecret string,
) {
	auth := middleware.AuthMiddleware(jwtSecret)
	adminOrStaff := middleware.RequireRole("admin", "staff")

	// Admin dan Staff: proses serah terima kendaraan
	mux.Handle("POST /api/rentals",
		auth(adminOrStaff(http.HandlerFunc(handler.CreateRental))),
	)

	// Admin dan Staff: melihat seluruh data rental
	mux.Handle("GET /api/rentals",
		auth(adminOrStaff(http.HandlerFunc(handler.GetAllRentals))),
	)

	mux.Handle("GET /api/rentals/{id}",
		auth(adminOrStaff(http.HandlerFunc(handler.GetRentalByID))),
	)

	// User yang login: handler harus memeriksa kepemilikan reservasi
	mux.Handle("GET /api/reservations/{reservation_id}/rental",
		auth(http.HandlerFunc(handler.GetRentalByReservationID)),
	)
}

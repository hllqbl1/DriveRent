package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/reservations"
)

func RegisterReservationRoutes(
	mux *http.ServeMux,
	handler reservations.ReservationHandler,
	jwtSecret string,
) {
	auth := middleware.AuthMiddleware(jwtSecret)
	adminOrStaff := middleware.RequireRole("admin", "staff")

	// Customer yang sudah login bisa membuat reservasi
	mux.Handle("POST /api/reservations",
		auth(http.HandlerFunc(handler.CreateReservation)),
	)

	// Admin dan Staff
	mux.Handle("GET /api/reservations",
		auth(adminOrStaff(http.HandlerFunc(handler.GetAllReservations))),
	)

	mux.Handle("GET /api/reservations/{id}",
		auth(adminOrStaff(http.HandlerFunc(handler.GetReservationByID))),
	)

	mux.Handle("GET /api/users/{user_id}/reservations",
		auth(adminOrStaff(http.HandlerFunc(handler.GetReservationsByUserID))),
	)

	mux.Handle("PATCH /api/reservations/{id}/status",
		auth(adminOrStaff(http.HandlerFunc(handler.UpdateReservationStatus))),
	)

	// Customer yang login juga bisa membatalkan reservasinya
	mux.Handle("PATCH /api/reservations/{id}/cancel",
		auth(http.HandlerFunc(handler.CancelReservation)),
	)
}

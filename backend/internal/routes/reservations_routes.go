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
	adminOnly := middleware.RequireRole("admin")

	// Customer yang sudah login bisa membuat reservasi.
	mux.Handle(
		"POST /api/reservations",
		auth(http.HandlerFunc(handler.CreateReservation)),
	)

	// Admin bisa melihat seluruh reservasi.
	mux.Handle(
		"GET /api/reservations",
		auth(adminOnly(http.HandlerFunc(handler.GetAllReservations))),
	)

	// Admin bisa melihat detail reservasi.
	mux.Handle(
		"GET /api/reservations/{id}",
		auth(adminOnly(http.HandlerFunc(handler.GetReservationByID))),
	)

	// Customer bisa melihat reservasi miliknya.
	// Admin juga bisa melihat reservasi berdasarkan user ID.
	// Validasi kepemilikan customer wajib dilakukan di handler.
	mux.Handle(
		"GET /api/users/{user_id}/reservations",
		auth(http.HandlerFunc(handler.GetReservationsByUserID)),
	)

	// Hanya admin yang bisa mengubah status reservasi.
	mux.Handle(
		"PATCH /api/reservations/{id}/status",
		auth(adminOnly(http.HandlerFunc(handler.UpdateReservationStatus))),
	)

	// Customer yang login bisa membatalkan reservasinya sendiri.
	// Validasi kepemilikan wajib dilakukan di handler.
	mux.Handle(
		"PATCH /api/reservations/{id}/cancel",
		auth(http.HandlerFunc(handler.CancelReservation)),
	)
}

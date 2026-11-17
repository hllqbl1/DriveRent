package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/modules/reservations"
)

func RegisterReservationRoutes(
	mux *http.ServeMux,
	handler reservations.ReservationHandler,
) {
	// Reservation Routes
	mux.HandleFunc("POST /api/reservations", handler.CreateReservation)
	mux.HandleFunc("GET /api/reservations", handler.GetAllReservations)
	mux.HandleFunc("GET /api/reservations/{id}", handler.GetReservationByID)

	// User Reservation Routes
	mux.HandleFunc("GET /api/users/{user_id}/reservations", handler.GetReservationsByUserID)

	// Reservation Status Routes
	mux.HandleFunc("PATCH /api/reservations/{id}/status", handler.UpdateReservationStatus)
	mux.HandleFunc("PATCH /api/reservations/{id}/cancel", handler.CancelReservation)
}

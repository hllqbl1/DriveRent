package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/returns"
)

func RegisterReturnRoutes(
	mux *http.ServeMux,
	handler returns.ReturnHandler,
	jwtSecret string,
) {
	auth := middleware.AuthMiddleware(jwtSecret)
	adminOrStaff := middleware.RequireRole("admin", "staff")

	// Admin dan Staff: memproses pengembalian mobil
	mux.Handle("POST /api/returns",
		auth(adminOrStaff(http.HandlerFunc(handler.CreateReturn))),
	)

	// Admin dan Staff: melihat data pengembalian
	mux.Handle("GET /api/returns",
		auth(adminOrStaff(http.HandlerFunc(handler.GetAllReturns))),
	)

	mux.Handle("GET /api/returns/{id}",
		auth(adminOrStaff(http.HandlerFunc(handler.GetReturnByID))),
	)

	// User yang login: handler harus memeriksa kepemilikan rental
	mux.Handle("GET /api/rentals/{rental_id}/return",
		auth(http.HandlerFunc(handler.GetReturnByRentalID)),
	)
}

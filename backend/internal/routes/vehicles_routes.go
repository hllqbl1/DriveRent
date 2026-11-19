package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/vehicles"
)

func RegisterVehicleRoutes(
	mux *http.ServeMux,
	handler vehicles.VehicleHandler,
	jwtSecret string,
) {
	auth := middleware.AuthMiddleware(jwtSecret)
	adminOrStaff := middleware.RequireRole("admin", "staff")
	adminOnly := middleware.RequireRole("admin")

	// Public: semua orang bisa melihat kendaraan
	mux.HandleFunc("GET /api/vehicles", handler.GetAllVehicles)
	mux.HandleFunc("GET /api/vehicles/{id}", handler.GetVehicleByID)

	// Admin dan Staff
	mux.Handle("POST /api/vehicles",
		auth(adminOrStaff(http.HandlerFunc(handler.CreateVehicle))),
	)
	mux.Handle("PUT /api/vehicles/{id}",
		auth(adminOrStaff(http.HandlerFunc(handler.UpdateVehicle))),
	)
	mux.Handle("PATCH /api/vehicles/{id}/status",
		auth(adminOrStaff(http.HandlerFunc(handler.UpdateVehicleStatus))),
	)

	// Admin saja
	mux.Handle("DELETE /api/vehicles/{id}",
		auth(adminOnly(http.HandlerFunc(handler.DeleteVehicle))),
	)
}

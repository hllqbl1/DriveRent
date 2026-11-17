package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/modules/vehicles"
)

func RegisterVehicleRoutes(
	mux *http.ServeMux,
	handler vehicles.VehicleHandler,
) {
	// Vehicle Routes
	mux.HandleFunc("POST /api/vehicles", handler.CreateVehicle)
	mux.HandleFunc("GET /api/vehicles", handler.GetAllVehicles)
	mux.HandleFunc("GET /api/vehicles/{id}", handler.GetVehicleByID)
	mux.HandleFunc("PUT /api/vehicles/{id}", handler.UpdateVehicle)
	mux.HandleFunc("PATCH /api/vehicles/{id}/status", handler.UpdateVehicleStatus)
	mux.HandleFunc("DELETE /api/vehicles/{id}", handler.DeleteVehicle)
}

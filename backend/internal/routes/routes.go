package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/modules/auth"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/payments"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/rentals"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/reservations"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/returns"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/users"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/vehicles"
)

func SetupRoutes(
	authHandler auth.AuthHandler,
	userHandler users.UserHandler,
	vehicleHandler vehicles.VehicleHandler,
	reservationHandler reservations.ReservationHandler,
	paymentHandler payments.PaymentHandler,
	rentalHandler rentals.RentalHandler,
	returnHandler returns.ReturnHandler,
	jwtSecret string,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Auth Routes
	RegisterAuthRoutes(mux, authHandler, jwtSecret)

	// User Routes
	RegisterUserRoutes(mux, userHandler, jwtSecret)

	// Vehicle Routes
	RegisterVehicleRoutes(mux, vehicleHandler, jwtSecret)

	// Reservation Routes
	RegisterReservationRoutes(mux, reservationHandler, jwtSecret)

	// Payment Routes
	RegisterPaymentRoutes(mux, paymentHandler, jwtSecret)

	// Rental Routes
	RegisterRentalRoutes(mux, rentalHandler, jwtSecret)

	// Return Routes
	RegisterReturnRoutes(mux, returnHandler, jwtSecret)

	return mux
}

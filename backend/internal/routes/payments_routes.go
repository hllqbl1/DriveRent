package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/payments"
)

func RegisterPaymentRoutes(
	mux *http.ServeMux,
	handler payments.PaymentHandler,
	jwtSecret string,
) {
	auth := middleware.AuthMiddleware(jwtSecret)
	adminOrStaff := middleware.RequireRole("admin", "staff")

	// Customer yang sudah login bisa membuat pembayaran
	mux.Handle("POST /api/payments",
		auth(http.HandlerFunc(handler.CreatePayment)),
	)

	// Admin dan Staff
	mux.Handle("GET /api/payments",
		auth(adminOrStaff(http.HandlerFunc(handler.GetAllPayments))),
	)

	mux.Handle("GET /api/payments/{id}",
		auth(adminOrStaff(http.HandlerFunc(handler.GetPaymentByID))),
	)

	// Customer yang login bisa melihat riwayat pembayaran.
	// Handler tetap harus memeriksa kepemilikan reservasi.
	mux.Handle("GET /api/reservations/{reservation_id}/payments",
		auth(http.HandlerFunc(handler.GetPaymentsByReservationID)),
	)

	// Admin dan Staff
	mux.Handle("PATCH /api/payments/{id}/verify",
		auth(adminOrStaff(http.HandlerFunc(handler.VerifyPayment))),
	)

	mux.Handle("PATCH /api/payments/{id}/reject",
		auth(adminOrStaff(http.HandlerFunc(handler.RejectPayment))),
	)
}

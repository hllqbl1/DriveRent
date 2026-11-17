package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/modules/payments"
)

func RegisterPaymentRoutes(
	mux *http.ServeMux,
	handler payments.PaymentHandler,
) {
	// Payment
	mux.HandleFunc("POST /api/payments", handler.CreatePayment)
	mux.HandleFunc("GET /api/payments", handler.GetAllPayments)
	mux.HandleFunc("GET /api/payments/{id}", handler.GetPaymentByID)

	// Payment history by reservation
	mux.HandleFunc(
		"GET /api/reservations/{reservation_id}/payments",
		handler.GetPaymentsByReservationID,
	)

	// Payment verification
	mux.HandleFunc(
		"PATCH /api/payments/{id}/verify",
		handler.VerifyPayment,
	)

	mux.HandleFunc(
		"PATCH /api/payments/{id}/reject",
		handler.RejectPayment,
	)
}

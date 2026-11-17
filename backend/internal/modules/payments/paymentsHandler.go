package payments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type PaymentHandler interface {
	CreatePayment(w http.ResponseWriter, r *http.Request)
	GetAllPayments(w http.ResponseWriter, r *http.Request)
	GetPaymentByID(w http.ResponseWriter, r *http.Request)
	GetPaymentsByReservationID(w http.ResponseWriter, r *http.Request)
	VerifyPayment(w http.ResponseWriter, r *http.Request)
	RejectPayment(w http.ResponseWriter, r *http.Request)
}

type paymentHandler struct {
	service PaymentService
}

func NewPaymentHandler(service PaymentService) PaymentHandler {
	return &paymentHandler{
		service: service,
	}
}

// Helper untuk mengirim JSON response.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

// Helper untuk mengirim error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

// CreatePayment membuat pembayaran baru.
func (h *paymentHandler) CreatePayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		ReservationID int64   `json:"reservation_id"`
		PaymentMethod string  `json:"payment_method"`
		PaymentProof  *string `json:"payment_proof"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Format JSON tidak valid",
		)
		return
	}

	err := h.service.CreatePayment(
		r.Context(),
		request.ReservationID,
		request.PaymentMethod,
		request.PaymentProof,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrReservationNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrReservationNotPending),
			errors.Is(err, ErrPaymentAlreadyExists):
			writeError(w, http.StatusConflict, err.Error())

		case errors.Is(err, ErrInvalidPayment),
			errors.Is(err, ErrInvalidMethod):
			writeError(w, http.StatusBadRequest, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal membuat pembayaran",
			)
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Pembayaran berhasil dibuat dan menunggu verifikasi",
	})
}

// GetAllPayments mengambil seluruh pembayaran.
func (h *paymentHandler) GetAllPayments(
	w http.ResponseWriter,
	r *http.Request,
) {
	paymentList, err := h.service.GetAllPayments(r.Context())

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil data pembayaran",
		)
		return
	}

	writeJSON(w, http.StatusOK, paymentList)
}

// GetPaymentByID mengambil detail pembayaran.
func (h *paymentHandler) GetPaymentByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"ID pembayaran tidak valid",
		)
		return
	}

	payment, err := h.service.GetPaymentByID(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil pembayaran",
		)
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

// GetPaymentsByReservationID mengambil riwayat pembayaran reservasi.
func (h *paymentHandler) GetPaymentsByReservationID(
	w http.ResponseWriter,
	r *http.Request,
) {
	reservationID, err := strconv.ParseInt(
		r.PathValue("reservation_id"),
		10,
		64,
	)

	if err != nil || reservationID <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"ID reservasi tidak valid",
		)
		return
	}

	paymentList, err := h.service.GetPaymentsByReservationID(
		r.Context(),
		reservationID,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil riwayat pembayaran",
		)
		return
	}

	writeJSON(w, http.StatusOK, paymentList)
}

// VerifyPayment memverifikasi pembayaran.
func (h *paymentHandler) VerifyPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"ID pembayaran tidak valid",
		)
		return
	}

	err = h.service.VerifyPayment(
		r.Context(),
		id,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrPaymentNotFound),
			errors.Is(err, ErrReservationNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrPaymentAlreadyProcessed),
			errors.Is(err, ErrReservationNotPending):
			writeError(w, http.StatusConflict, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal memverifikasi pembayaran",
			)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Pembayaran berhasil diverifikasi dan reservasi dikonfirmasi",
	})
}

// RejectPayment menolak pembayaran.
func (h *paymentHandler) RejectPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"ID pembayaran tidak valid",
		)
		return
	}

	err = h.service.RejectPayment(
		r.Context(),
		id,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrPaymentNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrPaymentAlreadyProcessed):
			writeError(w, http.StatusConflict, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal menolak pembayaran",
			)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Pembayaran berhasil ditolak",
	})
}

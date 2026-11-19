package reservations

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
)

type ReservationHandler interface {
	CreateReservation(w http.ResponseWriter, r *http.Request)
	GetAllReservations(w http.ResponseWriter, r *http.Request)
	GetReservationByID(w http.ResponseWriter, r *http.Request)
	GetReservationsByUserID(w http.ResponseWriter, r *http.Request)
	UpdateReservationStatus(w http.ResponseWriter, r *http.Request)
	CancelReservation(w http.ResponseWriter, r *http.Request)
}

type reservationHandler struct {
	service ReservationService
}

func NewReservationHandler(service ReservationService) ReservationHandler {
	return &reservationHandler{
		service: service,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

// CreateReservation membuat reservasi baru.
// CreateReservation membuat reservasi baru.
func (h *reservationHandler) CreateReservation(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		VehicleID int64  `json:"vehicle_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	// Ambil User ID dari JWT, bukan dari request body.
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		writeError(w, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}

	if request.VehicleID <= 0 {
		writeError(w, http.StatusBadRequest, "Vehicle ID tidak valid")
		return
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Format start_date harus YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Format end_date harus YYYY-MM-DD")
		return
	}

	err = h.service.CreateReservation(
		r.Context(),
		userID,
		request.VehicleID,
		startDate,
		endDate,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrVehicleUnavailable):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, ErrDateConflict):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, ErrInvalidReservation):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Reservasi berhasil dibuat",
	})
}

// GetAllReservations mengambil seluruh reservasi.
func (h *reservationHandler) GetAllReservations(
	w http.ResponseWriter,
	r *http.Request,
) {
	reservations, err := h.service.GetAllReservations(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil data reservasi")
		return
	}

	writeJSON(w, http.StatusOK, reservations)
}

// GetReservationByID mengambil detail reservasi.
func (h *reservationHandler) GetReservationByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ID reservasi tidak valid")
		return
	}

	reservation, err := h.service.GetReservationByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrReservationNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, "Gagal mengambil detail reservasi")
		return
	}

	writeJSON(w, http.StatusOK, reservation)
}

// GetReservationsByUserID mengambil reservasi berdasarkan user ID.
func (h *reservationHandler) GetReservationsByUserID(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		writeError(w, http.StatusBadRequest, "User ID tidak valid")
		return
	}

	reservations, err := h.service.GetReservationsByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil reservasi user")
		return
	}

	writeJSON(w, http.StatusOK, reservations)
}

// UpdateReservationStatus mengubah status reservasi.
func (h *reservationHandler) UpdateReservationStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ID reservasi tidak valid")
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	err = h.service.UpdateReservationStatus(r.Context(), id, request.Status)
	if err != nil {
		switch {
		case errors.Is(err, ErrReservationNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Status reservasi berhasil diperbarui",
	})
}

// CancelReservation membatalkan reservasi.
func (h *reservationHandler) CancelReservation(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ID reservasi tidak valid")
		return
	}

	err = h.service.CancelReservation(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrReservationNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Reservasi berhasil dibatalkan",
	})
}

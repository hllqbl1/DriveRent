package rentals

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type RentalHandler interface {
	CreateRental(w http.ResponseWriter, r *http.Request)
	GetAllRentals(w http.ResponseWriter, r *http.Request)
	GetRentalByID(w http.ResponseWriter, r *http.Request)
	GetRentalByReservationID(w http.ResponseWriter, r *http.Request)
}

type rentalHandler struct {
	service RentalService
}

func NewRentalHandler(service RentalService) RentalHandler {
	return &rentalHandler{
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

// CreateRental membuat rental baru.
func (h *rentalHandler) CreateRental(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		ReservationID   int64  `json:"reservation_id"`
		StaffID         int64  `json:"staff_id"`
		PickupDate      string `json:"pickup_date"`
		InitialFuel     string `json:"initial_fuel"`
		InitialOdometer int    `json:"initial_odometer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Format JSON tidak valid",
		)
		return
	}

	var pickupDate time.Time

	if request.PickupDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, request.PickupDate)
		if err != nil {
			writeError(
				w,
				http.StatusBadRequest,
				"Format pickup_date harus RFC3339",
			)
			return
		}

		pickupDate = parsedDate
	}

	rental := Rental{
		ReservationID:   request.ReservationID,
		StaffID:         request.StaffID,
		PickupDate:      pickupDate,
		InitialFuel:     request.InitialFuel,
		InitialOdometer: request.InitialOdometer,
	}

	err := h.service.CreateRental(r.Context(), rental)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRental),
			errors.Is(err, ErrInvalidInitialFuel):
			writeError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, ErrReservationNotFound),
			errors.Is(err, ErrRentalNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrReservationNotReady),
			errors.Is(err, ErrRentalAlreadyExists),
			errors.Is(err, ErrVehicleNotAvailable):
			writeError(w, http.StatusConflict, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal membuat rental",
			)
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Rental berhasil dibuat",
	})
}

// GetAllRentals mengambil seluruh data rental.
func (h *rentalHandler) GetAllRentals(
	w http.ResponseWriter,
	r *http.Request,
) {
	rentalList, err := h.service.GetAllRentals(r.Context())

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil data rental",
		)
		return
	}

	writeJSON(w, http.StatusOK, rentalList)
}

// GetRentalByID mengambil rental berdasarkan ID.
func (h *rentalHandler) GetRentalByID(
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
			"ID rental tidak valid",
		)
		return
	}

	rental, err := h.service.GetRentalByID(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrRentalNotFound) {
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
			"Gagal mengambil data rental",
		)
		return
	}

	writeJSON(w, http.StatusOK, rental)
}

// GetRentalByReservationID mengambil rental berdasarkan reservasi.
func (h *rentalHandler) GetRentalByReservationID(
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

	rental, err := h.service.GetRentalByReservationID(
		r.Context(),
		reservationID,
	)

	if err != nil {
		if errors.Is(err, ErrRentalNotFound) {
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
			"Gagal mengambil data rental",
		)
		return
	}

	writeJSON(w, http.StatusOK, rental)
}

package returns

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type ReturnHandler interface {
	CreateReturn(w http.ResponseWriter, r *http.Request)
	GetAllReturns(w http.ResponseWriter, r *http.Request)
	GetReturnByID(w http.ResponseWriter, r *http.Request)
	GetReturnByRentalID(w http.ResponseWriter, r *http.Request)
}

type returnHandler struct {
	service ReturnService
}

func NewReturnHandler(service ReturnService) ReturnHandler {
	return &returnHandler{
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

// CreateReturn membuat data pengembalian.
func (h *returnHandler) CreateReturn(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request struct {
		RentalID      int64   `json:"rental_id"`
		ReturnDate    string  `json:"return_date"`
		FinalFuel     string  `json:"final_fuel"`
		FinalOdometer int     `json:"final_odometer"`
		LateFee       float64 `json:"late_fee"`
		DamageFee     float64 `json:"damage_fee"`
		Notes         *string `json:"notes"`
		VehicleStatus string  `json:"vehicle_status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"Format JSON tidak valid",
		)
		return
	}

	var returnDate time.Time

	if request.ReturnDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, request.ReturnDate)

		if err != nil {
			writeError(
				w,
				http.StatusBadRequest,
				"Format return_date harus RFC3339",
			)
			return
		}

		returnDate = parsedDate
	}

	returnData := Return{
		RentalID:      request.RentalID,
		ReturnDate:    returnDate,
		FinalFuel:     request.FinalFuel,
		FinalOdometer: request.FinalOdometer,
		LateFee:       request.LateFee,
		DamageFee:     request.DamageFee,
		Notes:         request.Notes,
	}

	err := h.service.CreateReturn(
		r.Context(),
		returnData,
		request.VehicleStatus,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidReturn),
			errors.Is(err, ErrInvalidFinalFuel),
			errors.Is(err, ErrInvalidVehicleStatus):
			writeError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, ErrRentalNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, ErrRentalNotActive),
			errors.Is(err, ErrReturnAlreadyExists),
			errors.Is(err, ErrOdometerInvalid):
			writeError(w, http.StatusConflict, err.Error())

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"Gagal membuat data pengembalian",
			)
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Pengembalian kendaraan berhasil dicatat",
	})
}

// GetAllReturns mengambil seluruh data pengembalian.
func (h *returnHandler) GetAllReturns(
	w http.ResponseWriter,
	r *http.Request,
) {
	returnList, err := h.service.GetAllReturns(r.Context())

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Gagal mengambil data pengembalian",
		)
		return
	}

	writeJSON(w, http.StatusOK, returnList)
}

// GetReturnByID mengambil pengembalian berdasarkan ID.
func (h *returnHandler) GetReturnByID(
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
			"ID pengembalian tidak valid",
		)
		return
	}

	returnData, err := h.service.GetReturnByID(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, ErrReturnNotFound) {
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
			"Gagal mengambil data pengembalian",
		)
		return
	}

	writeJSON(w, http.StatusOK, returnData)
}

// GetReturnByRentalID mengambil pengembalian berdasarkan rental.
func (h *returnHandler) GetReturnByRentalID(
	w http.ResponseWriter,
	r *http.Request,
) {
	rentalID, err := strconv.ParseInt(
		r.PathValue("rental_id"),
		10,
		64,
	)

	if err != nil || rentalID <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"ID rental tidak valid",
		)
		return
	}

	returnData, err := h.service.GetReturnByRentalID(
		r.Context(),
		rentalID,
	)

	if err != nil {
		if errors.Is(err, ErrReturnNotFound) {
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
			"Gagal mengambil data pengembalian",
		)
		return
	}

	writeJSON(w, http.StatusOK, returnData)
}

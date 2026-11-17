package vehicles

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type VehicleHandler interface {
	CreateVehicle(w http.ResponseWriter, r *http.Request)
	GetAllVehicles(w http.ResponseWriter, r *http.Request)
	GetVehicleByID(w http.ResponseWriter, r *http.Request)
	UpdateVehicle(w http.ResponseWriter, r *http.Request)
	UpdateVehicleStatus(w http.ResponseWriter, r *http.Request)
	DeleteVehicle(w http.ResponseWriter, r *http.Request)
}

type vehicleHandler struct {
	service VehicleService
}

func NewVehicleHandler(service VehicleService) VehicleHandler {
	return &vehicleHandler{
		service: service,
	}
}

// CreateVehicle menambahkan kendaraan baru.
func (h *vehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var vehicle Vehicle

	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	err = h.service.CreateVehicle(r.Context(), vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kendaraan berhasil ditambahkan",
	})
}

// GetAllVehicles mengambil seluruh kendaraan.
func (h *vehicleHandler) GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.service.GetAllVehicles(r.Context())
	if err != nil {
		http.Error(w, "Gagal mengambil data kendaraan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(vehicles)
}

// GetVehicleByID mengambil kendaraan berdasarkan ID.
func (h *vehicleHandler) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "ID kendaraan tidak valid", http.StatusBadRequest)
		return
	}

	vehicle, err := h.service.GetVehicleByID(r.Context(), id)
	if err != nil {
		if err == ErrVehicleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Gagal mengambil data kendaraan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(vehicle)
}

// UpdateVehicle memperbarui informasi kendaraan.
func (h *vehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "ID kendaraan tidak valid", http.StatusBadRequest)
		return
	}

	var vehicle Vehicle

	err = json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	vehicle.ID = id

	err = h.service.UpdateVehicle(r.Context(), vehicle)
	if err != nil {
		switch err {
		case ErrVehicleNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case ErrLicensePlateUsed:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kendaraan berhasil diperbarui",
	})
}

// UpdateVehicleStatus mengubah status kendaraan.
func (h *vehicleHandler) UpdateVehicleStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "ID kendaraan tidak valid", http.StatusBadRequest)
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateVehicleStatus(r.Context(), id, request.Status)
	if err != nil {
		if err == ErrVehicleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Status kendaraan berhasil diperbarui",
	})
}

// DeleteVehicle menghapus kendaraan berdasarkan ID.
func (h *vehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "ID kendaraan tidak valid", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteVehicle(r.Context(), id)
	if err != nil {
		if err == ErrVehicleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Gagal menghapus kendaraan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kendaraan berhasil dihapus",
	})
}

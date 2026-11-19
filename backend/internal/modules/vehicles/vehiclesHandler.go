package vehicles

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxImageSize = 5 << 20 // 5 MB

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

// parseVehicleRequest membaca data kendaraan dari JSON atau multipart/form-data.
func (h *vehicleHandler) parseVehicleRequest(
	w http.ResponseWriter,
	r *http.Request,
) (Vehicle, string, error) {
	var vehicle Vehicle

	contentType := r.Header.Get("Content-Type")

	// Tetap mendukung JSON untuk kompatibilitas.
	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&vehicle)
		return vehicle, "", err
	}

	// Upload foto menggunakan multipart/form-data.
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return vehicle, "", fmt.Errorf("format request tidak didukung")
	}

	// Beri sedikit ruang tambahan untuk data form dan multipart boundary.
	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize+(1<<20))

	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		return vehicle, "", fmt.Errorf("gagal membaca form: ukuran request terlalu besar atau format tidak valid")
	}

	// Ambil data kendaraan.
	vehicle.Name = r.FormValue("name")
	vehicle.Brand = r.FormValue("brand")
	vehicle.Model = r.FormValue("model")
	vehicle.LicensePlate = r.FormValue("license_plate")
	vehicle.Category = r.FormValue("category")
	vehicle.Transmission = r.FormValue("transmission")
	vehicle.FuelType = r.FormValue("fuel_type")
	vehicle.Description = optionalString(r.FormValue("description"))

	var err error

	vehicle.Year, err = strconv.Atoi(r.FormValue("year"))
	if err != nil {
		return vehicle, "", fmt.Errorf("tahun kendaraan tidak valid")
	}

	vehicle.Capacity, err = strconv.Atoi(r.FormValue("capacity"))
	if err != nil {
		return vehicle, "", fmt.Errorf("kapasitas kendaraan tidak valid")
	}

	vehicle.PricePerDay, err = strconv.ParseFloat(r.FormValue("price_per_day"), 64)
	if err != nil {
		return vehicle, "", fmt.Errorf("harga sewa tidak valid")
	}

	// Status akan ditangani service saat kendaraan dibuat.
	vehicle.Status = r.FormValue("status")

	// Simpan foto jika admin memilih file.
	imagePath, err := saveVehicleImage(r)
	if err != nil {
		return vehicle, "", err
	}

	if imagePath != "" {
		vehicle.Image = &imagePath
	}

	return vehicle, imagePath, nil
}

// optionalString mengubah string kosong menjadi nil.
func optionalString(value string) *string {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil
	}

	return &value
}

// saveVehicleImage menyimpan foto kendaraan ke folder uploads/vehicles.
func saveVehicleImage(r *http.Request) (string, error) {
	file, _, err := r.FormFile("image")

	if err == http.ErrMissingFile {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("gagal membaca file foto")
	}
	defer file.Close()

	// Baca maksimal 5 MB + 1 byte untuk mendeteksi file terlalu besar.
	data, err := io.ReadAll(io.LimitReader(file, maxImageSize+1))
	if err != nil {
		return "", fmt.Errorf("gagal membaca foto")
	}

	if len(data) == 0 {
		return "", fmt.Errorf("file foto kosong")
	}

	if len(data) > maxImageSize {
		return "", fmt.Errorf("ukuran foto maksimal 5 MB")
	}

	// Periksa tipe file berdasarkan isi file, bukan nama file.
	contentType := http.DetectContentType(data)

	extension := ""

	switch contentType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "image/webp":
		extension = ".webp"
	default:
		return "", fmt.Errorf("format foto harus JPG, PNG, atau WebP")
	}

	// Buat nama file acak agar tidak bentrok.
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("gagal membuat nama file")
	}

	filename := hex.EncodeToString(randomBytes) + extension

	uploadDir := filepath.Join("uploads", "vehicles")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat folder upload")
	}

	filePath := filepath.Join(uploadDir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("gagal menyimpan foto")
	}

	// Path ini yang disimpan ke database.
	return "/uploads/vehicles/" + filename, nil
}

// CreateVehicle menambahkan kendaraan baru.
func (h *vehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	vehicle, savedImage, err := h.parseVehicleRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateVehicle(r.Context(), vehicle)
	if err != nil {
		// Hapus foto jika penyimpanan kendaraan gagal.
		if savedImage != "" {
			os.Remove(filepath.Join("uploads", "vehicles", filepath.Base(savedImage)))
		}

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

	// Ambil data kendaraan lama agar foto lama bisa dipertahankan.
	oldVehicle, err := h.service.GetVehicleByID(r.Context(), id)
	if err != nil {
		if err == ErrVehicleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Gagal mengambil data kendaraan", http.StatusInternalServerError)
		return
	}

	vehicle, savedImage, err := h.parseVehicleRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vehicle.ID = id

	// Jika tidak ada foto baru, gunakan foto lama.
	if vehicle.Image == nil {
		vehicle.Image = oldVehicle.Image
	}

	err = h.service.UpdateVehicle(r.Context(), vehicle)
	if err != nil {
		// Hapus foto baru jika update gagal.
		if savedImage != "" {
			os.Remove(filepath.Join("uploads", "vehicles", filepath.Base(savedImage)))
		}

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
		// Tampilkan penyebab asli kegagalan di terminal backend.
		log.Printf("Gagal menghapus kendaraan ID %d: %v", id, err)

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

package vehicles

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrVehicleNotFound  = errors.New("kendaraan tidak ditemukan")
	ErrLicensePlateUsed = errors.New("plat nomor sudah digunakan")
)

type VehicleService interface {
	CreateVehicle(ctx context.Context, vehicle Vehicle) error
	GetAllVehicles(ctx context.Context) ([]Vehicle, error)
	GetVehicleByID(ctx context.Context, id int64) (*Vehicle, error)
	UpdateVehicle(ctx context.Context, vehicle Vehicle) error
	UpdateVehicleStatus(ctx context.Context, id int64, status string) error
	DeleteVehicle(ctx context.Context, id int64) error
}

type vehicleService struct {
	repo VehicleRepository
}

func NewVehicleService(repo VehicleRepository) VehicleService {
	return &vehicleService{
		repo: repo,
	}
}

// CreateVehicle membuat kendaraan baru.
func (s *vehicleService) CreateVehicle(ctx context.Context, vehicle Vehicle) error {
	// Normalisasi data
	vehicle.Name = strings.TrimSpace(vehicle.Name)
	vehicle.Brand = strings.TrimSpace(vehicle.Brand)
	vehicle.Model = strings.TrimSpace(vehicle.Model)
	vehicle.LicensePlate = strings.ToUpper(strings.TrimSpace(vehicle.LicensePlate))
	vehicle.Category = strings.TrimSpace(vehicle.Category)
	vehicle.Transmission = strings.ToLower(strings.TrimSpace(vehicle.Transmission))
	vehicle.FuelType = strings.ToLower(strings.TrimSpace(vehicle.FuelType))

	// Validasi data
	if err := validateVehicle(vehicle); err != nil {
		return err
	}

	// Cek duplikasi plat nomor
	used, err := s.repo.IsLicensePlateUsed(ctx, vehicle.LicensePlate, 0)
	if err != nil {
		return err
	}

	if used {
		return ErrLicensePlateUsed
	}

	// Kendaraan baru selalu dimulai dengan status available.
	vehicle.Status = "available"

	return s.repo.CreateVehicle(ctx, vehicle)
}

// GetAllVehicles mengambil seluruh kendaraan.
func (s *vehicleService) GetAllVehicles(ctx context.Context) ([]Vehicle, error) {
	return s.repo.GetAllVehicles(ctx)
}

// GetVehicleByID mengambil kendaraan berdasarkan ID.
func (s *vehicleService) GetVehicleByID(ctx context.Context, id int64) (*Vehicle, error) {
	if id <= 0 {
		return nil, errors.New("ID kendaraan tidak valid")
	}

	vehicle, err := s.repo.GetVehicleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVehicleNotFound
		}

		return nil, err
	}

	return vehicle, nil
}

// UpdateVehicle memperbarui informasi kendaraan selain status.
func (s *vehicleService) UpdateVehicle(ctx context.Context, vehicle Vehicle) error {
	if vehicle.ID <= 0 {
		return errors.New("ID kendaraan tidak valid")
	}

	// Pastikan kendaraan tersedia.
	_, err := s.GetVehicleByID(ctx, vehicle.ID)
	if err != nil {
		return err
	}

	// Normalisasi data
	vehicle.Name = strings.TrimSpace(vehicle.Name)
	vehicle.Brand = strings.TrimSpace(vehicle.Brand)
	vehicle.Model = strings.TrimSpace(vehicle.Model)
	vehicle.LicensePlate = strings.ToUpper(strings.TrimSpace(vehicle.LicensePlate))
	vehicle.Category = strings.TrimSpace(vehicle.Category)
	vehicle.Transmission = strings.ToLower(strings.TrimSpace(vehicle.Transmission))
	vehicle.FuelType = strings.ToLower(strings.TrimSpace(vehicle.FuelType))

	// Validasi data
	if err := validateVehicle(vehicle); err != nil {
		return err
	}

	// Cek duplikasi plat nomor.
	// ID kendaraan sendiri dikecualikan.
	used, err := s.repo.IsLicensePlateUsed(
		ctx,
		vehicle.LicensePlate,
		vehicle.ID,
	)
	if err != nil {
		return err
	}

	if used {
		return ErrLicensePlateUsed
	}

	return s.repo.UpdateVehicle(ctx, vehicle)
}

// UpdateVehicleStatus mengubah status kendaraan.
func (s *vehicleService) UpdateVehicleStatus(ctx context.Context, id int64, status string) error {
	if id <= 0 {
		return errors.New("ID kendaraan tidak valid")
	}

	status = strings.ToLower(strings.TrimSpace(status))

	// Validasi status.
	if !isValidVehicleStatus(status) {
		return errors.New("status kendaraan tidak valid")
	}

	// Pastikan kendaraan tersedia.
	_, err := s.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.UpdateVehicleStatus(ctx, id, status)
}

// DeleteVehicle menghapus kendaraan berdasarkan ID.
func (s *vehicleService) DeleteVehicle(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("ID kendaraan tidak valid")
	}

	// Pastikan kendaraan ditemukan.
	_, err := s.GetVehicleByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteVehicle(ctx, id)
}

// validateVehicle memvalidasi data kendaraan.
func validateVehicle(vehicle Vehicle) error {
	if vehicle.Name == "" {
		return errors.New("nama kendaraan wajib diisi")
	}

	if utf8.RuneCountInString(vehicle.Name) > 100 {
		return errors.New("nama kendaraan maksimal 100 karakter")
	}

	if vehicle.Brand == "" {
		return errors.New("brand kendaraan wajib diisi")
	}

	if utf8.RuneCountInString(vehicle.Brand) > 50 {
		return errors.New("brand maksimal 50 karakter")
	}

	if vehicle.Model == "" {
		return errors.New("model kendaraan wajib diisi")
	}

	if utf8.RuneCountInString(vehicle.Model) > 50 {
		return errors.New("model maksimal 50 karakter")
	}

	currentYear := time.Now().Year()

	if vehicle.Year < 1886 || vehicle.Year > currentYear+1 {
		return errors.New("tahun kendaraan tidak valid")
	}

	if vehicle.LicensePlate == "" {
		return errors.New("plat nomor wajib diisi")
	}

	if utf8.RuneCountInString(vehicle.LicensePlate) > 20 {
		return errors.New("plat nomor maksimal 20 karakter")
	}

	if vehicle.Category == "" {
		return errors.New("kategori kendaraan wajib diisi")
	}

	if utf8.RuneCountInString(vehicle.Category) > 50 {
		return errors.New("kategori maksimal 50 karakter")
	}

	if vehicle.Transmission != "manual" && vehicle.Transmission != "automatic" {
		return errors.New("transmission harus manual atau automatic")
	}

	if vehicle.FuelType == "" {
		return errors.New("jenis bahan bakar wajib diisi")
	}

	if utf8.RuneCountInString(vehicle.FuelType) > 20 {
		return errors.New("jenis bahan bakar maksimal 20 karakter")
	}

	if vehicle.Capacity <= 0 {
		return errors.New("kapasitas kendaraan harus lebih dari 0")
	}

	if vehicle.PricePerDay <= 0 ||
		math.IsNaN(vehicle.PricePerDay) ||
		math.IsInf(vehicle.PricePerDay, 0) {
		return errors.New("harga sewa harus lebih dari 0 dan valid")
	}

	return nil
}

// isValidVehicleStatus memeriksa status kendaraan.
func isValidVehicleStatus(status string) bool {
	switch status {
	case "available", "rented", "maintenance":
		return true
	default:
		return false
	}
}

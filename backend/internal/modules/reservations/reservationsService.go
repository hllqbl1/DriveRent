package reservations

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/hllqbl1/DriveRent/backend/internal/modules/vehicles"
)

var (
	ErrInvalidReservation = errors.New("data reservasi tidak valid")
	ErrInvalidStatus      = errors.New("status reservasi tidak valid")
)

type ReservationService interface {
	CreateReservation(ctx context.Context, userID, vehicleID int64, startDate, endDate time.Time) error
	GetAllReservations(ctx context.Context) ([]Reservation, error)
	GetReservationByID(ctx context.Context, id int64) (*Reservation, error)
	GetReservationsByUserID(ctx context.Context, userID int64) ([]Reservation, error)
	UpdateReservationStatus(ctx context.Context, id int64, status string) error
	CancelReservation(ctx context.Context, id int64) error
}

type reservationService struct {
	repo        ReservationRepository
	vehicleRepo vehicles.VehicleRepository
}

func NewReservationService(
	repo ReservationRepository,
	vehicleRepo vehicles.VehicleRepository,
) ReservationService {
	return &reservationService{
		repo:        repo,
		vehicleRepo: vehicleRepo,
	}
}

// CreateReservation membuat reservasi baru.
func (s *reservationService) CreateReservation(
	ctx context.Context,
	userID, vehicleID int64,
	startDate, endDate time.Time,
) error {
	if userID <= 0 || vehicleID <= 0 {
		return ErrInvalidReservation
	}

	// Normalisasi tanggal agar hanya menggunakan tanggal kalender.
	startDate = dateOnly(startDate)
	endDate = dateOnly(endDate)

	if !endDate.After(startDate) {
		return errors.New("tanggal selesai harus setelah tanggal mulai")
	}

	// Reservasi tidak boleh dimulai pada tanggal yang sudah lewat.
	today := dateOnly(time.Now())

	if startDate.Before(today) {
		return errors.New("tanggal mulai tidak boleh di masa lalu")
	}

	// Ambil data kendaraan untuk mendapatkan harga sewa.
	vehicle, err := s.vehicleRepo.GetVehicleByID(ctx, vehicleID)
	if err != nil {
		return err
	}

	if vehicle.Status != "available" {
		return ErrVehicleUnavailable
	}

	if vehicle.PricePerDay <= 0 ||
		math.IsNaN(vehicle.PricePerDay) ||
		math.IsInf(vehicle.PricePerDay, 0) {
		return errors.New("harga kendaraan tidak valid")
	}

	// Hitung total hari dari backend.
	totalDays := int(endDate.Sub(startDate).Hours() / 24)

	if totalDays <= 0 {
		return ErrInvalidReservation
	}

	// Hitung harga dari backend, bukan dari frontend.
	totalPrice := math.Round(vehicle.PricePerDay*float64(totalDays)*100) / 100

	if totalPrice <= 0 || math.IsInf(totalPrice, 0) {
		return errors.New("total harga reservasi tidak valid")
	}

	reservation := Reservation{
		UserID:     userID,
		VehicleID:  vehicleID,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalDays:  totalDays,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	return s.repo.CreateReservation(ctx, reservation)
}

// GetAllReservations mengambil seluruh reservasi.
func (s *reservationService) GetAllReservations(
	ctx context.Context,
) ([]Reservation, error) {
	return s.repo.GetAllReservations(ctx)
}

// GetReservationByID mengambil reservasi berdasarkan ID.
func (s *reservationService) GetReservationByID(
	ctx context.Context,
	id int64,
) (*Reservation, error) {
	if id <= 0 {
		return nil, ErrInvalidReservation
	}

	return s.repo.GetReservationByID(ctx, id)
}

// GetReservationsByUserID mengambil reservasi milik user.
func (s *reservationService) GetReservationsByUserID(
	ctx context.Context,
	userID int64,
) ([]Reservation, error) {
	if userID <= 0 {
		return nil, ErrInvalidReservation
	}

	return s.repo.GetReservationsByUserID(ctx, userID)
}

// UpdateReservationStatus mengubah status reservasi.
func (s *reservationService) UpdateReservationStatus(
	ctx context.Context,
	id int64,
	status string,
) error {
	if id <= 0 {
		return ErrInvalidReservation
	}

	status = strings.ToLower(strings.TrimSpace(status))

	// Untuk saat ini, admin/staff hanya dapat memproses
	// reservasi pending menjadi confirmed atau rejected.
	if status != "confirmed" && status != "rejected" {
		return ErrInvalidStatus
	}

	reservation, err := s.repo.GetReservationByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.Status != "pending" {
		return errors.New("hanya reservasi pending yang dapat diproses")
	}

	return s.repo.UpdateReservationStatus(ctx, id, status)
}

// CancelReservation membatalkan reservasi.
func (s *reservationService) CancelReservation(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrInvalidReservation
	}

	reservation, err := s.repo.GetReservationByID(ctx, id)
	if err != nil {
		return err
	}

	if reservation.Status != "pending" &&
		reservation.Status != "confirmed" {
		return errors.New("reservasi tidak dapat dibatalkan")
	}

	return s.repo.CancelReservation(ctx, id)
}

// dateOnly menghilangkan komponen jam dari tanggal.
func dateOnly(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
}

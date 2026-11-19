package rentals

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidRental      = errors.New("data rental tidak valid")
	ErrInvalidInitialFuel = errors.New("kondisi bahan bakar wajib diisi")
)

type RentalService interface {
	CreateRental(ctx context.Context, rental Rental) error
	GetAllRentals(ctx context.Context) ([]Rental, error)
	GetRentalByID(ctx context.Context, id int64) (*Rental, error)
	GetRentalByReservationID(ctx context.Context, reservationID int64) (*Rental, error)
}

type rentalService struct {
	repo RentalRepository
}

func NewRentalService(repo RentalRepository) RentalService {
	return &rentalService{
		repo: repo,
	}
}

// CreateRental membuat rental baru.
func (s *rentalService) CreateRental(
	ctx context.Context,
	rental Rental,
) error {
	if rental.ReservationID <= 0 || rental.StaffID <= 0 {
		return ErrInvalidRental
	}

	if rental.InitialOdometer < 0 {
		return ErrInvalidRental
	}

	rental.InitialFuel = strings.TrimSpace(rental.InitialFuel)

	if rental.InitialFuel == "" {
		return ErrInvalidInitialFuel
	}

	if len(rental.InitialFuel) > 20 {
		return errors.New("kondisi bahan bakar maksimal 20 karakter")
	}

	// Jika pickup_date tidak diberikan, gunakan waktu sekarang.
	if rental.PickupDate.IsZero() {
		rental.PickupDate = time.Now()
	}

	return s.repo.CreateRental(ctx, rental)
}

// GetAllRentals mengambil seluruh data rental.
func (s *rentalService) GetAllRentals(
	ctx context.Context,
) ([]Rental, error) {
	return s.repo.GetAllRentals(ctx)
}

// GetRentalByID mengambil rental berdasarkan ID.
func (s *rentalService) GetRentalByID(
	ctx context.Context,
	id int64,
) (*Rental, error) {
	if id <= 0 {
		return nil, ErrInvalidRental
	}

	return s.repo.GetRentalByID(ctx, id)
}

// GetRentalByReservationID mengambil rental berdasarkan reservasi.
func (s *rentalService) GetRentalByReservationID(
	ctx context.Context,
	reservationID int64,
) (*Rental, error) {
	if reservationID <= 0 {
		return nil, ErrInvalidRental
	}

	return s.repo.GetRentalByReservationID(ctx, reservationID)
}

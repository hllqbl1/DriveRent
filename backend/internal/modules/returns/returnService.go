package returns

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidReturn        = errors.New("data pengembalian tidak valid")
	ErrInvalidFinalFuel     = errors.New("bahan bakar akhir wajib diisi")
	ErrInvalidVehicleStatus = errors.New("status kendaraan harus available atau maintenance")
)

type ReturnService interface {
	CreateReturn(ctx context.Context, returnData Return, vehicleStatus string) error
	GetAllReturns(ctx context.Context) ([]Return, error)
	GetReturnByID(ctx context.Context, id int64) (*Return, error)
	GetReturnByRentalID(ctx context.Context, rentalID int64) (*Return, error)
}

type returnService struct {
	repo ReturnRepository
}

func NewReturnService(repo ReturnRepository) ReturnService {
	return &returnService{
		repo: repo,
	}
}

// CreateReturn membuat data pengembalian baru.
func (s *returnService) CreateReturn(
	ctx context.Context,
	returnData Return,
	vehicleStatus string,
) error {
	if returnData.RentalID <= 0 {
		return ErrInvalidReturn
	}

	if returnData.FinalOdometer < 0 {
		return ErrInvalidReturn
	}

	if returnData.LateFee < 0 || returnData.DamageFee < 0 {
		return ErrInvalidReturn
	}

	returnData.FinalFuel = strings.TrimSpace(returnData.FinalFuel)

	if returnData.FinalFuel == "" {
		return ErrInvalidFinalFuel
	}

	if len(returnData.FinalFuel) > 20 {
		return errors.New("bahan bakar akhir maksimal 20 karakter")
	}

	if returnData.Notes != nil {
		notes := strings.TrimSpace(*returnData.Notes)

		if notes == "" {
			returnData.Notes = nil
		} else {
			returnData.Notes = &notes
		}
	}

	vehicleStatus = strings.TrimSpace(strings.ToLower(vehicleStatus))

	if vehicleStatus != "available" && vehicleStatus != "maintenance" {
		return ErrInvalidVehicleStatus
	}

	if returnData.ReturnDate.IsZero() {
		returnData.ReturnDate = time.Now()
	}

	return s.repo.CreateReturn(ctx, returnData, vehicleStatus)
}

// GetAllReturns mengambil seluruh data pengembalian.
func (s *returnService) GetAllReturns(
	ctx context.Context,
) ([]Return, error) {
	return s.repo.GetAllReturns(ctx)
}

// GetReturnByID mengambil pengembalian berdasarkan ID.
func (s *returnService) GetReturnByID(
	ctx context.Context,
	id int64,
) (*Return, error) {
	if id <= 0 {
		return nil, ErrInvalidReturn
	}

	return s.repo.GetReturnByID(ctx, id)
}

// GetReturnByRentalID mengambil pengembalian berdasarkan rental.
func (s *returnService) GetReturnByRentalID(
	ctx context.Context,
	rentalID int64,
) (*Return, error) {
	if rentalID <= 0 {
		return nil, ErrInvalidReturn
	}

	return s.repo.GetReturnByRentalID(ctx, rentalID)
}

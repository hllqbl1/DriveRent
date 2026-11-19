package rentals

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrRentalNotFound      = errors.New("rental tidak ditemukan")
	ErrReservationNotFound = errors.New("reservasi tidak ditemukan")
	ErrReservationNotReady = errors.New("reservasi belum dikonfirmasi")
	ErrRentalAlreadyExists = errors.New("rental untuk reservasi ini sudah ada")
	ErrVehicleNotAvailable = errors.New("kendaraan tidak tersedia")
)

type RentalRepository interface {
	CreateRental(ctx context.Context, rental Rental) error
	GetAllRentals(ctx context.Context) ([]Rental, error)
	GetRentalByID(ctx context.Context, id int64) (*Rental, error)
	GetRentalByReservationID(ctx context.Context, reservationID int64) (*Rental, error)
}

type rentalRepository struct {
	db *sql.DB
}

func NewRentalRepository(db *sql.DB) RentalRepository {
	return &rentalRepository{
		db: db,
	}
}

// CreateRental membuat rental baru berdasarkan reservasi yang confirmed.
func (r *rentalRepository) CreateRental(
	ctx context.Context,
	rental Rental,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var vehicleID int64
	var reservationStatus string

	// Lock reservasi agar tidak diproses bersamaan.
	err = tx.QueryRowContext(ctx, `
		SELECT vehicle_id, status
		FROM reservations
		WHERE id = $1
		FOR UPDATE
	`, rental.ReservationID).Scan(
		&vehicleID,
		&reservationStatus,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReservationNotFound
		}
		return err
	}

	if reservationStatus != "confirmed" {
		return ErrReservationNotReady
	}

	// Pastikan belum ada rental untuk reservasi ini.
	var exists bool

	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM rentals
			WHERE reservation_id = $1
		)
	`, rental.ReservationID).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return ErrRentalAlreadyExists
	}

	// Lock kendaraan dan periksa statusnya.
	var vehicleStatus string

	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM vehicles
		WHERE id = $1
		FOR UPDATE
	`, vehicleID).Scan(&vehicleStatus)

	if err != nil {
		return err
	}

	if vehicleStatus != "available" {
		return ErrVehicleNotAvailable
	}

	// Buat rental.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO rentals (
			reservation_id,
			staff_id,
			pickup_date,
			initial_fuel,
			initial_odometer,
			status,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			'active',
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		)
	`,
		rental.ReservationID,
		rental.StaffID,
		rental.PickupDate,
		rental.InitialFuel,
		rental.InitialOdometer,
	)

	if err != nil {
		return err
	}

	// Ubah status kendaraan menjadi rented.
	_, err = tx.ExecContext(ctx, `
		UPDATE vehicles
		SET status = 'rented'
		WHERE id = $1
	`, vehicleID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetAllRentals mengambil seluruh data rental.
func (r *rentalRepository) GetAllRentals(
	ctx context.Context,
) ([]Rental, error) {
	query := `
		SELECT
			id,
			reservation_id,
			staff_id,
			pickup_date,
			initial_fuel,
			initial_odometer,
			status,
			created_at,
			updated_at
		FROM rentals
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rentals := make([]Rental, 0)

	for rows.Next() {
		var rental Rental

		err := rows.Scan(
			&rental.ID,
			&rental.ReservationID,
			&rental.StaffID,
			&rental.PickupDate,
			&rental.InitialFuel,
			&rental.InitialOdometer,
			&rental.Status,
			&rental.CreatedAt,
			&rental.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		rentals = append(rentals, rental)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rentals, nil
}

// GetRentalByID mengambil rental berdasarkan ID.
func (r *rentalRepository) GetRentalByID(
	ctx context.Context,
	id int64,
) (*Rental, error) {
	query := `
		SELECT
			id,
			reservation_id,
			staff_id,
			pickup_date,
			initial_fuel,
			initial_odometer,
			status,
			created_at,
			updated_at
		FROM rentals
		WHERE id = $1
	`

	var rental Rental

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rental.ID,
		&rental.ReservationID,
		&rental.StaffID,
		&rental.PickupDate,
		&rental.InitialFuel,
		&rental.InitialOdometer,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRentalNotFound
		}
		return nil, err
	}

	return &rental, nil
}

// GetRentalByReservationID mengambil rental berdasarkan ID reservasi.
func (r *rentalRepository) GetRentalByReservationID(
	ctx context.Context,
	reservationID int64,
) (*Rental, error) {
	query := `
		SELECT
			id,
			reservation_id,
			staff_id,
			pickup_date,
			initial_fuel,
			initial_odometer,
			status,
			created_at,
			updated_at
		FROM rentals
		WHERE reservation_id = $1
	`

	var rental Rental

	err := r.db.QueryRowContext(ctx, query, reservationID).Scan(
		&rental.ID,
		&rental.ReservationID,
		&rental.StaffID,
		&rental.PickupDate,
		&rental.InitialFuel,
		&rental.InitialOdometer,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRentalNotFound
		}
		return nil, err
	}

	return &rental, nil
}

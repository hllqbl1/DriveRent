package returns

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrReturnNotFound      = errors.New("data pengembalian tidak ditemukan")
	ErrRentalNotFound      = errors.New("rental tidak ditemukan")
	ErrRentalNotActive     = errors.New("rental tidak dalam status active")
	ErrReturnAlreadyExists = errors.New("rental ini sudah memiliki data pengembalian")
	ErrOdometerInvalid     = errors.New("odometer akhir lebih kecil dari odometer awal")
)

type ReturnRepository interface {
	CreateReturn(ctx context.Context, returnData Return, vehicleStatus string) error
	GetAllReturns(ctx context.Context) ([]Return, error)
	GetReturnByID(ctx context.Context, id int64) (*Return, error)
	GetReturnByRentalID(ctx context.Context, rentalID int64) (*Return, error)
}

type returnRepository struct {
	db *sql.DB
}

func NewReturnRepository(db *sql.DB) ReturnRepository {
	return &returnRepository{
		db: db,
	}
}

// CreateReturn membuat data pengembalian dan menyelesaikan rental.
func (r *returnRepository) CreateReturn(
	ctx context.Context,
	returnData Return,
	vehicleStatus string,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var initialOdometer int
	var rentalStatus string
	var vehicleID int64

	// Lock rental dan ambil data kendaraan.
	err = tx.QueryRowContext(ctx, `
		SELECT
			r.initial_odometer,
			r.status,
			res.vehicle_id
		FROM rentals r
		JOIN reservations res
			ON res.id = r.reservation_id
		WHERE r.id = $1
		FOR UPDATE OF r
	`, returnData.RentalID).Scan(
		&initialOdometer,
		&rentalStatus,
		&vehicleID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRentalNotFound
		}
		return err
	}

	if rentalStatus != "active" {
		return ErrRentalNotActive
	}

	// Pastikan odometer akhir tidak lebih kecil dari odometer awal.
	if returnData.FinalOdometer < initialOdometer {
		return ErrOdometerInvalid
	}

	// Lock kendaraan agar perubahan status tidak bertabrakan.
	var currentVehicleStatus string

	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM vehicles
		WHERE id = $1
		FOR UPDATE
	`, vehicleID).Scan(&currentVehicleStatus)

	if err != nil {
		return err
	}

	// Simpan data pengembalian.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO returns (
			rental_id,
			return_date,
			final_fuel,
			final_odometer,
			late_fee,
			damage_fee,
			notes,
			status,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			'completed',
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		)
	`,
		returnData.RentalID,
		returnData.ReturnDate,
		returnData.FinalFuel,
		returnData.FinalOdometer,
		returnData.LateFee,
		returnData.DamageFee,
		returnData.Notes,
	)

	if err != nil {
		return err
	}

	// Selesaikan rental.
	_, err = tx.ExecContext(ctx, `
		UPDATE rentals
		SET status = 'completed',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, returnData.RentalID)

	if err != nil {
		return err
	}

	// Perbarui status kendaraan.
	_, err = tx.ExecContext(ctx, `
		UPDATE vehicles
		SET status = $1
		WHERE id = $2
	`, vehicleStatus, vehicleID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetAllReturns mengambil seluruh data pengembalian.
func (r *returnRepository) GetAllReturns(
	ctx context.Context,
) ([]Return, error) {
	query := `
		SELECT
			id,
			rental_id,
			return_date,
			final_fuel,
			final_odometer,
			late_fee,
			damage_fee,
			notes,
			status,
			created_at,
			updated_at
		FROM returns
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	returns := make([]Return, 0)

	for rows.Next() {
		var returnData Return

		err := rows.Scan(
			&returnData.ID,
			&returnData.RentalID,
			&returnData.ReturnDate,
			&returnData.FinalFuel,
			&returnData.FinalOdometer,
			&returnData.LateFee,
			&returnData.DamageFee,
			&returnData.Notes,
			&returnData.Status,
			&returnData.CreatedAt,
			&returnData.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		returns = append(returns, returnData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return returns, nil
}

// GetReturnByID mengambil data pengembalian berdasarkan ID.
func (r *returnRepository) GetReturnByID(
	ctx context.Context,
	id int64,
) (*Return, error) {
	query := `
		SELECT
			id,
			rental_id,
			return_date,
			final_fuel,
			final_odometer,
			late_fee,
			damage_fee,
			notes,
			status,
			created_at,
			updated_at
		FROM returns
		WHERE id = $1
	`

	var returnData Return

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&returnData.ID,
		&returnData.RentalID,
		&returnData.ReturnDate,
		&returnData.FinalFuel,
		&returnData.FinalOdometer,
		&returnData.LateFee,
		&returnData.DamageFee,
		&returnData.Notes,
		&returnData.Status,
		&returnData.CreatedAt,
		&returnData.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReturnNotFound
		}
		return nil, err
	}

	return &returnData, nil
}

// GetReturnByRentalID mengambil data pengembalian berdasarkan rental.
func (r *returnRepository) GetReturnByRentalID(
	ctx context.Context,
	rentalID int64,
) (*Return, error) {
	query := `
		SELECT
			id,
			rental_id,
			return_date,
			final_fuel,
			final_odometer,
			late_fee,
			damage_fee,
			notes,
			status,
			created_at,
			updated_at
		FROM returns
		WHERE rental_id = $1
	`

	var returnData Return

	err := r.db.QueryRowContext(ctx, query, rentalID).Scan(
		&returnData.ID,
		&returnData.RentalID,
		&returnData.ReturnDate,
		&returnData.FinalFuel,
		&returnData.FinalOdometer,
		&returnData.LateFee,
		&returnData.DamageFee,
		&returnData.Notes,
		&returnData.Status,
		&returnData.CreatedAt,
		&returnData.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReturnNotFound
		}
		return nil, err
	}

	return &returnData, nil
}

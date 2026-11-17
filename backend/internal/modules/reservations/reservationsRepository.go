package reservations

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrReservationNotFound = errors.New("reservasi tidak ditemukan")
	ErrVehicleUnavailable  = errors.New("kendaraan tidak tersedia")
	ErrDateConflict        = errors.New("kendaraan sudah dipesan pada tanggal tersebut")
)

type ReservationRepository interface {
	CreateReservation(ctx context.Context, reservation Reservation) error
	GetAllReservations(ctx context.Context) ([]Reservation, error)
	GetReservationByID(ctx context.Context, id int64) (*Reservation, error)
	GetReservationsByUserID(ctx context.Context, userID int64) ([]Reservation, error)
	UpdateReservationStatus(ctx context.Context, id int64, status string) error
	CancelReservation(ctx context.Context, id int64) error
}

type reservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepository{
		db: db,
	}
}

// CreateReservation membuat reservasi dengan transaksi database.
func (r *reservationRepository) CreateReservation(
	ctx context.Context,
	reservation Reservation,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock kendaraan agar dua transaksi tidak memesan mobil
	// yang sama secara bersamaan.
	var vehicleStatus string

	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM vehicles
		WHERE id = $1
		FOR UPDATE
	`, reservation.VehicleID).Scan(&vehicleStatus)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrVehicleUnavailable
		}
		return err
	}

	if vehicleStatus != "available" {
		return ErrVehicleUnavailable
	}

	// Periksa apakah ada reservasi yang tanggalnya bertabrakan.
	var conflict bool

	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM reservations
			WHERE vehicle_id = $1
			AND status IN ('pending', 'confirmed')
			AND start_date < $3
			AND end_date > $2
		)
	`, reservation.VehicleID, reservation.StartDate, reservation.EndDate).Scan(&conflict)

	if err != nil {
		return err
	}

	if conflict {
		return ErrDateConflict
	}

	// Simpan reservasi.
	query := `
		INSERT INTO reservations (
			user_id,
			vehicle_id,
			start_date,
			end_date,
			total_days,
			total_price,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.ExecContext(ctx, query,
		reservation.UserID,
		reservation.VehicleID,
		reservation.StartDate,
		reservation.EndDate,
		reservation.TotalDays,
		reservation.TotalPrice,
		"pending",
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetAllReservations mengambil seluruh reservasi.
func (r *reservationRepository) GetAllReservations(
	ctx context.Context,
) ([]Reservation, error) {
	query := `
		SELECT
			id, user_id, vehicle_id,
			start_date, end_date,
			total_days, total_price,
			status, created_at, updated_at
		FROM reservations
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation

	for rows.Next() {
		var reservation Reservation

		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.VehicleID,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.TotalDays,
			&reservation.TotalPrice,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationByID mengambil reservasi berdasarkan ID.
func (r *reservationRepository) GetReservationByID(
	ctx context.Context,
	id int64,
) (*Reservation, error) {
	query := `
		SELECT
			id, user_id, vehicle_id,
			start_date, end_date,
			total_days, total_price,
			status, created_at, updated_at
		FROM reservations
		WHERE id = $1
	`

	var reservation Reservation

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.VehicleID,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.TotalDays,
		&reservation.TotalPrice,
		&reservation.Status,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	return &reservation, nil
}

// GetReservationsByUserID mengambil reservasi milik user.
func (r *reservationRepository) GetReservationsByUserID(
	ctx context.Context,
	userID int64,
) ([]Reservation, error) {
	query := `
		SELECT
			id, user_id, vehicle_id,
			start_date, end_date,
			total_days, total_price,
			status, created_at, updated_at
		FROM reservations
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation

	for rows.Next() {
		var reservation Reservation

		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.VehicleID,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.TotalDays,
			&reservation.TotalPrice,
			&reservation.Status,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// UpdateReservationStatus mengubah status reservasi.
func (r *reservationRepository) UpdateReservationStatus(
	ctx context.Context,
	id int64,
	status string,
) error {
	query := `
		UPDATE reservations
		SET status = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrReservationNotFound
	}

	return nil
}

// CancelReservation membatalkan reservasi.
func (r *reservationRepository) CancelReservation(
	ctx context.Context,
	id int64,
) error {
	query := `
		UPDATE reservations
		SET status = 'cancelled',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		  AND status IN ('pending', 'confirmed')
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Bedakan reservasi yang tidak ada dengan status
		// yang memang tidak bisa dibatalkan.
		var exists bool

		err = r.db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM reservations WHERE id = $1
			)
		`, id).Scan(&exists)

		if err != nil {
			return err
		}

		if !exists {
			return ErrReservationNotFound
		}

		return errors.New("reservasi tidak dapat dibatalkan")
	}

	return nil
}

// Menjaga import time tetap digunakan jika nanti diperlukan
// untuk query atau helper tanggal pada repository ini.
var _ = time.Time{}

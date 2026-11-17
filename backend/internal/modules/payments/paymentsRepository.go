package payments

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrPaymentNotFound         = errors.New("pembayaran tidak ditemukan")
	ErrReservationNotFound     = errors.New("reservasi tidak ditemukan")
	ErrReservationNotPending   = errors.New("reservasi tidak dalam status pending")
	ErrPaymentAlreadyExists    = errors.New("masih ada pembayaran aktif untuk reservasi ini")
	ErrPaymentAlreadyProcessed = errors.New("pembayaran sudah diproses")
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, reservationID int64, method string, proof *string) error
	GetAllPayments(ctx context.Context) ([]Payment, error)
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	GetPaymentsByReservationID(ctx context.Context, reservationID int64) ([]Payment, error)
	VerifyPayment(ctx context.Context, id int64) error
	RejectPayment(ctx context.Context, id int64) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

// CreatePayment membuat pembayaran berdasarkan total harga reservasi.
func (r *paymentRepository) CreatePayment(
	ctx context.Context,
	reservationID int64,
	method string,
	proof *string,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var amount float64
	var reservationStatus string

	// Lock reservasi untuk mencegah pembuatan pembayaran bersamaan.
	err = tx.QueryRowContext(ctx, `
		SELECT total_price, status
		FROM reservations
		WHERE id = $1
		FOR UPDATE
	`, reservationID).Scan(&amount, &reservationStatus)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReservationNotFound
		}
		return err
	}

	if reservationStatus != "pending" {
		return ErrReservationNotPending
	}

	// Periksa apakah sudah ada pembayaran pending atau verified.
	var exists bool

	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM payments
			WHERE reservation_id = $1
			AND status IN ('pending', 'verified')
		)
	`, reservationID).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return ErrPaymentAlreadyExists
	}

	// Simpan pembayaran.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO payments (
			reservation_id,
			amount,
			payment_method,
			payment_proof,
			status,
			paid_at,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4,
			'pending',
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		)
	`, reservationID, amount, method, proof)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetAllPayments mengambil seluruh pembayaran.
func (r *paymentRepository) GetAllPayments(
	ctx context.Context,
) ([]Payment, error) {
	query := `
		SELECT
			id,
			reservation_id,
			amount,
			payment_method,
			payment_proof,
			status,
			paid_at,
			created_at,
			updated_at
		FROM payments
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]Payment, 0)

	for rows.Next() {
		var payment Payment

		err := rows.Scan(
			&payment.ID,
			&payment.ReservationID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentProof,
			&payment.Status,
			&payment.PaidAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetPaymentByID mengambil pembayaran berdasarkan ID.
func (r *paymentRepository) GetPaymentByID(
	ctx context.Context,
	id int64,
) (*Payment, error) {
	query := `
		SELECT
			id,
			reservation_id,
			amount,
			payment_method,
			payment_proof,
			status,
			paid_at,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`

	var payment Payment

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.ReservationID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentProof,
		&payment.Status,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

// GetPaymentsByReservationID mengambil riwayat pembayaran reservasi.
func (r *paymentRepository) GetPaymentsByReservationID(
	ctx context.Context,
	reservationID int64,
) ([]Payment, error) {
	query := `
		SELECT
			id,
			reservation_id,
			amount,
			payment_method,
			payment_proof,
			status,
			paid_at,
			created_at,
			updated_at
		FROM payments
		WHERE reservation_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]Payment, 0)

	for rows.Next() {
		var payment Payment

		err := rows.Scan(
			&payment.ID,
			&payment.ReservationID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentProof,
			&payment.Status,
			&payment.PaidAt,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// VerifyPayment memverifikasi pembayaran dan mengonfirmasi reservasi.
func (r *paymentRepository) VerifyPayment(
	ctx context.Context,
	id int64,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var reservationID int64
	var paymentStatus string

	// Lock pembayaran.
	err = tx.QueryRowContext(ctx, `
		SELECT reservation_id, status
		FROM payments
		WHERE id = $1
		FOR UPDATE
	`, id).Scan(&reservationID, &paymentStatus)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPaymentNotFound
		}
		return err
	}

	if paymentStatus != "pending" {
		return ErrPaymentAlreadyProcessed
	}

	var reservationStatus string

	// Lock reservasi.
	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM reservations
		WHERE id = $1
		FOR UPDATE
	`, reservationID).Scan(&reservationStatus)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrReservationNotFound
		}
		return err
	}

	if reservationStatus != "pending" {
		return ErrReservationNotPending
	}

	// Update status pembayaran.
	_, err = tx.ExecContext(ctx, `
		UPDATE payments
		SET status = 'verified',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	// Update status reservasi.
	_, err = tx.ExecContext(ctx, `
		UPDATE reservations
		SET status = 'confirmed',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, reservationID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// RejectPayment menolak pembayaran.
func (r *paymentRepository) RejectPayment(
	ctx context.Context,
	id int64,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string

	// Lock pembayaran agar tidak diproses bersamaan.
	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM payments
		WHERE id = $1
		FOR UPDATE
	`, id).Scan(&status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPaymentNotFound
		}
		return err
	}

	if status != "pending" {
		return ErrPaymentAlreadyProcessed
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE payments
		SET status = 'rejected',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

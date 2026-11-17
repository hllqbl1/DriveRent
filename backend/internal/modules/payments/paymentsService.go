package payments

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidPayment = errors.New("data pembayaran tidak valid")
	ErrInvalidMethod  = errors.New("metode pembayaran wajib diisi")
)

type PaymentService interface {
	CreatePayment(ctx context.Context, reservationID int64, method string, proof *string) error
	GetAllPayments(ctx context.Context) ([]Payment, error)
	GetPaymentByID(ctx context.Context, id int64) (*Payment, error)
	GetPaymentsByReservationID(ctx context.Context, reservationID int64) ([]Payment, error)
	VerifyPayment(ctx context.Context, id int64) error
	RejectPayment(ctx context.Context, id int64) error
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}

// CreatePayment membuat pembayaran baru.
func (s *paymentService) CreatePayment(
	ctx context.Context,
	reservationID int64,
	method string,
	proof *string,
) error {
	if reservationID <= 0 {
		return ErrInvalidPayment
	}

	method = strings.TrimSpace(method)

	if method == "" {
		return ErrInvalidMethod
	}

	if len(method) > 30 {
		return errors.New("metode pembayaran maksimal 30 karakter")
	}

	// Bersihkan URL bukti pembayaran.
	if proof != nil {
		cleanProof := strings.TrimSpace(*proof)

		if cleanProof == "" {
			proof = nil
		} else {
			proof = &cleanProof
		}
	}

	return s.repo.CreatePayment(
		ctx,
		reservationID,
		method,
		proof,
	)
}

// GetAllPayments mengambil seluruh pembayaran.
func (s *paymentService) GetAllPayments(
	ctx context.Context,
) ([]Payment, error) {
	return s.repo.GetAllPayments(ctx)
}

// GetPaymentByID mengambil pembayaran berdasarkan ID.
func (s *paymentService) GetPaymentByID(
	ctx context.Context,
	id int64,
) (*Payment, error) {
	if id <= 0 {
		return nil, ErrInvalidPayment
	}

	return s.repo.GetPaymentByID(ctx, id)
}

// GetPaymentsByReservationID mengambil riwayat pembayaran reservasi.
func (s *paymentService) GetPaymentsByReservationID(
	ctx context.Context,
	reservationID int64,
) ([]Payment, error) {
	if reservationID <= 0 {
		return nil, ErrInvalidPayment
	}

	return s.repo.GetPaymentsByReservationID(
		ctx,
		reservationID,
	)
}

// VerifyPayment memverifikasi pembayaran.
func (s *paymentService) VerifyPayment(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrInvalidPayment
	}

	return s.repo.VerifyPayment(ctx, id)
}

// RejectPayment menolak pembayaran.
func (s *paymentService) RejectPayment(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrInvalidPayment
	}

	return s.repo.RejectPayment(ctx, id)
}

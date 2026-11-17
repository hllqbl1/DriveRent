package payments

import "time"

type Payment struct {
	ID            int64      `json:"id"`
	ReservationID int64      `json:"reservation_id"`
	Amount        float64    `json:"amount"`
	PaymentMethod string     `json:"payment_method"`
	PaymentProof  *string    `json:"payment_proof"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

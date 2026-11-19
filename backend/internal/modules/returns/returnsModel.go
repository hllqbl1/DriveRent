package returns

import "time"

type Return struct {
	ID            int64     `json:"id"`
	RentalID      int64     `json:"rental_id"`
	ReturnDate    time.Time `json:"return_date"`
	FinalFuel     string    `json:"final_fuel"`
	FinalOdometer int       `json:"final_odometer"`
	LateFee       float64   `json:"late_fee"`
	DamageFee     float64   `json:"damage_fee"`
	Notes         *string   `json:"notes"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

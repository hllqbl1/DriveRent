package rentals

import "time"

type Rental struct {
	ID              int64     `json:"id"`
	ReservationID   int64     `json:"reservation_id"`
	StaffID         int64     `json:"staff_id"`
	PickupDate      time.Time `json:"pickup_date"`
	InitialFuel     string    `json:"initial_fuel"`
	InitialOdometer int       `json:"initial_odometer"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

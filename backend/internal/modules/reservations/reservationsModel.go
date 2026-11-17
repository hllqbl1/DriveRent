package reservations

import "time"

type Reservation struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	VehicleID  int64     `json:"vehicle_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalDays  int       `json:"total_days"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

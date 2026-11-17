package vehicles

import "time"

type Vehicle struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	LicensePlate string    `json:"license_plate"`
	Category     string    `json:"category"`
	Transmission string    `json:"transmission"`
	FuelType     string    `json:"fuel_type"`
	Capacity     int       `json:"capacity"`
	PricePerDay  float64   `json:"price_per_day"`
	Image        *string   `json:"image"`
	Description  *string   `json:"description"`
	Status       string    `json:"status"`
	Created_at   time.Time `json:"created_at"`
	Update_at    time.Time `json:"updated_at"`
}

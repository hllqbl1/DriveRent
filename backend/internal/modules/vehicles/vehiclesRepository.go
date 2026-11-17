package vehicles

import (
	"context"
	"database/sql"
)

type VehicleRepository interface {
	CreateVehicle(ctx context.Context, vehicle Vehicle) error
	GetAllVehicles(ctx context.Context) ([]Vehicle, error)
	GetVehicleByID(ctx context.Context, id int64) (*Vehicle, error)
	UpdateVehicle(ctx context.Context, vehicle Vehicle) error
	UpdateVehicleStatus(ctx context.Context, id int64, status string) error
	DeleteVehicle(ctx context.Context, id int64) error
	IsLicensePlateUsed(ctx context.Context, licensePlate string, excludeID int64) (bool, error)
}

type vehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) VehicleRepository {
	return &vehicleRepository{
		db: db,
	}
}

// CreateVehicle menambahkan kendaraan ke database.
func (r *vehicleRepository) CreateVehicle(ctx context.Context, vehicle Vehicle) error {
	query := `
		INSERT INTO vehicles (
			name, brand, model, year, license_plate,
			category, transmission, fuel_type, capacity,
			price_per_day, image, description, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		vehicle.Name,
		vehicle.Brand,
		vehicle.Model,
		vehicle.Year,
		vehicle.LicensePlate,
		vehicle.Category,
		vehicle.Transmission,
		vehicle.FuelType,
		vehicle.Capacity,
		vehicle.PricePerDay,
		vehicle.Image,
		vehicle.Description,
		vehicle.Status,
	)

	return err
}

// GetAllVehicles mengambil seluruh kendaraan.
func (r *vehicleRepository) GetAllVehicles(ctx context.Context) ([]Vehicle, error) {
	query := `
		SELECT
			id, name, brand, model, year, license_plate,
			category, transmission, fuel_type, capacity,
			price_per_day, image, description, status, created_at
		FROM vehicles
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle

	for rows.Next() {
		var vehicle Vehicle

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.Name,
			&vehicle.Brand,
			&vehicle.Model,
			&vehicle.Year,
			&vehicle.LicensePlate,
			&vehicle.Category,
			&vehicle.Transmission,
			&vehicle.FuelType,
			&vehicle.Capacity,
			&vehicle.PricePerDay,
			&vehicle.Image,
			&vehicle.Description,
			&vehicle.Status,
			&vehicle.Created_at,
		)

		if err != nil {
			return nil, err
		}

		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

// GetVehicleByID mengambil kendaraan berdasarkan ID.
func (r *vehicleRepository) GetVehicleByID(ctx context.Context, id int64) (*Vehicle, error) {
	query := `
		SELECT
			id, name, brand, model, year, license_plate,
			category, transmission, fuel_type, capacity,
			price_per_day, image, description, status, created_at
		FROM vehicles
		WHERE id = $1
	`

	var vehicle Vehicle

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vehicle.ID,
		&vehicle.Name,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.LicensePlate,
		&vehicle.Category,
		&vehicle.Transmission,
		&vehicle.FuelType,
		&vehicle.Capacity,
		&vehicle.PricePerDay,
		&vehicle.Image,
		&vehicle.Description,
		&vehicle.Status,
		&vehicle.Created_at,
	)

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}

// UpdateVehicle memperbarui informasi kendaraan selain status.
func (r *vehicleRepository) UpdateVehicle(ctx context.Context, vehicle Vehicle) error {
	query := `
		UPDATE vehicles SET
			name = $1,
			brand = $2,
			model = $3,
			year = $4,
			license_plate = $5,
			category = $6,
			transmission = $7,
			fuel_type = $8,
			capacity = $9,
			price_per_day = $10,
			image = $11,
			description = $12
		WHERE id = $13
	`

	result, err := r.db.ExecContext(ctx, query,
		vehicle.Name,
		vehicle.Brand,
		vehicle.Model,
		vehicle.Year,
		vehicle.LicensePlate,
		vehicle.Category,
		vehicle.Transmission,
		vehicle.FuelType,
		vehicle.Capacity,
		vehicle.PricePerDay,
		vehicle.Image,
		vehicle.Description,
		vehicle.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateVehicleStatus mengubah status kendaraan.
func (r *vehicleRepository) UpdateVehicleStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE vehicles
		SET status = $1
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
		return sql.ErrNoRows
	}

	return nil
}

// DeleteVehicle menghapus kendaraan berdasarkan ID.
func (r *vehicleRepository) DeleteVehicle(ctx context.Context, id int64) error {
	query := `DELETE FROM vehicles WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// IsLicensePlateUsed memeriksa apakah plat nomor sudah digunakan.
func (r *vehicleRepository) IsLicensePlateUsed(
	ctx context.Context,
	licensePlate string,
	excludeID int64,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM vehicles
			WHERE license_plate = $1
			AND id != $2
		)
	`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, licensePlate, excludeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

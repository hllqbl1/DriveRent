CREATE TABLE vehicles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INTEGER NOT NULL,
    license_plate VARCHAR(20) NOT NULL UNIQUE,
    category VARCHAR(50) NOT NULL,
    transmission VARCHAR(20) NOT NULL,
    fuel_type VARCHAR(20) NOT NULL,
    capacity INTEGER NOT NULL,
    price_per_day NUMERIC(12, 2) NOT NULL,
    image TEXT,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_vehicles_transmission
        CHECK (transmission IN ('manual', 'automatic')),
    CONSTRAINT check_vehicles_capacity
        CHECK (capacity > 0),
    CONSTRAINT check_vehicles_price
        CHECK (price_per_day > 0),
    CONSTRAINT check_vehicles_status
        CHECK (status IN ('available', 'rented', 'maintenance'))
);

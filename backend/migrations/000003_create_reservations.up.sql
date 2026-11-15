CREATE TABLE reservations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    vehicle_id BIGINT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INTEGER NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_reservations_user
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_reservations_vehicle
        FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
    CONSTRAINT check_reservations_dates
        CHECK (end_date > start_date),
    CONSTRAINT check_reservations_total_days
        CHECK (total_days > 0),
    CONSTRAINT check_reservations_total_price
        CHECK (total_price > 0),
    CONSTRAINT check_reservations_status
        CHECK (status IN ('pending', 'confirmed', 'rejected', 'cancelled'))
);

CREATE TABLE returns (
    id BIGSERIAL PRIMARY KEY,
    rental_id BIGINT NOT NULL UNIQUE,
    return_date TIMESTAMP NOT NULL,
    final_fuel VARCHAR(20) NOT NULL,
    final_odometer INTEGER NOT NULL,
    late_fee NUMERIC(12, 2) NOT NULL DEFAULT 0,
    damage_fee NUMERIC(12, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_returns_rental
        FOREIGN KEY (rental_id) REFERENCES rentals(id),
    CONSTRAINT check_returns_final_odometer
        CHECK (final_odometer >= 0),
    CONSTRAINT check_returns_late_fee
        CHECK (late_fee >= 0),
    CONSTRAINT check_returns_damage_fee
        CHECK (damage_fee >= 0),
    CONSTRAINT check_returns_status
        CHECK (status IN ('completed'))
);

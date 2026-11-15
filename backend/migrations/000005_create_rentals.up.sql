CREATE TABLE rentals (
    id BIGSERIAL PRIMARY KEY,
    reservation_id BIGINT NOT NULL UNIQUE,
    staff_id BIGINT NOT NULL,
    pickup_date TIMESTAMP NOT NULL,
    initial_fuel VARCHAR(20) NOT NULL,
    initial_odometer INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_rentals_reservation
        FOREIGN KEY (reservation_id) REFERENCES reservations(id),
    CONSTRAINT fk_rentals_staff
        FOREIGN KEY (staff_id) REFERENCES users(id),
    CONSTRAINT check_rentals_initial_odometer
        CHECK (initial_odometer >= 0),
    CONSTRAINT check_rentals_status
        CHECK (status IN ('active', 'completed'))
);

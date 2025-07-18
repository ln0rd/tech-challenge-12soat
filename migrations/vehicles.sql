CREATE TABLE vehicles (
    id SERIAL PRIMARY KEY,
    model VARCHAR NOT NULL,
    brand VARCHAR NOT NULL,
    release_year INTEGER NOT NULL,
    vehicle_identification_number VARCHAR NOT NULL,
    number_plate VARCHAR NOT NULL,
    color VARCHAR NOT NULL,
    customer_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model VARCHAR NOT NULL,
    brand VARCHAR NOT NULL,
    release_year INTEGER NOT NULL,
    vehicle_identification_number VARCHAR NOT NULL,
    number_plate VARCHAR NOT NULL UNIQUE,
    color VARCHAR NOT NULL,
    customer_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE order_supplies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    supply_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    total_value NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
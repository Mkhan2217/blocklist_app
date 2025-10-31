-- Database schema for CheckGuard blocked_numbers table.
-- The phone_number CHECK enforces international format starting with '+'
-- followed by 10-15 digits and first digit after + must not be 0.
CREATE TABLE IF NOT EXISTS blocked_numbers (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(15) UNIQUE NOT NULL CHECK (phone_number ~ '^\+[1-9][0-9]{9,14}$'),
    reason VARCHAR(100) NOT NULL,
    store_location VARCHAR(100) NOT NULL,
    incident_date DATE NOT NULL DEFAULT CURRENT_DATE,
    check_amount NUMERIC(12,2),
    notes VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_phone_number ON blocked_numbers(phone_number);
CREATE INDEX IF NOT EXISTS idx_store_location ON blocked_numbers(store_location);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create patients table
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL,
    gender VARCHAR(50) NOT NULL,
    contact_info VARCHAR(255) NOT NULL,
    medical_notes TEXT,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on patients name for faster search
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(name);

-- Create index on patients gender for faster filtering
CREATE INDEX IF NOT EXISTS idx_patients_gender ON patients(gender);

-- Create index on patients age for faster filtering
CREATE INDEX IF NOT EXISTS idx_patients_age ON patients(age);

-- Insert default admin user (password: admin123)
INSERT INTO users (username, password_hash, role)
VALUES ('admin', '$2a$10$1qAz2wSx3eDc4rFv5tGb5edDmJMuYFJx4hQ/g8MgbxEP6Y.M5uy7y', 'receptionist')
ON CONFLICT (username) DO NOTHING;

-- Insert default doctor user (password: doctor123)
INSERT INTO users (username, password_hash, role)
VALUES ('doctor', '$2a$10$37OXvIOPu7vEPCazGczL3.Qm.7tXUUayL1.7F/hYU.5LOOFBK9.Oa', 'doctor')
ON CONFLICT (username) DO NOTHING;

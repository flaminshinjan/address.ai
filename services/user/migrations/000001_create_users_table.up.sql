CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'guest',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create admin user with password 'admin123'
INSERT INTO users (id, username, email, password, first_name, last_name, role, created_at, updated_at)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'admin',
    'admin@hotel.com',
    '$2a$10$1/dXSLDDyI1bJ9X8/2HE9eF.fEqBQeJtjlY7y99PGR8tEbGi1QpLm', -- hashed 'admin123'
    'Admin',
    'User',
    'admin',
    NOW(),
    NOW()
) ON CONFLICT DO NOTHING; 
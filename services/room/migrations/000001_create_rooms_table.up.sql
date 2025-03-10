CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(36) PRIMARY KEY,
    number VARCHAR(20) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    floor INT NOT NULL,
    description TEXT,
    capacity INT NOT NULL,
    price_per_day DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create some initial rooms
INSERT INTO rooms (id, number, type, floor, description, capacity, price_per_day, status, created_at, updated_at)
VALUES 
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '101', 'Standard', 1, 'Standard room with a queen bed', 2, 100.00, 'available', NOW(), NOW()),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '102', 'Standard', 1, 'Standard room with two twin beds', 2, 100.00, 'available', NOW(), NOW()),
    ('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', '201', 'Deluxe', 2, 'Deluxe room with a king bed and city view', 2, 150.00, 'available', NOW(), NOW()),
    ('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', '202', 'Deluxe', 2, 'Deluxe room with a king bed and ocean view', 2, 180.00, 'available', NOW(), NOW()),
    ('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', '301', 'Suite', 3, 'Suite with a king bed, living room, and kitchenette', 4, 250.00, 'available', NOW(), NOW()),
    ('b5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', '302', 'Suite', 3, 'Suite with a king bed, living room, and ocean view', 4, 300.00, 'available', NOW(), NOW()),
    ('b6eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', '401', 'Presidential Suite', 4, 'Presidential suite with two bedrooms, living room, dining room, and panoramic view', 6, 500.00, 'available', NOW(), NOW())
ON CONFLICT DO NOTHING; 
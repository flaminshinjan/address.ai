CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_number VARCHAR(10) NOT NULL UNIQUE,
    room_type VARCHAR(50) NOT NULL,
    description TEXT,
    price_per_night DECIMAL(10,2) NOT NULL,
    capacity INT NOT NULL,
    amenities TEXT[],
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert sample rooms
INSERT INTO rooms (id, room_number, room_type, description, price_per_night, capacity, amenities, status)
VALUES 
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '101', 'Deluxe King', 'Spacious room with king-size bed and city view', 200.00, 2, ARRAY['King Bed', 'City View', 'Mini Bar', 'Wi-Fi', 'TV'], 'available'),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '102', 'Twin Deluxe', 'Comfortable room with two twin beds', 180.00, 2, ARRAY['Twin Beds', 'Garden View', 'Wi-Fi', 'TV'], 'available'),
    ('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', '201', 'Executive Suite', 'Luxury suite with separate living area', 350.00, 3, ARRAY['King Bed', 'Living Room', 'Mini Bar', 'Wi-Fi', 'TV', 'Bathtub'], 'available'),
    ('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', '202', 'Family Suite', 'Perfect for families with extra space', 400.00, 4, ARRAY['King Bed', 'Sofa Bed', 'Kitchen', 'Wi-Fi', 'TV', 'Bathtub'], 'available'),
    ('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', '301', 'Presidential Suite', 'Our finest suite with premium amenities', 600.00, 4, ARRAY['King Bed', 'Living Room', 'Kitchen', 'Mini Bar', 'Wi-Fi', 'TV', 'Jacuzzi'], 'available')
ON CONFLICT DO NOTHING; 
CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create some initial suppliers
INSERT INTO suppliers (id, name, email, phone, address, description, is_active)
VALUES 
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Global Foods Supply Co.', 'contact@globalfoods.com', '+1-555-0123', '123 Supply St, Business District', 'Major food ingredients supplier', true),
    ('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Hotel Essentials Ltd.', 'orders@hotelessentials.com', '+1-555-0124', '456 Hotel Ave, Commerce Park', 'Hotel amenities and supplies', true),
    ('d2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Fresh Produce Inc.', 'sales@freshproduce.com', '+1-555-0125', '789 Fresh Rd, Market District', 'Fresh fruits and vegetables supplier', true),
    ('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'Cleaning Solutions Co.', 'info@cleaningsolutions.com', '+1-555-0126', '321 Clean St, Industrial Park', 'Cleaning supplies and equipment', true)
ON CONFLICT DO NOTHING; 
CREATE TABLE IF NOT EXISTS inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    minimum_quantity INTEGER NOT NULL DEFAULT 0,
    supplier_id UUID REFERENCES suppliers(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create some initial inventory items
INSERT INTO inventory_items (id, name, description, category, unit, quantity, minimum_quantity, supplier_id)
VALUES 
    ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Bath Towels', 'Luxury cotton bath towels', 'Linens', 'piece', 200, 50, 'd1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'),
    ('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Hand Soap', 'Antibacterial hand soap', 'Toiletries', 'bottle', 150, 30, 'd1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'),
    ('e2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Toilet Paper', 'Premium toilet paper rolls', 'Supplies', 'roll', 500, 100, 'd1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'),
    ('e3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'All-Purpose Cleaner', 'Multi-surface cleaning solution', 'Cleaning', 'bottle', 80, 20, 'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14'),
    ('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'Flour', 'All-purpose flour', 'Food', 'kg', 100, 25, 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
    ('e5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'Sugar', 'Granulated sugar', 'Food', 'kg', 80, 20, 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
    ('e6eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'Tomatoes', 'Fresh tomatoes', 'Produce', 'kg', 50, 15, 'd2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13'),
    ('e7eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'Lettuce', 'Fresh lettuce', 'Produce', 'kg', 30, 10, 'd2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13')
ON CONFLICT DO NOTHING; 
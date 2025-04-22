CREATE TABLE IF NOT EXISTS menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    preparation_time INT NOT NULL, -- in minutes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert sample menu items
INSERT INTO menu_items (id, name, description, category, price, preparation_time, is_available)
VALUES 
    ('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Classic Burger', 'Juicy beef patty with lettuce, tomato, and cheese', 'Main Course', 18.99, 20, true),
    ('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Caesar Salad', 'Crisp romaine lettuce with parmesan and croutons', 'Starters', 12.99, 10, true),
    ('f2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Margherita Pizza', 'Fresh tomatoes, mozzarella, and basil', 'Main Course', 16.99, 25, true),
    ('f3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'Chocolate Brownie', 'Warm chocolate brownie with vanilla ice cream', 'Dessert', 8.99, 15, true),
    ('f4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'Club Sandwich', 'Triple-decker with chicken, bacon, and egg', 'Main Course', 15.99, 15, true),
    ('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'Fruit Platter', 'Selection of fresh seasonal fruits', 'Dessert', 10.99, 10, true),
    ('f6eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'Tomato Soup', 'Creamy tomato soup with garlic bread', 'Starters', 9.99, 15, true),
    ('f7eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'Grilled Salmon', 'Fresh salmon with seasonal vegetables', 'Main Course', 24.99, 25, true)
ON CONFLICT DO NOTHING; 
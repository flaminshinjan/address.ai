CREATE TABLE IF NOT EXISTS menu_items (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create some initial menu items
INSERT INTO menu_items (id, name, description, category, price, is_available, created_at, updated_at)
VALUES 
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Margherita Pizza', 'Classic pizza with tomato sauce, mozzarella, and basil', 'Pizza', 12.99, true, NOW(), NOW()),
    ('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Pepperoni Pizza', 'Pizza with tomato sauce, mozzarella, and pepperoni', 'Pizza', 14.99, true, NOW(), NOW()),
    ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Caesar Salad', 'Romaine lettuce, croutons, parmesan cheese, and Caesar dressing', 'Salad', 8.99, true, NOW(), NOW()),
    ('c3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'Greek Salad', 'Tomatoes, cucumbers, olives, feta cheese, and olive oil', 'Salad', 9.99, true, NOW(), NOW()),
    ('c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'Cheeseburger', 'Beef patty, cheddar cheese, lettuce, tomato, and special sauce', 'Burger', 11.99, true, NOW(), NOW()),
    ('c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'Veggie Burger', 'Plant-based patty, lettuce, tomato, and special sauce', 'Burger', 10.99, true, NOW(), NOW()),
    ('c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', 'Chocolate Cake', 'Rich chocolate cake with chocolate ganache', 'Dessert', 6.99, true, NOW(), NOW()),
    ('c7eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', 'Cheesecake', 'New York style cheesecake with berry compote', 'Dessert', 7.99, true, NOW(), NOW()),
    ('c8eebc99-9c0b-4ef8-bb6d-6bb9bd380a19', 'Coca-Cola', 'Classic Coca-Cola', 'Beverage', 2.99, true, NOW(), NOW()),
    ('c9eebc99-9c0b-4ef8-bb6d-6bb9bd380a20', 'Iced Tea', 'Freshly brewed iced tea', 'Beverage', 2.99, true, NOW(), NOW())
ON CONFLICT DO NOTHING; 
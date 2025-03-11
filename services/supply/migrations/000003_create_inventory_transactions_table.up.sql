CREATE TABLE IF NOT EXISTS inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_item_id UUID NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'in' or 'out'
    source VARCHAR(100) NOT NULL, -- 'purchase', 'adjustment', 'consumption', etc.
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
); 
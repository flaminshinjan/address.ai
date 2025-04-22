-- Enable Row Level Security
ALTER TABLE rooms ENABLE ROW LEVEL SECURITY;
ALTER TABLE bookings ENABLE ROW LEVEL SECURITY;

-- Create policies for rooms
CREATE POLICY "Allow public read access" ON rooms
  FOR SELECT USING (true);

CREATE POLICY "Allow authenticated users to book rooms" ON rooms
  FOR UPDATE USING (auth.role() = 'authenticated')
  WITH CHECK (status = 'available');

-- Create policies for bookings
CREATE POLICY "Users can view their own bookings" ON bookings
  FOR SELECT USING (auth.uid() = user_id::uuid);

CREATE POLICY "Users can create their own bookings" ON bookings
  FOR INSERT WITH CHECK (auth.uid() = user_id::uuid);

-- Create a view that joins bookings with rooms
CREATE VIEW user_bookings AS
SELECT 
  b.*,
  r.room_number,
  r.room_type,
  r.description,
  r.amenities
FROM bookings b
JOIN rooms r ON b.room_id = r.id;

-- Enable RLS on the view
ALTER VIEW user_bookings ENABLE ROW LEVEL SECURITY;

-- Create policy for the view
CREATE POLICY "Users can view their own booking details" ON user_bookings
  FOR SELECT USING (auth.uid() = user_id::uuid); 
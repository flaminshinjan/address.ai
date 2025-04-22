import React, { useState, useEffect } from 'react';
import {
  Container,
  Grid,
  Card,
  CardContent,
  CardMedia,
  Typography,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Box,
  Chip,
  Alert,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { roomApi } from '../services/api';
import { format } from 'date-fns';

interface Room {
  id: string;
  room_number: string;
  room_type: string;
  description: string;
  price_per_night: number;
  capacity: number;
  amenities: string[];
  status: string;
}

const Rooms = () => {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);
  const [checkIn, setCheckIn] = useState<Date | null>(null);
  const [checkOut, setCheckOut] = useState<Date | null>(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [bookingSuccess, setBookingSuccess] = useState(false);

  useEffect(() => {
    fetchRooms();
  }, []);

  const fetchRooms = async () => {
    try {
      const response = await roomApi.get('/rooms?select=*');
      setRooms(response.data);
      setLoading(false);
    } catch (err) {
      setError('Failed to load rooms');
      setLoading(false);
      console.error('Error fetching rooms:', err);
    }
  };

  const handleBookRoom = async () => {
    if (!selectedRoom || !checkIn || !checkOut) return;

    try {
      const nights = Math.ceil((checkOut.getTime() - checkIn.getTime()) / (1000 * 60 * 60 * 24));
      const totalPrice = selectedRoom.price_per_night * nights;

      // Create booking
      await roomApi.post('/bookings', {
        room_id: selectedRoom.id,
        user_id: localStorage.getItem('userId'), // Make sure this is set during login
        check_in_date: format(checkIn, 'yyyy-MM-dd'),
        check_out_date: format(checkOut, 'yyyy-MM-dd'),
        total_price: totalPrice,
        status: 'confirmed'
      });

      // Update room status
      await roomApi.patch(`/rooms?id=eq.${selectedRoom.id}`, {
        status: 'booked'
      });

      setBookingSuccess(true);
      setTimeout(() => {
        setOpenDialog(false);
        setBookingSuccess(false);
        fetchRooms(); // Refresh rooms list
      }, 2000);
    } catch (err) {
      setError('Failed to book room');
      console.error('Error booking room:', err);
    }
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Alert severity="error">{error}</Alert>;

  return (
    <Container>
      <Typography variant="h4" gutterBottom sx={{ my: 4 }}>
        Available Rooms
      </Typography>
      <Grid container spacing={4}>
        {rooms.map((room) => (
          <Grid item xs={12} sm={6} md={4} key={room.id}>
            <Card>
              <CardMedia
                component="img"
                height="200"
                image={`/room-${room.room_type.toLowerCase().replace(/\s+/g, '-')}.jpg`}
                alt={room.room_type}
              />
              <CardContent>
                  <Typography variant="h6" gutterBottom>
                  {room.room_type} - Room {room.room_number}
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  {room.description}
                </Typography>
                <Box sx={{ mb: 2 }}>
                  {room.amenities.map((amenity) => (
                    <Chip
                      key={amenity}
                      label={amenity}
                      size="small"
                      sx={{ mr: 0.5, mb: 0.5 }}
                    />
                  ))}
                </Box>
                <Typography variant="h6" color="primary" gutterBottom>
                  ${room.price_per_night} per night
                </Typography>
                  <Button
                    variant="contained"
                    fullWidth
                  disabled={room.status !== 'available'}
                  onClick={() => {
                    setSelectedRoom(room);
                    setOpenDialog(true);
                  }}
                  >
                  {room.status === 'available' ? 'Book Now' : 'Not Available'}
                  </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
        <DialogTitle>Book Room</DialogTitle>
        <DialogContent>
          {bookingSuccess ? (
            <Alert severity="success">Booking successful!</Alert>
          ) : (
            <>
              <Typography variant="body1" gutterBottom>
                {selectedRoom?.room_type} - Room {selectedRoom?.room_number}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                ${selectedRoom?.price_per_night} per night
              </Typography>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <Box sx={{ mt: 2 }}>
                  <DatePicker
                    label="Check-in Date"
                    value={checkIn}
                    onChange={(newValue) => setCheckIn(newValue)}
                    minDate={new Date()}
                  />
                </Box>
                <Box sx={{ mt: 2 }}>
                  <DatePicker
                    label="Check-out Date"
                    value={checkOut}
                    onChange={(newValue) => setCheckOut(newValue)}
                    minDate={checkIn || new Date()}
      />
    </Box>
              </LocalizationProvider>
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button
            onClick={handleBookRoom}
            variant="contained"
            disabled={!checkIn || !checkOut || bookingSuccess}
          >
            Confirm Booking
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Rooms; 
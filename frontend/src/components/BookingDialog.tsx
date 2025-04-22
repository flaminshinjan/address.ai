import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Box,
  TextField,
  Typography,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { roomApi } from '../services/api';

interface BookingDialogProps {
  open: boolean;
  onClose: () => void;
}

interface Room {
  id: string;
  name: string;
  price_per_night: number;
}

const BookingDialog: React.FC<BookingDialogProps> = ({ open, onClose }) => {
  const queryClient = useQueryClient();
  const [roomId, setRoomId] = useState('');
  const [checkInDate, setCheckInDate] = useState<Date | null>(null);
  const [checkOutDate, setCheckOutDate] = useState<Date | null>(null);
  const [specialRequests, setSpecialRequests] = useState('');

  const { data: rooms } = useQuery<Room[]>({
    queryKey: ['rooms'],
    queryFn: async () => {
      const response = await roomApi.get('/rooms?select=*');
      return response.data;
    },
  });

  const selectedRoom = rooms?.find(room => room.id === roomId);

  const calculateTotalPrice = () => {
    if (!checkInDate || !checkOutDate || !selectedRoom) return 0;
    const days = Math.ceil((checkOutDate.getTime() - checkInDate.getTime()) / (1000 * 60 * 60 * 24));
    return days * selectedRoom.price_per_night;
  };

  const createBooking = useMutation({
    mutationFn: async (data: {
      room_id: string;
      check_in_date: string;
      check_out_date: string;
      total_price: number;
      special_requests?: string;
    }) => {
      return roomApi.post('/bookings', data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['bookings'] });
      handleClose();
    },
  });

  const handleClose = () => {
    setRoomId('');
    setCheckInDate(null);
    setCheckOutDate(null);
    setSpecialRequests('');
    onClose();
  };

  const handleSubmit = () => {
    if (!roomId || !checkInDate || !checkOutDate) return;

    createBooking.mutate({
      room_id: roomId,
      check_in_date: checkInDate.toISOString().split('T')[0],
      check_out_date: checkOutDate.toISOString().split('T')[0],
      total_price: calculateTotalPrice(),
      special_requests: specialRequests || undefined,
    });
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>New Booking</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
          <FormControl fullWidth>
            <InputLabel>Room</InputLabel>
            <Select
              value={roomId}
              label="Room"
              onChange={(e) => setRoomId(e.target.value)}
            >
              {rooms?.map((room) => (
                <MenuItem key={room.id} value={room.id}>
                  {room.name} - ${room.price_per_night}/night
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <DatePicker
            label="Check-in Date"
            value={checkInDate}
            onChange={(newValue) => setCheckInDate(newValue)}
            disablePast
            slotProps={{
              textField: {
                fullWidth: true,
              },
            }}
          />

          <DatePicker
            label="Check-out Date"
            value={checkOutDate}
            onChange={(newValue) => setCheckOutDate(newValue)}
            disablePast
            minDate={checkInDate || undefined}
            slotProps={{
              textField: {
                fullWidth: true,
              },
            }}
          />

          <TextField
            label="Special Requests"
            multiline
            rows={3}
            fullWidth
            value={specialRequests}
            onChange={(e) => setSpecialRequests(e.target.value)}
          />

          {selectedRoom && checkInDate && checkOutDate && (
            <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
              <Typography variant="subtitle1">
                Total Price: ${calculateTotalPrice().toFixed(2)}
              </Typography>
            </Box>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={!roomId || !checkInDate || !checkOutDate}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default BookingDialog; 
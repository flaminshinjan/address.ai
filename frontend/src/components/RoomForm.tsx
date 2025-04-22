import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Grid,
  MenuItem,
} from '@mui/material';
import { Room } from '../types';

interface RoomFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: Omit<Room, 'id' | 'createdAt'>) => void;
  initialData?: Room;
}

export const RoomForm: React.FC<RoomFormProps> = ({
  open,
  onClose,
  onSubmit,
  initialData,
}) => {
  const [formData, setFormData] = React.useState({
    number: initialData?.number || '',
    type: initialData?.type || '',
    floor: initialData?.floor || 1,
    description: initialData?.description || '',
    capacity: initialData?.capacity || 2,
    pricePerDay: initialData?.pricePerDay || 0,
    status: initialData?.status || 'available',
  });

  React.useEffect(() => {
    if (initialData) {
      setFormData({
        number: initialData.number,
        type: initialData.type,
        floor: initialData.floor,
        description: initialData.description,
        capacity: initialData.capacity,
        pricePerDay: initialData.pricePerDay,
        status: initialData.status,
      });
    }
  }, [initialData]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'floor' || name === 'capacity' || name === 'pricePerDay'
        ? Number(value)
        : value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{initialData ? 'Edit Room' : 'Add Room'}</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                name="number"
                label="Room Number"
                value={formData.number}
                onChange={handleChange}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                name="type"
                label="Room Type"
                value={formData.type}
                onChange={handleChange}
                fullWidth
                required
                select
              >
                <MenuItem value="standard">Standard</MenuItem>
                <MenuItem value="deluxe">Deluxe</MenuItem>
                <MenuItem value="suite">Suite</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                name="floor"
                label="Floor"
                type="number"
                value={formData.floor}
                onChange={handleChange}
                fullWidth
                required
                inputProps={{ min: 1 }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                name="capacity"
                label="Capacity"
                type="number"
                value={formData.capacity}
                onChange={handleChange}
                fullWidth
                required
                inputProps={{ min: 1 }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                name="description"
                label="Description"
                value={formData.description}
                onChange={handleChange}
                fullWidth
                multiline
                rows={3}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                name="pricePerDay"
                label="Price per Day"
                type="number"
                value={formData.pricePerDay}
                onChange={handleChange}
                fullWidth
                required
                inputProps={{ min: 0, step: 0.01 }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                name="status"
                label="Status"
                value={formData.status}
                onChange={handleChange}
                fullWidth
                required
                select
              >
                <MenuItem value="available">Available</MenuItem>
                <MenuItem value="occupied">Occupied</MenuItem>
                <MenuItem value="maintenance">Maintenance</MenuItem>
              </TextField>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained" color="primary">
            {initialData ? 'Save' : 'Add'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}; 
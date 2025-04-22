import React, { useState } from 'react';
import { Box, Button, Container, Typography } from '@mui/material';
import { DataGrid, GridColDef, GridRenderCellParams } from '@mui/x-data-grid';
import { useQuery } from '@tanstack/react-query';
import { roomApi } from '../services/api';
import BookingDialog from '../components/BookingDialog';

interface Booking {
  id: string;
  room_id: string;
  room: {
    id: string;
    name: string;
  };
  user_id: string;
  check_in_date: string;
  check_out_date: string;
  total_price: number;
  status: string;
  special_requests?: string;
  created_at: string;
  updated_at: string;
}

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', width: 90 },
  {
    field: 'room',
    headerName: 'Room',
    width: 130,
    valueGetter: (params: GridRenderCellParams<Booking>) => params.row.room?.name || 'N/A',
  },
  {
    field: 'check_in_date',
    headerName: 'Check In',
    width: 130,
    valueGetter: (params: GridRenderCellParams<Booking>) =>
      new Date(params.row.check_in_date).toLocaleDateString(),
  },
  {
    field: 'check_out_date',
    headerName: 'Check Out',
    width: 130,
    valueGetter: (params: GridRenderCellParams<Booking>) =>
      new Date(params.row.check_out_date).toLocaleDateString(),
  },
  {
    field: 'total_price',
    headerName: 'Price',
    width: 100,
    valueFormatter: (params) => {
      if (params.value == null) return '';
      return `$${params.value.toFixed(2)}`;
    },
  },
  {
    field: 'status',
    headerName: 'Status',
    width: 130,
    renderCell: (params: GridRenderCellParams<Booking>) => (
      <Box
        sx={{
          backgroundColor: params.value === 'confirmed' ? 'success.main' : 
                          params.value === 'pending' ? 'warning.main' : 
                          params.value === 'cancelled' ? 'error.main' : 'info.main',
          color: '#fff',
          padding: '4px 8px',
          borderRadius: '4px',
          fontSize: '0.875rem',
        }}
      >
        {params.value.charAt(0).toUpperCase() + params.value.slice(1)}
      </Box>
    ),
  },
  {
    field: 'special_requests',
    headerName: 'Special Requests',
    width: 200,
    valueGetter: (params: GridRenderCellParams<Booking>) => 
      params.row.special_requests || 'None',
  },
  {
    field: 'created_at',
    headerName: 'Booked On',
    width: 180,
    valueGetter: (params: GridRenderCellParams<Booking>) =>
      new Date(params.row.created_at).toLocaleString(),
  },
];

const Bookings: React.FC = () => {
  const [dialogOpen, setDialogOpen] = useState(false);

  const { data: bookings, isLoading, error } = useQuery<Booking[]>({
    queryKey: ['bookings'],
    queryFn: async () => {
      try {
        const response = await roomApi.get('/bookings', {
          params: {
            select: '*,room(*)',
            order: 'created_at.desc',
          },
        });
        return response.data;
      } catch (err) {
        console.error('Error fetching bookings:', err);
        throw err;
      }
    },
  });

  if (error) {
    return (
      <Container maxWidth="lg">
        <Box sx={{ mt: 4, mb: 4 }}>
          <Typography color="error">
            Error loading bookings. Please try again later.
          </Typography>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            mb: 2,
          }}
        >
          <Typography variant="h4" component="h1">
            Bookings
          </Typography>
          <Button
            variant="contained"
            onClick={() => setDialogOpen(true)}
          >
            New Booking
          </Button>
        </Box>

        <DataGrid
          rows={bookings || []}
          columns={columns}
          loading={isLoading}
          initialState={{
            pagination: {
              paginationModel: { page: 0, pageSize: 10 },
            },
          }}
          pageSizeOptions={[10, 25, 50]}
          autoHeight
          getRowId={(row) => row.id}
          sx={{
            '& .MuiDataGrid-cell': {
              fontSize: '0.875rem',
            },
          }}
        />

        <BookingDialog
          open={dialogOpen}
          onClose={() => setDialogOpen(false)}
        />
      </Box>
    </Container>
  );
};

export default Bookings; 
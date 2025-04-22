import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Chip,
} from '@mui/material';
import { format } from 'date-fns';
import { foodService } from '../services/food';
import { Order, OrderItem } from '../types';

const Orders: React.FC = () => {
  const { data: orders, isLoading, error } = useQuery({
    queryKey: ['orders'],
    queryFn: foodService.getOrders,
  });

  if (isLoading) return <Box>Loading...</Box>;
  if (error) return <Box>Error loading orders</Box>;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Orders
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Items</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Total Price</TableCell>
              <TableCell>Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {orders?.data?.map((order: Order) => (
              <TableRow key={order.id}>
                <TableCell>{order.id}</TableCell>
                <TableCell>
                  {order.items.map((item: OrderItem) => (
                    <Box key={item.menuItemId}>
                      {item.quantity}x {item.menuItemId}
                    </Box>
                  ))}
                </TableCell>
                <TableCell>
                  <Chip
                    label={order.status}
                    color={
                      order.status === 'pending'
                        ? 'warning'
                        : order.status === 'preparing'
                        ? 'info'
                        : order.status === 'ready'
                        ? 'success'
                        : 'default'
                    }
                  />
                </TableCell>
                <TableCell>${order.totalPrice}</TableCell>
                <TableCell>
                  {format(new Date(order.createdAt), 'MMM dd, yyyy HH:mm')}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export default Orders; 
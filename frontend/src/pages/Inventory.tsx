import React, { useState } from 'react';
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
  CircularProgress,
  Alert,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import { supplyService } from '../services/supply';
import { InventoryItem, PurchaseOrderItem } from '../types';

const Inventory: React.FC = () => {
  const [selectedItem, setSelectedItem] = useState<InventoryItem | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [openDialog, setOpenDialog] = useState(false);
  const [unitPrice, setUnitPrice] = useState(0);

  const { data: inventoryItems, isLoading, error } = useQuery({
    queryKey: ['inventory'],
    queryFn: supplyService.getInventoryItems,
  });

  const handleOrder = (item: InventoryItem) => {
    setSelectedItem(item);
    setOpenDialog(true);
  };

  const handleConfirmOrder = async () => {
    if (!selectedItem) return;

    try {
      const orderItem: PurchaseOrderItem = {
        inventoryItemId: selectedItem.id,
        quantity,
        unitPrice,
      };
      await supplyService.createPurchaseOrder(selectedItem.supplierId, [orderItem]);
      setOpenDialog(false);
      setUnitPrice(0);
    } catch (error) {
      console.error('Failed to create purchase order:', error);
    }
  };

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        Error loading inventory
      </Alert>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Inventory
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Category</TableCell>
              <TableCell>Quantity</TableCell>
              <TableCell>Unit</TableCell>
              <TableCell>Minimum Quantity</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {inventoryItems?.data?.map((item: InventoryItem) => (
              <TableRow key={item.id}>
                <TableCell>{item.name}</TableCell>
                <TableCell>{item.category}</TableCell>
                <TableCell>{item.quantity}</TableCell>
                <TableCell>{item.unit}</TableCell>
                <TableCell>{item.minimumQuantity}</TableCell>
                <TableCell>
                  <Chip
                    label={item.quantity <= item.minimumQuantity ? 'Low Stock' : 'In Stock'}
                    color={item.quantity <= item.minimumQuantity ? 'warning' : 'success'}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={() => handleOrder(item)}
                  >
                    Order
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
        <DialogTitle>Create Purchase Order</DialogTitle>
        <DialogContent>
          <Typography gutterBottom>
            {selectedItem?.name}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Current Stock: {selectedItem?.quantity} {selectedItem?.unit}
          </Typography>
          <TextField
            label="Quantity"
            type="number"
            value={quantity}
            onChange={(e) => setQuantity(Number(e.target.value))}
            fullWidth
            margin="normal"
          />
          <TextField
            label="Unit Price"
            type="number"
            value={unitPrice}
            onChange={(e) => setUnitPrice(Number(e.target.value))}
            fullWidth
            margin="normal"
            inputProps={{ min: 0, step: 0.01 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button onClick={handleConfirmOrder} color="primary">
            Confirm
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Inventory; 
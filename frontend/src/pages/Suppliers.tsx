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
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import { supplyService } from '../services/supply';
import { Supplier } from '../types';

const Suppliers: React.FC = () => {
  const { data: suppliers, isLoading, error } = useQuery({
    queryKey: ['suppliers'],
    queryFn: supplyService.getSuppliers,
  });

  const [selectedSupplier, setSelectedSupplier] = useState<Supplier | null>(null);
  const [openDialog, setOpenDialog] = useState(false);

  const handleViewDetails = (supplier: Supplier) => {
    setSelectedSupplier(supplier);
    setOpenDialog(true);
  };

  if (isLoading) return <Box>Loading...</Box>;
  if (error) return <Box>Error loading suppliers</Box>;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Suppliers
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Contact Person</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Phone</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {suppliers?.data?.map((supplier: Supplier) => (
              <TableRow key={supplier.id}>
                <TableCell>{supplier.name}</TableCell>
                <TableCell>{supplier.contactPerson}</TableCell>
                <TableCell>{supplier.email}</TableCell>
                <TableCell>{supplier.phone}</TableCell>
                <TableCell>
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => handleViewDetails(supplier)}
                  >
                    View Details
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
        <DialogTitle>Supplier Details</DialogTitle>
        <DialogContent>
          <Typography variant="h6" gutterBottom>
            {selectedSupplier?.name}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Contact Person: {selectedSupplier?.contactPerson}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Email: {selectedSupplier?.email}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Phone: {selectedSupplier?.phone}
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Address: {selectedSupplier?.address}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Suppliers; 
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
  TextField,
  Box,
  IconButton,
  Alert,
  Snackbar,
  Chip,
  Divider,
} from '@mui/material';
import { Add as AddIcon, Remove as RemoveIcon } from '@mui/icons-material';
import { foodApi } from '../services/api';

interface MenuItem {
  id: string;
  name: string;
  description: string;
  category: string;
  price: number;
  preparation_time: number;
  is_available: boolean;
}

interface CartItem extends MenuItem {
  quantity: number;
}

const Menu = () => {
  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [openCart, setOpenCart] = useState(false);
  const [roomNumber, setRoomNumber] = useState('');
  const [specialInstructions, setSpecialInstructions] = useState('');
  const [orderSuccess, setOrderSuccess] = useState(false);

  useEffect(() => {
    fetchMenuItems();
  }, []);

  const fetchMenuItems = async () => {
    try {
      const response = await foodApi.get('/menu-items');
      setMenuItems(response.data);
      setLoading(false);
    } catch (err) {
      setError('Failed to load menu items');
      setLoading(false);
    }
  };

  const addToCart = (item: MenuItem) => {
    setCart((prevCart) => {
      const existingItem = prevCart.find((cartItem) => cartItem.id === item.id);
      if (existingItem) {
        return prevCart.map((cartItem) =>
          cartItem.id === item.id
            ? { ...cartItem, quantity: cartItem.quantity + 1 }
            : cartItem
        );
      }
      return [...prevCart, { ...item, quantity: 1 }];
    });
  };

  const removeFromCart = (itemId: string) => {
    setCart((prevCart) => {
      const existingItem = prevCart.find((item) => item.id === itemId);
      if (existingItem && existingItem.quantity > 1) {
        return prevCart.map((item) =>
          item.id === itemId ? { ...item, quantity: item.quantity - 1 } : item
        );
      }
      return prevCart.filter((item) => item.id !== itemId);
    });
  };

  const getTotalPrice = () => {
    return cart.reduce((total, item) => total + item.price * item.quantity, 0);
  };

  const handlePlaceOrder = async () => {
    if (!roomNumber) {
      setError('Please enter your room number');
      return;
    }

    try {
      await foodApi.post('/orders', {
        room_number: roomNumber,
        items: cart.map((item) => ({
          menu_item_id: item.id,
          quantity: item.quantity,
          unit_price: item.price,
        })),
        total_amount: getTotalPrice(),
        special_instructions: specialInstructions,
      });

      setOrderSuccess(true);
      setCart([]);
      setRoomNumber('');
      setSpecialInstructions('');
      setTimeout(() => {
        setOpenCart(false);
        setOrderSuccess(false);
      }, 2000);
    } catch (err) {
      setError('Failed to place order');
    }
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Alert severity="error">{error}</Alert>;

  const categories = Array.from(new Set(menuItems.map((item) => item.category)));

  return (
    <Container>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4" gutterBottom sx={{ my: 4 }}>
          Room Service Menu
        </Typography>
        <Button
          variant="contained"
          onClick={() => setOpenCart(true)}
          disabled={cart.length === 0}
        >
          View Cart ({cart.length} items)
        </Button>
      </Box>

      {categories.map((category) => (
        <Box key={category} sx={{ mb: 6 }}>
          <Typography variant="h5" gutterBottom color="primary">
            {category}
          </Typography>
          <Divider sx={{ mb: 2 }} />
          <Grid container spacing={4}>
            {menuItems
              .filter((item) => item.category === category)
              .map((item) => (
                <Grid item xs={12} sm={6} md={4} key={item.id}>
                  <Card>
                    <CardMedia
                      component="img"
                      height="200"
                      image={`/food-${item.name.toLowerCase().replace(/\s+/g, '-')}.jpg`}
                      alt={item.name}
                    />
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        {item.name}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" paragraph>
                        {item.description}
                      </Typography>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Typography variant="h6" color="primary">
                          ${item.price.toFixed(2)}
                        </Typography>
                        <Chip
                          label={`${item.preparation_time} mins`}
                          size="small"
                          color="secondary"
                        />
                      </Box>
                      <Button
                        variant="contained"
                        fullWidth
                        sx={{ mt: 2 }}
                        disabled={!item.is_available}
                        onClick={() => addToCart(item)}
                      >
                        {item.is_available ? 'Add to Cart' : 'Not Available'}
                      </Button>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
          </Grid>
        </Box>
      ))}

      <Dialog open={openCart} onClose={() => setOpenCart(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Your Cart</DialogTitle>
        <DialogContent>
          {orderSuccess ? (
            <Alert severity="success">Order placed successfully!</Alert>
          ) : (
            <>
              {cart.map((item) => (
                <Box
                  key={item.id}
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    mb: 2,
                  }}
                >
                  <Box>
                    <Typography variant="body1">{item.name}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      ${item.price.toFixed(2)} x {item.quantity}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <IconButton size="small" onClick={() => removeFromCart(item.id)}>
                      <RemoveIcon />
                    </IconButton>
                    <Typography sx={{ mx: 1 }}>{item.quantity}</Typography>
                    <IconButton size="small" onClick={() => addToCart(item)}>
                      <AddIcon />
                    </IconButton>
                  </Box>
                </Box>
              ))}
              <Divider sx={{ my: 2 }} />
              <Typography variant="h6" gutterBottom>
                Total: ${getTotalPrice().toFixed(2)}
              </Typography>
              <TextField
                fullWidth
                label="Room Number"
                value={roomNumber}
                onChange={(e) => setRoomNumber(e.target.value)}
                margin="normal"
                required
              />
              <TextField
                fullWidth
                label="Special Instructions"
                value={specialInstructions}
                onChange={(e) => setSpecialInstructions(e.target.value)}
                margin="normal"
                multiline
                rows={3}
              />
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCart(false)}>Cancel</Button>
          <Button
            onClick={handlePlaceOrder}
            variant="contained"
            disabled={cart.length === 0 || !roomNumber || orderSuccess}
          >
            Place Order
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={!!error}
        autoHideDuration={6000}
        onClose={() => setError(null)}
        message={error}
      />
    </Container>
  );
};

export default Menu; 
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  useTheme,
  alpha,
  styled,
  Avatar,
  IconButton,
  Tab,
  Tabs,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  LinearProgress,
} from '@mui/material';
import {
  Room as RoomIcon,
  Menu as MenuIcon,
  Inventory as InventoryIcon,
  LocalShipping as SupplyIcon,
  Warning as WarningIcon,
  MoreVert as MoreVertIcon,
  KeyboardArrowUp as KeyboardArrowUpIcon,
  KeyboardArrowDown as KeyboardArrowDownIcon,
} from '@mui/icons-material';
import { roomService } from '../services/room';
import { foodService } from '../services/food';
import { supplyService } from '../services/supply';
import { LoadingSpinner, ErrorMessage } from '../components';
import { Room, Booking, MenuItem, Order, InventoryItem } from '../types';
import { BookingTrendChart, PlatformBookingChart, WeeklyVisitorsChart } from '../components/charts';

const StyledCard = styled(Card)(({ theme }) => ({
  height: '100%',
  background: theme.palette.background.paper,
  borderRadius: theme.shape.borderRadius * 2,
  border: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
}));

const StatBox = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
  borderRadius: theme.shape.borderRadius,
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(2),
}));

const Dashboard: React.FC = () => {
  const theme = useTheme();
  const [timeRange, setTimeRange] = React.useState('Month');

  const { data: rooms, isLoading: roomsLoading, error: roomsError } = useQuery<Room[]>({
    queryKey: ['rooms'],
    queryFn: () => roomService.getRooms(),
  });

  const { data: bookings, isLoading: bookingsLoading } = useQuery<Booking[]>({
    queryKey: ['bookings'],
    queryFn: () => roomService.getBookings(),
  });

  const { data: menuItems, isLoading: menuLoading } = useQuery<MenuItem[]>({
    queryKey: ['menuItems'],
    queryFn: async () => {
      const response = await foodService.getMenuItems();
      return response.data;
    },
  });

  const { data: orders, isLoading: ordersLoading } = useQuery<Order[]>({
    queryKey: ['orders'],
    queryFn: async () => {
      const response = await foodService.getOrders();
      return response.data;
    },
  });

  const { data: inventory, isLoading: inventoryLoading } = useQuery<InventoryItem[]>({
    queryKey: ['inventory'],
    queryFn: async () => {
      const response = await supplyService.getInventoryItems();
      return response.data;
    },
  });

  const { data: lowStock, isLoading: lowStockLoading } = useQuery<InventoryItem[]>({
    queryKey: ['lowStock'],
    queryFn: async () => {
      const response = await supplyService.getLowStockItems();
      return response.data;
    },
  });

  const isLoading = roomsLoading || bookingsLoading || menuLoading || ordersLoading || inventoryLoading || lowStockLoading;

  if (isLoading) return <LoadingSpinner />;
  if (roomsError) return <ErrorMessage message="Failed to load dashboard data" />;

  return (
    <Box>
      <Grid container spacing={3}>
        {/* Available Rooms Card */}
        <Grid item xs={12} md={6}>
          <StyledCard>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Box>
                  <Typography variant="h6" sx={{ mb: 0.5 }}>Available room</Typography>
                  <Typography variant="h3" sx={{ fontWeight: 700 }}>
                    126
                    <Typography component="span" color="success.main" sx={{ fontSize: '1rem', ml: 1 }}>
                      +30%
                    </Typography>
                  </Typography>
                </Box>
                <RoomIcon sx={{ fontSize: 40, color: 'primary.main' }} />
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Tabs value={timeRange} onChange={(e, v) => setTimeRange(v)}>
                  <Tab label="Month" value="Month" />
                  <Tab label="Year" value="Year" />
                </Tabs>
              </Box>

              <Box sx={{ height: 200 }}>
                <BookingTrendChart />
              </Box>
            </CardContent>
          </StyledCard>
        </Grid>

        {/* Booking by Platform */}
        <Grid item xs={12} md={6}>
          <StyledCard>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Typography variant="h6">Booking by Platform</Typography>
                <IconButton size="small">
                  <MoreVertIcon />
                </IconButton>
              </Box>
              <Box sx={{ height: 300 }}>
                <PlatformBookingChart />
              </Box>
            </CardContent>
          </StyledCard>
        </Grid>

        {/* Staff Section */}
        <Grid item xs={12}>
          <StyledCard>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Typography variant="h6">Staff</Typography>
                <Box sx={{ display: 'flex', gap: 2 }}>
                  <Typography variant="body2" color="text.secondary">
                    Staff productivity overview
                  </Typography>
                  <IconButton size="small">
                    <MoreVertIcon />
                  </IconButton>
                </Box>
              </Box>

              <Grid container spacing={3} sx={{ mb: 4 }}>
                <Grid item xs={12} sm={3}>
                  <StatBox>
                    <Avatar>32</Avatar>
                    <Box>
                      <Typography variant="body2" color="text.secondary">Total</Typography>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>employee</Typography>
                    </Box>
                  </StatBox>
                </Grid>
                <Grid item xs={12} sm={3}>
                  <StatBox>
                    <Avatar sx={{ bgcolor: 'success.light' }}>42</Avatar>
                    <Box>
                      <Typography variant="body2" color="text.secondary">Cleaned</Typography>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>rooms</Typography>
                    </Box>
                  </StatBox>
                </Grid>
                <Grid item xs={12} sm={3}>
                  <StatBox>
                    <Avatar sx={{ bgcolor: 'warning.light' }}>24</Avatar>
                    <Box>
                      <Typography variant="body2" color="text.secondary">Pending</Typography>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>rooms</Typography>
                    </Box>
                  </StatBox>
                </Grid>
                <Grid item xs={12} sm={3}>
                  <StatBox>
                    <Box sx={{ width: '100%' }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                        <Typography variant="body2" color="text.secondary">Progress</Typography>
                        <Typography variant="body2" color="success.main">62%</Typography>
                      </Box>
                      <LinearProgress variant="determinate" value={62} sx={{ height: 8, borderRadius: 4 }} />
                    </Box>
                  </StatBox>
                </Grid>
              </Grid>

              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Role</TableCell>
                      <TableCell>Total room</TableCell>
                      <TableCell>Start work</TableCell>
                      <TableCell>End work</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {[
                      { name: 'Guy Hawkins', role: 'Chef', rooms: 28, start: '8:00 AM', end: '8:00 PM' },
                      { name: 'Eleanor Pena', role: 'Receptionist', rooms: 16, start: '8:00 AM', end: '4:00 PM' },
                      { name: 'Robert Fox', role: 'Cleaner', rooms: 20, start: '8:00 AM', end: '6:00 PM' },
                    ].map((staff) => (
                      <TableRow key={staff.name}>
                        <TableCell>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                            <Avatar>{staff.name[0]}</Avatar>
                            <Typography>{staff.name}</Typography>
                          </Box>
                        </TableCell>
                        <TableCell>{staff.role}</TableCell>
                        <TableCell>{staff.rooms}</TableCell>
                        <TableCell>{staff.start}</TableCell>
                        <TableCell>{staff.end}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </StyledCard>
        </Grid>

        {/* Weekly Visitors */}
        <Grid item xs={12} md={6}>
          <StyledCard>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Box>
                  <Typography variant="h6">Weekly Visitors</Typography>
                  <Typography variant="h4" sx={{ fontWeight: 700, my: 2 }}>
                    120
                    <Typography component="span" color="primary" sx={{ fontSize: '1rem', ml: 1 }}>
                      visitors
                    </Typography>
                  </Typography>
                  <Typography variant="body2" color="success.main">
                    keep it up! 👍
                  </Typography>
                </Box>
                <Box sx={{ width: '50%' }}>
                  <WeeklyVisitorsChart />
                </Box>
              </Box>
            </CardContent>
          </StyledCard>
        </Grid>

        {/* Booking List */}
        <Grid item xs={12} md={6}>
          <StyledCard>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Typography variant="h6">Booking List</Typography>
                <Typography variant="body2" color="text.secondary">
                  All bookings at a glance
                </Typography>
              </Box>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Room No.</TableCell>
                      <TableCell>Room</TableCell>
                      <TableCell>Check in</TableCell>
                      <TableCell>Check out</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {[
                      { name: 'Floyd Miles', room: '113', type: 'King', checkIn: '5 March', checkOut: '8 March' },
                      { name: 'Devon Lane', room: '101', type: 'Deluxe', checkIn: '9 March', checkOut: '10 March' },
                    ].map((booking) => (
                      <TableRow key={booking.name}>
                        <TableCell>{booking.name}</TableCell>
                        <TableCell>{booking.room}</TableCell>
                        <TableCell>{booking.type}</TableCell>
                        <TableCell>{booking.checkIn}</TableCell>
                        <TableCell>{booking.checkOut}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </StyledCard>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard; 
import axios from 'axios';
import { MenuItem, Order } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const foodService = {
  getMenuItems: async () => {
    const response = await axios.get<MenuItem[]>(`${API_URL}/api/menu`);
    return response;
  },

  getOrders: async () => {
    const response = await axios.get<Order[]>(`${API_URL}/api/orders`);
    return response;
  },

  createOrder: async (items: { menuItemId: string; quantity: number }[]) => {
    const response = await axios.post<Order>(`${API_URL}/api/orders`, { items });
    return response;
  },

  updateOrderStatus: async (orderId: string, status: string) => {
    const response = await axios.patch<Order>(`${API_URL}/api/orders/${orderId}`, { status });
    return response;
  },
}; 
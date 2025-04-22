import axios from 'axios';
import { Room, Booking } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const roomService = {
  getRooms: async (): Promise<Room[]> => {
    const response = await axios.get<Room[]>(`${API_URL}/api/rooms`);
    return response.data;
  },

  getAvailableRooms: async (startDate: Date, endDate: Date): Promise<Room[]> => {
    const response = await axios.get<Room[]>(`${API_URL}/api/rooms/available`, {
      params: {
        startDate: startDate.toISOString(),
        endDate: endDate.toISOString(),
      },
    });
    return response.data;
  },

  getBookings: async (): Promise<Booking[]> => {
    const response = await axios.get<Booking[]>(`${API_URL}/api/bookings`);
    return response.data;
  },

  createBooking: async (roomId: string, startDate: Date, endDate: Date): Promise<Booking> => {
    const response = await axios.post<Booking>(`${API_URL}/api/bookings`, {
      roomId,
      startDate: startDate.toISOString(),
      endDate: endDate.toISOString(),
    });
    return response.data;
  },

  cancelBooking: async (bookingId: string): Promise<void> => {
    await axios.delete(`${API_URL}/api/bookings/${bookingId}`);
  },

  createRoom: async (room: Omit<Room, 'id' | 'createdAt'>): Promise<Room> => {
    const response = await axios.post<Room>(`${API_URL}/api/rooms`, room);
    return response.data;
  },

  updateRoom: async (id: string, room: Partial<Room>): Promise<Room> => {
    const response = await axios.put<Room>(`${API_URL}/api/rooms/${id}`, room);
    return response.data;
  },

  deleteRoom: async (id: string): Promise<void> => {
    await axios.delete(`${API_URL}/api/rooms/${id}`);
  },
}; 
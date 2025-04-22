// Room types
export interface Room {
  id: string;
  number: string;
  type: string;
  floor: number;
  description: string;
  capacity: number;
  pricePerDay: number;
  status: 'available' | 'occupied' | 'maintenance';
}

export interface Booking {
  id: string;
  roomId: string;
  userId: string;
  startDate: string;
  endDate: string;
  status: 'pending' | 'confirmed' | 'cancelled' | 'completed';
  totalPrice: number;
  createdAt: string;
}

// Food types
export interface MenuItem {
  id: string;
  name: string;
  description: string;
  price: number;
  category: string;
  imageUrl?: string;
  isAvailable: boolean;
}

export interface Order {
  id: string;
  userId: string;
  items: OrderItem[];
  status: 'pending' | 'preparing' | 'ready' | 'delivered' | 'cancelled';
  totalPrice: number;
  createdAt: string;
}

export interface OrderItem {
  menuItemId: string;
  quantity: number;
  price: number;
}

// Supply types
export interface Supplier {
  id: string;
  name: string;
  contactPerson: string;
  email: string;
  phone: string;
  address: string;
}

export interface InventoryItem {
  id: string;
  name: string;
  description: string;
  category: string;
  quantity: number;
  unit: string;
  minimumQuantity: number;
  supplierId: string;
}

export interface PurchaseOrder {
  id: string;
  supplierId: string;
  items: PurchaseOrderItem[];
  status: 'pending' | 'approved' | 'ordered' | 'received' | 'cancelled';
  totalAmount: number;
  createdAt: string;
}

export interface PurchaseOrderItem {
  inventoryItemId: string;
  quantity: number;
  unitPrice: number;
}

// Auth types
export interface User {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  data: User;
  token: string;
} 
import axios from 'axios';
import { InventoryItem, Supplier, PurchaseOrder } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const supplyService = {
  getInventoryItems: async () => {
    const response = await axios.get<InventoryItem[]>(`${API_URL}/api/inventory`);
    return response;
  },

  getLowStockItems: async () => {
    const response = await axios.get<InventoryItem[]>(`${API_URL}/api/inventory/low-stock`);
    return response;
  },

  getSuppliers: async () => {
    const response = await axios.get<Supplier[]>(`${API_URL}/api/suppliers`);
    return response;
  },

  getPurchaseOrders: async () => {
    const response = await axios.get<PurchaseOrder[]>(`${API_URL}/api/purchase-orders`);
    return response;
  },

  createPurchaseOrder: async (supplierId: string, items: { inventoryItemId: string; quantity: number }[]) => {
    const response = await axios.post<PurchaseOrder>(`${API_URL}/api/purchase-orders`, {
      supplierId,
      items,
    });
    return response;
  },

  updatePurchaseOrderStatus: async (orderId: string, status: string) => {
    const response = await axios.patch<PurchaseOrder>(`${API_URL}/api/purchase-orders/${orderId}`, { status });
    return response;
  },
}; 
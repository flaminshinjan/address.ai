import axios from 'axios';

const SUPABASE_URL = 'https://iurddbrzlrhqcdwjauvw.supabase.co';
const SUPABASE_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Iml1cmRkYnJ6bHJocWNkd2phdXZ3Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDUyODgzMzIsImV4cCI6MjA2MDg2NDMzMn0.lvzzluq_W9Nj4YLJMKCwq_8Ei0YAhCv9Q4gjhjlBReI';

const createApiInstance = (baseURL: string) => {
  const instance = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
      'apikey': SUPABASE_KEY,
      'Authorization': `Bearer ${SUPABASE_KEY}`,
      'Prefer': 'return=representation',
    },
  });

  // Add request interceptor for authentication
  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Add response interceptor for error handling
  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        // Handle unauthorized access
        localStorage.removeItem('token');
        localStorage.removeItem('userId');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }
  );

  return instance;
};

export const userApi = createApiInstance(`${SUPABASE_URL}/rest/v1`);
export const roomApi = createApiInstance(`${SUPABASE_URL}/rest/v1`);
export const foodApi = createApiInstance(`${SUPABASE_URL}/rest/v1`);
export const supplyApi = createApiInstance(`${SUPABASE_URL}/rest/v1`); 
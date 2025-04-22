import { userApi } from './api';

const SUPABASE_URL = 'https://iurddbrzlrhqcdwjauvw.supabase.co';
const SUPABASE_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Iml1cmRkYnJ6bHJocWNkd2phdXZ3Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDUyODgzMzIsImV4cCI6MjA2MDg2NDMzMn0.lvzzluq_W9Nj4YLJMKCwq_8Ei0YAhCv9Q4gjhjlBReI';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface AuthResponse {
  user: {
    id: string;
    email: string;
    user_metadata: {
      firstName: string;
      lastName: string;
    };
  };
  session: {
    access_token: string;
  };
}

export const authService = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await fetch(`${SUPABASE_URL}/auth/v1/token?grant_type=password`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'apikey': SUPABASE_KEY,
      },
      body: JSON.stringify({
        email: credentials.email,
        password: credentials.password,
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error_description || 'Failed to login');
    }

    const data = await response.json();
    localStorage.setItem('token', data.access_token);
    localStorage.setItem('userId', data.user.id);
    localStorage.setItem('user', JSON.stringify(data.user));
    return data;
  },

  async register(data: RegisterData): Promise<AuthResponse> {
    const response = await fetch(`${SUPABASE_URL}/auth/v1/signup`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'apikey': SUPABASE_KEY,
      },
      body: JSON.stringify({
        email: data.email,
        password: data.password,
        data: {
          firstName: data.firstName,
          lastName: data.lastName,
        },
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error_description || 'Failed to register');
    }

    const responseData = await response.json();
    return this.login({ email: data.email, password: data.password });
  },

  logout(): void {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    localStorage.removeItem('user');
  },

  getCurrentUser(): any {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  },
}; 
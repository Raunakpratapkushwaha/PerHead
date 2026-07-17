import { create } from 'zustand';
import { apiCall } from '../services/api';

interface User {
  id: number;
  name: string;
  email: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isLoading: false,
  error: null,

  login: async (email, password) => {
    set({ isLoading: true, error: null });
    try {
      // 🚀 Hits your Go backend running on your Mac!
      const data = await apiCall('/auth/login', 'POST', { email, password });
      
      // If successful, save the token and user data to global state
      set({ 
        token: data.access_token, 
        user: data.user, 
        isLoading: false 
      });
      
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  logout: () => {
    set({ user: null, token: null, error: null });
  }
}));
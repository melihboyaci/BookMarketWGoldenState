import axios from 'axios';

// Ensure the trailing slash is omitted, or handle it consistently.
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to attach JWT token
api.interceptors.request.use(
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

// Response interceptor to handle 401s (optional, but good practice)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token might be expired or invalid
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      // window.location.href = '/login'; // Optional: Redirect to login
    }
    return Promise.reject(error);
  }
);

export const authService = {
  login: async (email, password) => {
    const response = await api.post('/auth/login', { email, password });
    return response.data; // { token: '...', role: '...' }
  },
};

export const bookService = {
  getBooks: async () => {
    // Backend doesn't have a GET /books endpoint yet according to main.go, 
    // but Phase 2 plan requests this function.
    try {
      const response = await api.get('/books');
      return response.data;
    } catch (error) {
      // Mocked data if endpoint is missing to allow frontend testing
      console.warn("API /books failed, returning mock data.", error);
      return [];
    }
  },
};

export const systemService = {
  restoreGoldenState: async () => {
    const response = await api.post('/system/restore');
    return response.data;
  },
};

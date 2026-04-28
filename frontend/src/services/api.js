import axios from 'axios';

// Ensure the trailing slash is omitted, or handle it consistently.
const API_URL = import.meta.env.VITE_API_URL || '/api/v1';

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
  register: async (username, email, password) => {
    const response = await api.post('/auth/register', { username, email, password });
    return response.data; // { token: '...', role: '...' }
  },
};

export const bookService = {
  getBooks: async () => {
    const response = await api.get('/books');
    return response.data;
  },
  createBook: async (bookData) => {
    const response = await api.post('/books', bookData);
    return response.data;
  },
  updateBook: async (id, bookData) => {
    const response = await api.put(`/books/${id}`, bookData);
    return response.data;
  },
  deleteBook: async (id) => {
    const response = await api.delete(`/books/${id}`);
    return response.data;
  }
};

export const cartService = {
  checkout: async (items) => {
    const response = await api.post('/checkout', { items });
    return response.data;
  }
};

export const systemService = {
  getSales: async () => {
    const response = await api.get('/sales');
    return response.data;
  },
  restoreGoldenState: async () => {
    const response = await api.post('/admin/system/restore');
    return response.data;
  },
};

export const adminService = {
  getUsers: async () => {
    const response = await api.get('/admin/users');
    return response.data;
  },
};

import React, { createContext, useState, useEffect } from 'react';
import { authService } from '../services/api';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    const savedUser = localStorage.getItem('user');
    return savedUser ? JSON.parse(savedUser) : null;
  });
  const [token, setToken] = useState(() => localStorage.getItem('token') || null);
  const [loading, setLoading] = useState(false);

  const _persist = (token, userData) => {
    setToken(token);
    setUser(userData);
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(userData));
  };

  const login = async (email, password) => {
    setLoading(true);
    try {
      const data = await authService.login(email, password);
      _persist(data.token, { email, role: data.role });
      return { success: true, role: data.role };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Giriş yapılamadı. Lütfen bilgilerinizi kontrol edin.',
      };
    } finally {
      setLoading(false);
    }
  };

  const register = async (username, email, password) => {
    setLoading(true);
    try {
      const data = await authService.register(username, email, password);
      _persist(data.token, { email, role: data.role });
      return { success: true, role: data.role };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Kayıt oluşturulamadı. Lütfen tekrar deneyin.',
      };
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

import React, { createContext, useContext, useState, useEffect } from 'react';
import authApi from '../services/api';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    if (token && userData) {
      try {
        const parsedUser = JSON.parse(userData);
        setUser(parsedUser);
        authApi.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      } catch (error) {
        console.error('Error parsing user data:', error);
        logout();
      }
    }
    setLoading(false);
  }, []);

  const login = async (username, password) => {
    try {
      const response = await authApi.post('/login', { username, password });
      console.log('=== LOGIN DEBUG ===');
      console.log('Full response object:', response);
      console.log('Response data:', response.data);
      console.log('Response data type:', typeof response.data);
      console.log('Response data keys:', Object.keys(response.data || {}));
      
      const { accessToken, user: userData } = response.data;
      
      console.log('Extracted accessToken:', accessToken);
      console.log('Extracted user:', userData);
      console.log('User type:', typeof userData);
      console.log('==================');
      
      if (!accessToken) {
        console.error('No access token in response!');
        return { success: false, error: 'No access token received' };
      }
      
      if (!userData) {
        console.error('No user data in response!');
        return { success: false, error: 'No user data received' };
      }
      
      localStorage.setItem('token', accessToken);
      localStorage.setItem('user', JSON.stringify(userData));
      authApi.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
      setUser(userData);
      
      return { success: true, user: userData };
    } catch (error) {
      console.error('Login error:', error);
      return { 
        success: false, 
        error: error.response?.data?.error || 'Login failed' 
      };
    }
  };

  const register = async (userData) => {
    try {
      const response = await authApi.post('/register', userData);
      const { accessToken, user: newUser } = response.data;
      
      localStorage.setItem('token', accessToken);
      localStorage.setItem('user', JSON.stringify(newUser));
      authApi.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
      setUser(newUser);
      
      return { success: true, user: newUser, data: response.data };
    } catch (error) {
      console.error('Registration error:', error);
      return { 
        success: false, 
        error: error.response?.data?.error || 'Registration failed' 
      };
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    delete authApi.defaults.headers.common['Authorization'];
    setUser(null);
  };

  const value = {
    user,
    login,
    register,
    logout,
    loading
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
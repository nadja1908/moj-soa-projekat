import axios from 'axios';

// API Gateway base URL
const API_GATEWAY_URL = 'http://localhost:8000';

// Create axios instance for API Gateway
const apiGateway = axios.create({
  baseURL: `${API_GATEWAY_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Auth API through Gateway
const authApi = axios.create({
  baseURL: `${API_GATEWAY_URL}/api/auth`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Users API through Gateway  
const usersApi = axios.create({
  baseURL: `${API_GATEWAY_URL}/api/users`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Blog API through Gateway
const blogApi = axios.create({
  baseURL: `${API_GATEWAY_URL}/api/blog`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Admin API through Gateway
const adminApi = axios.create({
  baseURL: `${API_GATEWAY_URL}/api/admin`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
const requestInterceptor = (config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

// Error interceptor for logout on 401
const handleAuthError = (error) => {
  if (error.response?.status === 401) {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
  }
  return Promise.reject(error);
};

// Add request interceptors to protected APIs
usersApi.interceptors.request.use(requestInterceptor);
blogApi.interceptors.request.use(requestInterceptor);
adminApi.interceptors.request.use(requestInterceptor);

// Add error interceptors to all APIs
authApi.interceptors.response.use((response) => response, handleAuthError);
usersApi.interceptors.response.use((response) => response, handleAuthError);
//blogApi.interceptors.response.use((response) => response, handleAuthError);
adminApi.interceptors.response.use((response) => response, handleAuthError);
// umesto ovoga dodaj poseban interceptor SAMO za logovanje greške blogApi-ja:
blogApi.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('blogApi error:', {
      url: error.config?.url,
      method: error.config?.method,
      status: error.response?.status,
      data: error.response?.data,
      headers: error.response?.headers,
    });
    return Promise.reject(error);
  })

// Backward compatibility exports
const stakeholdersApi = usersApi;

export default authApi;
export { authApi, usersApi, blogApi, adminApi, stakeholdersApi };
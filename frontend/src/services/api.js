import axios from 'axios';

// API base URL - use environment variable or default to localhost
export const API_BASE_URL = import.meta.env.VITE_API_URL || 
  (import.meta.env.DEV 
    ? 'http://localhost:3000' 
    : 'https://your-backend.railway.app');

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token
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

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      // Handle 401 Unauthorized - token expired or invalid
      if (error.response.status === 401) {
        // Only redirect if we're not already on the auth page
        // This prevents page reload during login attempts
        const currentPath = window.location.pathname;
        if (currentPath !== '/auth' && currentPath !== '/') {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          window.location.href = '/auth';
        }
      }
      
      // Extract error message from response
      const errorMessage = error.response.data?.error?.message || 
                          error.response.data?.message || 
                          'An error occurred';
      
      // Preserve status/code on the rewritten error so callers can branch on
      // them (the original axios error fields are otherwise lost here).
      const rewritten = new Error(errorMessage);
      rewritten.status = error.response.status;
      return Promise.reject(rewritten);
    } else if (error.request) {
      // Timeouts and unreachable servers both land here; keep the axios
      // error code (e.g. 'ECONNABORTED') for the same reason.
      const rewritten = new Error('No response from server');
      rewritten.code = error.code;
      return Promise.reject(rewritten);
    } else {
      return Promise.reject(error);
    }
  }
);

export default api;

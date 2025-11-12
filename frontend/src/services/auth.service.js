import api from './api';
import { getUserFromToken, isTokenExpired } from '../utils/jwt';

const authService = {
  /**
   * Register a new user
   * @param {string} username 
   * @param {string} email 
   * @param {string} password 
   * @returns {Promise<{user: Object, token: string}>}
   */
  async register(username, email, password) {
    const response = await api.post('/api/auth/register', {
      username,
      email,
      password,
    });
    
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      // Store user data from API response
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    
    return response.data;
  },

  /**
   * Login with username and password
   * @param {string} username 
   * @param {string} password 
   * @returns {Promise<{user: Object, token: string}>}
   */
  async login(username, password) {
    const response = await api.post('/api/auth/login', {
      username,
      password,
    });
    
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      // Store user data from API response
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    
    return response.data;
  },

  /**
   * Logout current user
   */
  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  /**
   * Get current user info
   * @returns {Promise<Object>}
   */
  async getCurrentUser() {
    const response = await api.get('/api/auth/me');
    return response.data;
  },

  /**
   * Get stored user from localStorage
   * Falls back to parsing from JWT token if user data is missing
   * @returns {Object|null}
   */
  getStoredUser() {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      return JSON.parse(userStr);
    }
    
    // Fallback: try to get user info from token
    const token = this.getStoredToken();
    if (token && !isTokenExpired(token)) {
      return getUserFromToken(token);
    }
    
    return null;
  },

  /**
   * Get stored token from localStorage
   * @returns {string|null}
   */
  getStoredToken() {
    return localStorage.getItem('token');
  },

  /**
   * Check if user is authenticated
   * @returns {boolean}
   */
  isAuthenticated() {
    const token = this.getStoredToken();
    return !!token && !isTokenExpired(token);
  },
};

export default authService;

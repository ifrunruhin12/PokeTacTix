import api from './api';

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
   * @returns {Object|null}
   */
  getStoredUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
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
    return !!this.getStoredToken();
  },
};

export default authService;

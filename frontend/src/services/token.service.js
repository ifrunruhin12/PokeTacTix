import api from './api';

/**
 * Token Service
 * Handles all token-related API calls
 */
const tokenService = {
  /**
   * Get current user's token balance and reset information
   * @returns {Promise<Object>} Token balance data
   */
  async getTokenBalance() {
    try {
      const response = await api.get('/api/tokens/balance');
      return response.data;
    } catch (error) {
      console.error('Error fetching token balance:', error);
      throw error;
    }
  },
};

export default tokenService;

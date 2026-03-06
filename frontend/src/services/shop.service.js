import api from './api';

/**
 * Shop Service
 * Handles all shop-related API calls
 */

/**
 * Get current shop inventory
 * @returns {Promise<Object>} Shop inventory with items, discount info, and refresh time
 */
export const getInventory = async () => {
  const response = await api.get('/api/shop/inventory');
  return response.data;
};

/**
 * Purchase a Pokemon card from the shop
 * @param {string} pokemonName - Name of the Pokemon to purchase
 * @returns {Promise<Object>} Purchase result with card and remaining coins
 */
export const purchaseCard = async (pokemonName) => {
  const response = await api.post('/api/shop/purchase', {
    pokemon_name: pokemonName,
  });
  return response.data;
};

/**
 * Purchase game tokens from the shop
 * @param {number} quantity - Number of tokens to purchase (1-10)
 * @returns {Promise<Object>} Purchase result with tokens added and remaining coins
 */
export const purchaseTokens = async (quantity) => {
  const response = await api.post('/api/shop/tokens/purchase', {
    quantity,
  });
  return response.data;
};

const shopService = {
  getInventory,
  purchaseCard,
  purchaseTokens,
};

export default shopService;

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

/**
 * Purchase a game item (e.g. evolution stone) from the shop
 * @param {string} itemId - Item id (slug), e.g. 'thunder-stone'
 * @param {number} [quantity=1] - Quantity to purchase (1-99)
 * @returns {Promise<Object>} Purchase result with new quantity and remaining coins
 */
export const purchaseItem = async (itemId, quantity = 1) => {
  const response = await api.post('/api/shop/items/purchase', {
    item_id: itemId,
    quantity,
  });
  return response.data;
};

const shopService = {
  getInventory,
  purchaseCard,
  purchaseTokens,
  purchaseItem,
};

export default shopService;

import api from './api';

/**
 * Item Service
 * Handles item inventory and booster usage API calls
 */

/**
 * Get the user's item inventory and active boosts
 * @returns {Promise<Object>} { inventory: Array, activeBoosts: Array }
 */
export const getInventory = async () => {
  const response = await api.get('/api/items/inventory');
  return {
    inventory: response.data.inventory || [],
    activeBoosts: response.data.active_boosts || [],
  };
};

/**
 * Use (activate) a booster item from the inventory
 * @param {string} itemId - Item id (slug), e.g. 'attack-booster'
 * @returns {Promise<Object>} Response with message and active boost
 */
export const useItem = async (itemId) => {
  const response = await api.post(`/api/items/${itemId}/use`);
  return response.data;
};

export default {
  getInventory,
  useItem,
};

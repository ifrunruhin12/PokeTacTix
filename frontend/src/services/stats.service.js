import api from './api';

/**
 * Stats Service
 * Handles all statistics and profile-related API calls
 */

/**
 * Get player statistics
 * @returns {Promise<Object>} Player stats including battles, wins, losses, etc.
 */
export const getPlayerStats = async () => {
  const response = await api.get('/api/profile/stats');
  return response.data;
};

/**
 * Get battle history
 * @param {number} limit - Number of battles to retrieve (default 20, max 100)
 * @returns {Promise<Object>} Battle history with count
 */
export const getBattleHistory = async (limit = 20) => {
  const response = await api.get(`/api/profile/history?limit=${limit}`);
  return response.data;
};

/**
 * Get achievements
 * @returns {Promise<Object>} All achievements with unlock status
 */
export const getAchievements = async () => {
  const response = await api.get('/api/profile/achievements');
  return response.data;
};

/**
 * Check and unlock new achievements
 * @returns {Promise<Object>} Newly unlocked achievements
 */
export const checkAchievements = async () => {
  const response = await api.post('/api/profile/achievements/check');
  return response.data;
};

export default {
  getPlayerStats,
  getBattleHistory,
  getAchievements,
  checkAchievements,
};

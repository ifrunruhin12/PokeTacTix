import api from './api';

/**
 * Battle Service
 * Handles all battle-related API calls
 */

/**
 * Start a new battle
 * @param {string} mode - Battle mode ('1v1' or '5v5')
 * @returns {Promise<Object>} Battle state
 */
export const startBattle = async (mode = '5v5') => {
  const response = await api.post('/api/battle/start', { mode });
  return response.data;
};

/**
 * Submit a move in the current battle
 * @param {string} battleId - Battle session ID
 * @param {string} move - Move type ('attack', 'defend', 'pass', 'sacrifice', 'surrender')
 * @param {number} moveIdx - Move index (required for 'attack')
 * @returns {Promise<Object>} Battle result with updated state
 */
export const submitMove = async (battleId, move, moveIdx = null) => {
  const response = await api.post('/api/battle/move', {
    battleId,
    move,
    moveIdx
  });
  return response.data;
};

/**
 * Get current battle state
 * @param {string} battleId - Battle session ID
 * @returns {Promise<Object>} Current battle state
 */
export const getBattleState = async (battleId) => {
  const response = await api.get(`/api/battle/state?battleId=${battleId}`);
  return response.data;
};

/**
 * Switch active Pokemon
 * @param {string} battleId - Battle session ID
 * @param {number} newIdx - Index of Pokemon to switch to
 * @returns {Promise<Object>} Updated battle state
 */
export const switchPokemon = async (battleId, newIdx) => {
  const response = await api.post('/api/battle/switch', {
    battleId,
    newIdx
  });
  return response.data;
};

/**
 * Select reward Pokemon after 5v5 victory
 * @param {string} battleId - Battle session ID
 * @param {number} pokemonIdx - Index of AI Pokemon to select
 * @returns {Promise<Object>} Selected Pokemon card
 */
export const selectReward = async (battleId, pokemonIdx) => {
  const response = await api.post('/api/battle/select-reward', {
    battleId,
    pokemonIdx
  });
  return response.data;
};

export default {
  startBattle,
  submitMove,
  getBattleState,
  switchPokemon,
  selectReward
};

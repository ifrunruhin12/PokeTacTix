import api from './api';

/**
 * Card Service
 * Handles all card and deck management API calls
 */

/**
 * Get all cards in user's collection
 * @returns {Promise<Object>} Object with cards array
 */
export const getUserCards = async () => {
  const response = await api.get('/api/cards');
  return response.data;
};

/**
 * Get user's current deck (5 cards)
 * @returns {Promise<Array>} Array of deck cards
 */
export const getUserDeck = async () => {
  const response = await api.get('/api/cards/deck');
  return response.data.deck || [];
};

/**
 * Update user's deck configuration
 * @param {Array<number>} cardIds - Array of exactly 5 card IDs
 * @returns {Promise<Object>} Updated deck and success message
 */
export const updateDeck = async (cardIds) => {
  if (!Array.isArray(cardIds) || cardIds.length !== 5) {
    throw new Error('Deck must contain exactly 5 cards');
  }
  
  const response = await api.put('/api/cards/deck', {
    card_ids: cardIds
  });
  
  return response.data;
};

/**
 * Get a specific card by ID
 * @param {number} cardId - Card ID
 * @returns {Promise<Object>} Card data
 */
export const getCardById = async (cardId) => {
  const response = await api.get(`/api/cards/${cardId}`);
  return response.data.card;
};

/**
 * Calculate current stats for a card based on level
 * @param {Object} card - Player card object
 * @returns {Object} Current stats
 */
export const calculateCurrentStats = (card) => {
  const levelMultiplier = card.level - 1;
  
  const hp = Math.round(card.base_hp * (1.0 + levelMultiplier * 0.03));
  const attack = Math.round(card.base_attack * (1.0 + levelMultiplier * 0.02));
  const defense = Math.round(card.base_defense * (1.0 + levelMultiplier * 0.02));
  const speed = Math.round(card.base_speed * (1.0 + levelMultiplier * 0.01));
  const stamina = speed * 2;
  
  return {
    hp,
    hp_max: hp,
    attack,
    defense,
    speed,
    stamina,
    stamina_max: stamina
  };
};

/**
 * Transform backend card data to frontend format
 * @param {Object} card - Backend card object
 * @returns {Object} Frontend card object
 */
export const transformCardData = (card) => {
  const currentStats = calculateCurrentStats(card);
  
  return {
    id: card.id,
    name: card.pokemon_name,
    pokemon_name: card.pokemon_name,
    level: card.level,
    xp: card.xp,
    base_hp: card.base_hp,
    base_attack: card.base_attack,
    base_defense: card.base_defense,
    base_speed: card.base_speed,
    types: typeof card.types === 'string' ? JSON.parse(card.types) : card.types,
    moves: typeof card.moves === 'string' ? JSON.parse(card.moves) : card.moves,
    sprite: card.sprite,
    is_legendary: card.is_legendary,
    is_mythical: card.is_mythical,
    in_deck: card.in_deck,
    deck_position: card.deck_position,
    ...currentStats
  };
};

export default {
  getUserCards,
  getUserDeck,
  updateDeck,
  getCardById,
  calculateCurrentStats,
  transformCardData
};

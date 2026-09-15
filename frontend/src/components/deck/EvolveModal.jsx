import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import PropTypes from 'prop-types';
import { getEvolutionInfo, evolveCardWithItem } from '../../services/card.service';

/**
 * EvolveModal Component
 * Shows how a Pokemon can evolve (level-based or item-based) and lets the
 * player trigger item-based evolution from their deck.
 */
export default function EvolveModal({ isOpen, onClose, card, onEvolved }) {
  const [info, setInfo] = useState(null);
  const [loading, setLoading] = useState(false);
  const [evolving, setEvolving] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  useEffect(() => {
    if (isOpen && card) {
      loadInfo();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, card?.id]);

  const loadInfo = async () => {
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);
      const data = await getEvolutionInfo(card.id);
      setInfo(data);
    } catch (err) {
      setError(err.message || 'Failed to load evolution info');
    } finally {
      setLoading(false);
    }
  };

  const handleEvolve = async (itemId) => {
    try {
      setEvolving(true);
      setError(null);
      setSuccess(null);

      const result = await evolveCardWithItem(card.id, itemId);

      setSuccess(`${result.evolution.evolved_from} evolved into ${result.evolution.evolved_into}!`);

      // Refresh info to reflect the new species
      await loadInfo();

      if (onEvolved) {
        onEvolved(result.evolution);
      }
    } catch (err) {
      setError(err.message || 'Evolution failed');
    } finally {
      setEvolving(false);
    }
  };

  const handleClose = () => {
    if (!evolving) {
      onClose();
    }
  };

  return (
    <AnimatePresence>
      {isOpen && card && (
        <motion.div
          className="fixed inset-0 bg-black/70 flex items-center justify-center p-4 z-50"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onClick={handleClose}
        >
          <motion.div
            className="bg-gray-800 rounded-lg border border-gray-700 max-w-lg w-full p-6 max-h-[85vh] overflow-y-auto"
            initial={{ scale: 0.9, y: 20 }}
            animate={{ scale: 1, y: 0 }}
            exit={{ scale: 0.9, y: 20 }}
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-bold text-white">
                Evolution — {card.pokemon_name}
              </h2>
              <button
                onClick={handleClose}
                disabled={evolving}
                className="text-gray-400 hover:text-white text-xl font-bold disabled:opacity-50"
              >
                ✕
              </button>
            </div>

            {loading && (
              <div className="py-8 text-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-3"></div>
                <p className="text-gray-400">Checking evolution possibilities...</p>
              </div>
            )}

            {!loading && error && !success && (
              <div className="bg-red-900/50 border border-red-500 rounded-lg p-4 mb-4">
                <p className="text-red-200">{error}</p>
              </div>
            )}

            {success && (
              <div className="bg-green-900/50 border border-green-500 rounded-lg p-4 mb-4">
                <p className="text-green-200 font-semibold">✨ {success}</p>
              </div>
            )}

            {!loading && info && (
              <>
                {(info.options || []).length === 0 ? (
                  <div className="text-center py-6">
                    <p className="text-gray-400 text-lg mb-2">
                      {info.pokemon_name} cannot evolve further.
                    </p>
                    <p className="text-gray-500 text-sm">This Pokemon has no evolution path.</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {(info.options || []).map((option, index) => (
                      <div
                        key={index}
                        className="bg-gray-700/50 rounded-lg p-4 border border-gray-600"
                      >
                        <div className="flex items-center gap-4">
                          {option.target_sprite ? (
                            <img
                              src={option.target_sprite}
                              alt={option.target_name}
                              className="h-16 w-16 object-contain"
                              onError={(e) => { e.target.style.display = 'none'; }}
                            />
                          ) : (
                            <span className="text-3xl">❓</span>
                          )}
                          <div className="flex-1">
                            <div className="flex items-center gap-2 flex-wrap">
                              <span className="text-white font-bold">{option.target_name}</span>
                              {option.method === 'item' ? (
                                <span className="bg-amber-500 text-white text-xs font-bold px-2 py-0.5 rounded-full">
                                  ITEM
                                </span>
                              ) : (
                                <span className="bg-blue-500 text-white text-xs font-bold px-2 py-0.5 rounded-full">
                                  LEVEL
                                </span>
                              )}
                            </div>

                            {option.method === 'level' && (
                              <p className="text-gray-400 text-sm mt-1">
                                {option.eligible
                                  ? `Evolved automatically on the next battle reward (level ${option.min_level} reached).`
                                  : `Reach level ${option.min_level} through battles — current level ${info.level}.`}
                              </p>
                            )}

                            {option.method === 'item' && (
                              <p className="text-gray-400 text-sm mt-1">
                                Requires <span className="text-amber-300 font-semibold">{option.required_item_name || option.required_item}</span>
                                {' '}(owned: {option.owned_quantity})
                              </p>
                            )}
                          </div>

                          {option.method === 'item' && (
                            <button
                              onClick={() => handleEvolve(option.required_item)}
                              disabled={!option.eligible || evolving}
                              className={`px-4 py-2 rounded-lg font-bold text-sm whitespace-nowrap transition-colors ${
                                option.eligible && !evolving
                                  ? 'bg-green-600 hover:bg-green-700 text-white'
                                  : 'bg-gray-600 text-gray-400 cursor-not-allowed'
                              }`}
                            >
                              {!option.eligible
                                ? 'Item required'
                                : evolving
                                  ? 'Evolving...'
                                  : 'Evolve'}
                            </button>
                          )}
                        </div>

                        {option.method === 'item' && !option.eligible && (
                          <p className="text-xs text-gray-500 mt-2">
                            Purchase {option.required_item_name || option.required_item} in the Shop's Items category.
                          </p>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </>
            )}

            <div className="mt-6 text-center">
              <button
                onClick={handleClose}
                disabled={evolving}
                className="px-6 py-2 bg-gray-600 hover:bg-gray-700 rounded-lg font-semibold text-white transition-colors disabled:opacity-50"
              >
                Close
              </button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}

EvolveModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  card: PropTypes.shape({
    id: PropTypes.number.isRequired,
    pokemon_name: PropTypes.string.isRequired,
  }),
  onEvolved: PropTypes.func,
};

import PropTypes from 'prop-types';
import { motion, AnimatePresence } from 'framer-motion';

/**
 * MoveSelector Component
 * Modal for selecting attack moves
 * Displays move details (name, type, power, stamina cost)
 */
const MoveSelector = ({ 
  isOpen = false, 
  moves = [], 
  currentStamina = 0,
  onSelectMove, 
  onCancel 
}) => {
  // Get type color
  const getTypeColor = (type) => {
    const typeColors = {
      normal: '#A8A878',
      fire: '#F08030',
      water: '#6890F0',
      electric: '#F8D030',
      grass: '#78C850',
      ice: '#98D8D8',
      fighting: '#C03028',
      poison: '#A040A0',
      ground: '#E0C068',
      flying: '#A890F0',
      psychic: '#F85888',
      bug: '#A8B820',
      rock: '#B8A038',
      ghost: '#705898',
      dragon: '#7038F8',
      dark: '#705848',
      steel: '#B8B8D0',
      fairy: '#EE99AC'
    };
    return typeColors[type?.toLowerCase()] || '#A8A878';
  };

  // Check if move can be used
  const canUseMove = (move) => {
    return currentStamina >= (move.stamina_cost || 0);
  };

  if (!isOpen) return null;

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40"
            onClick={onCancel}
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            transition={{ type: 'spring', stiffness: 300, damping: 30 }}
            className="fixed inset-0 flex items-center justify-center z-50 p-4"
          >
            <div className="bg-gray-800 rounded-xl shadow-2xl max-w-md w-full max-h-[80vh] overflow-hidden border-2 border-gray-700">
              {/* Header */}
              <div className="bg-gradient-to-r from-blue-600 to-purple-600 p-4">
                <h2 className="text-2xl font-bold text-white text-center">
                  ⚔️ Select Your Move
                </h2>
                <p className="text-blue-100 text-sm text-center mt-1">
                  Current Stamina: {currentStamina}
                </p>
              </div>

              {/* Moves list */}
              <div className="p-4 space-y-3 max-h-96 overflow-y-auto">
                {moves && moves.length > 0 ? (
                  moves.map((move, index) => {
                    const usable = canUseMove(move);
                    return (
                      <motion.button
                        key={index}
                        whileHover={usable ? { scale: 1.02, x: 5 } : {}}
                        whileTap={usable ? { scale: 0.98 } : {}}
                        onClick={() => usable && onSelectMove(index)}
                        disabled={!usable}
                        className={`w-full p-4 rounded-lg border-2 text-left transition-all ${
                          usable
                            ? 'bg-gray-700 border-gray-600 hover:border-blue-500 hover:bg-gray-650 cursor-pointer'
                            : 'bg-gray-800 border-gray-700 opacity-50 cursor-not-allowed'
                        }`}
                      >
                        <div className="flex items-start justify-between mb-2">
                          <h3 className="text-lg font-bold text-white">
                            {move.name || 'Unknown Move'}
                          </h3>
                          {!usable && (
                            <span className="text-xs text-red-400 font-semibold bg-red-900/30 px-2 py-1 rounded">
                              Not enough stamina
                            </span>
                          )}
                        </div>

                        <div className="flex items-center gap-3 flex-wrap">
                          {/* Type badge */}
                          {move.type && (
                            <span
                              className="text-xs px-3 py-1 rounded-full text-white font-semibold"
                              style={{ backgroundColor: getTypeColor(move.type) }}
                            >
                              {move.type.toUpperCase()}
                            </span>
                          )}

                          {/* Power */}
                          {move.power !== undefined && (
                            <span className="text-sm text-orange-400 font-semibold">
                              ⚡ Power: {move.power}
                            </span>
                          )}

                          {/* Stamina cost */}
                          {move.stamina_cost !== undefined && (
                            <span className={`text-sm font-semibold ${
                              usable ? 'text-blue-400' : 'text-red-400'
                            }`}>
                              💧 Cost: {move.stamina_cost}
                            </span>
                          )}
                        </div>

                        {/* Move description if available */}
                        {move.description && (
                          <p className="text-xs text-gray-400 mt-2">
                            {move.description}
                          </p>
                        )}
                      </motion.button>
                    );
                  })
                ) : (
                  <p className="text-gray-400 text-center py-8">
                    No moves available
                  </p>
                )}
              </div>

              {/* Footer */}
              <div className="p-4 bg-gray-900 border-t border-gray-700">
                <button
                  onClick={onCancel}
                  className="w-full py-3 px-4 bg-gray-700 hover:bg-gray-600 text-white font-semibold rounded-lg transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
};

MoveSelector.propTypes = {
  isOpen: PropTypes.bool,
  moves: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string,
      type: PropTypes.string,
      power: PropTypes.number,
      stamina_cost: PropTypes.number,
      description: PropTypes.string
    })
  ),
  currentStamina: PropTypes.number,
  onSelectMove: PropTypes.func.isRequired,
  onCancel: PropTypes.func.isRequired
};

export default MoveSelector;

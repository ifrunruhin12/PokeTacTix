import { motion, AnimatePresence } from 'framer-motion';
import PropTypes from 'prop-types';

/**
 * Type color mapping for Pokemon types
 */
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
  fairy: '#EE99AC',
};

/**
 * PurchaseModal Component
 * Modal for confirming Pokemon card purchases
 */
export default function PurchaseModal({ 
  isOpen, 
  onClose, 
  item, 
  userCoins, 
  onConfirm, 
  isProcessing,
  error,
}) {
  if (!item) return null;

  const canAfford = userCoins >= item.price;
  const remainingCoins = userCoins - item.price;

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/70 z-40 backdrop-blur-sm"
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4"
          >
            <div className="bg-gray-800 rounded-lg shadow-2xl max-w-md w-full border-2 border-gray-700 overflow-hidden">
              {/* Header */}
              <div className="bg-gradient-to-r from-blue-600 to-purple-600 p-4">
                <h2 className="text-2xl font-bold text-white text-center">
                  Confirm Purchase
                </h2>
              </div>

              {/* Content */}
              <div className="p-6">
                {/* Pokemon Display */}
                <div className="flex flex-col items-center mb-6">
                  <img
                    src={item.sprite}
                    alt={item.pokemon_name}
                    className="w-32 h-32 object-contain mb-3"
                    onError={(e) => {
                      e.target.src = '/assets/pokeball.png';
                    }}
                  />
                  <h3 className="text-2xl font-bold capitalize mb-2">
                    {item.pokemon_name}
                  </h3>

                  {/* Types */}
                  <div className="flex gap-2 mb-3">
                    {item.types.map((type) => (
                      <span
                        key={type}
                        className="px-3 py-1 rounded-full text-sm font-semibold text-white capitalize"
                        style={{ backgroundColor: typeColors[type.toLowerCase()] || '#777' }}
                      >
                        {type}
                      </span>
                    ))}
                  </div>

                  {/* Rarity Badge */}
                  {(item.is_legendary || item.is_mythical) && (
                    <span className={`px-4 py-1 rounded-full text-sm font-bold ${
                      item.is_legendary 
                        ? 'bg-yellow-500 text-black' 
                        : 'bg-gradient-to-r from-purple-500 to-pink-500 text-white'
                    }`}>
                      {item.is_legendary ? '⭐ LEGENDARY' : '✨ MYTHICAL'}
                    </span>
                  )}
                </div>

                {/* Stats */}
                <div className="bg-gray-900 rounded-lg p-4 mb-4">
                  <h4 className="text-sm font-semibold text-gray-400 mb-2">Base Stats</h4>
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <span className="text-gray-400 text-sm">HP:</span>
                      <span className="text-green-400 font-bold ml-2">{item.base_hp}</span>
                    </div>
                    <div>
                      <span className="text-gray-400 text-sm">Attack:</span>
                      <span className="text-red-400 font-bold ml-2">{item.base_attack}</span>
                    </div>
                    <div>
                      <span className="text-gray-400 text-sm">Defense:</span>
                      <span className="text-yellow-400 font-bold ml-2">{item.base_defense}</span>
                    </div>
                    <div>
                      <span className="text-gray-400 text-sm">Speed:</span>
                      <span className="text-purple-400 font-bold ml-2">{item.base_speed}</span>
                    </div>
                  </div>
                </div>

                {/* Price Info */}
                <div className="bg-gray-900 rounded-lg p-4 mb-4">
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-gray-400">Price:</span>
                    <span className="text-yellow-400 font-bold text-xl">
                      {item.price} 🪙
                    </span>
                  </div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-gray-400">Your Coins:</span>
                    <span className="text-white font-bold">{userCoins} 🪙</span>
                  </div>
                  <div className="border-t border-gray-700 pt-2 mt-2">
                    <div className="flex justify-between items-center">
                      <span className="text-gray-400">After Purchase:</span>
                      <span className={`font-bold ${canAfford ? 'text-green-400' : 'text-red-400'}`}>
                        {canAfford ? `${remainingCoins} 🪙` : 'Insufficient Coins'}
                      </span>
                    </div>
                  </div>
                </div>

                {/* Error Message */}
                {error && (
                  <div className="bg-red-900/50 border border-red-500 rounded-lg p-3 mb-4">
                    <p className="text-red-200 text-sm">{error}</p>
                  </div>
                )}

                {/* Buttons */}
                <div className="flex gap-3">
                  <button
                    onClick={onClose}
                    disabled={isProcessing}
                    className="flex-1 bg-gray-700 hover:bg-gray-600 text-white font-bold py-3 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={onConfirm}
                    disabled={!canAfford || isProcessing}
                    className={`flex-1 font-bold py-3 rounded-lg transition-colors ${
                      canAfford && !isProcessing
                        ? 'bg-blue-600 hover:bg-blue-700 text-white'
                        : 'bg-gray-700 text-gray-500 cursor-not-allowed'
                    }`}
                  >
                    {isProcessing ? (
                      <span className="flex items-center justify-center">
                        <svg className="animate-spin h-5 w-5 mr-2" viewBox="0 0 24 24">
                          <circle
                            className="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            strokeWidth="4"
                            fill="none"
                          />
                          <path
                            className="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                          />
                        </svg>
                        Processing...
                      </span>
                    ) : (
                      'Confirm Purchase'
                    )}
                  </button>
                </div>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}

PurchaseModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  item: PropTypes.shape({
    pokemon_name: PropTypes.string.isRequired,
    price: PropTypes.number.isRequired,
    rarity: PropTypes.string.isRequired,
    is_legendary: PropTypes.bool.isRequired,
    is_mythical: PropTypes.bool.isRequired,
    sprite: PropTypes.string.isRequired,
    types: PropTypes.arrayOf(PropTypes.string).isRequired,
    base_hp: PropTypes.number.isRequired,
    base_attack: PropTypes.number.isRequired,
    base_defense: PropTypes.number.isRequired,
    base_speed: PropTypes.number.isRequired,
  }),
  userCoins: PropTypes.number.isRequired,
  onConfirm: PropTypes.func.isRequired,
  isProcessing: PropTypes.bool,
  error: PropTypes.string,
};

PurchaseModal.defaultProps = {
  item: null,
  isProcessing: false,
  error: null,
};

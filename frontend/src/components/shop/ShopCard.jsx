import { motion } from 'framer-motion';
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
 * ShopCard Component
 * Displays a Pokemon card in the shop with price and rarity indicators
 */
export default function ShopCard({ item, onPurchase, userCoins, isOwned, originalPrice }) {
  const canAfford = userCoins >= item.price;
  const isDiscounted = originalPrice && originalPrice > item.price;
  
  // Determine border color based on rarity
  const getBorderClass = () => {
    if (item.is_legendary) {
      return 'border-yellow-400 shadow-yellow-400/50';
    }
    if (item.is_mythical) {
      return 'border-purple-500 shadow-purple-500/50 bg-gradient-to-br from-purple-900/20 to-pink-900/20';
    }
    if (item.rarity === 'rare') {
      return 'border-blue-400 shadow-blue-400/30';
    }
    if (item.rarity === 'uncommon') {
      return 'border-green-400 shadow-green-400/30';
    }
    return 'border-gray-600 shadow-gray-600/30';
  };

  // Get rarity badge color
  const getRarityBadgeClass = () => {
    if (item.is_legendary) return 'bg-yellow-500 text-black';
    if (item.is_mythical) return 'bg-gradient-to-r from-purple-500 to-pink-500 text-white';
    if (item.rarity === 'rare') return 'bg-blue-500 text-white';
    if (item.rarity === 'uncommon') return 'bg-green-500 text-white';
    return 'bg-gray-500 text-white';
  };

  return (
    <motion.div
      className={`relative bg-gray-800 rounded-lg border-2 ${getBorderClass()} overflow-hidden shadow-lg`}
      whileHover={{ scale: 1.05, y: -5 }}
      transition={{ duration: 0.2 }}
    >
      {/* Rare Appearance Badge */}
      {(item.is_legendary || item.is_mythical) && (
        <div className="absolute top-2 left-2 z-10">
          <span className="bg-red-600 text-white text-xs font-bold px-2 py-1 rounded-full animate-pulse">
            ⭐ RARE
          </span>
        </div>
      )}

      {/* Owned Indicator */}
      {isOwned && (
        <div className="absolute top-2 right-2 z-10">
          <span className="bg-green-600 text-white text-xs font-bold px-2 py-1 rounded-full">
            ✓ OWNED
          </span>
        </div>
      )}

      {/* Card Content */}
      <div className="p-4">
        {/* Pokemon Sprite */}
        <div className="flex justify-center mb-3">
          <img
            src={item.sprite}
            alt={item.pokemon_name}
            className="w-24 h-24 object-contain"
            onError={(e) => {
              e.target.src = '/assets/pokeball.png';
            }}
          />
        </div>

        {/* Pokemon Name */}
        <h3 className="text-xl font-bold text-center mb-2 capitalize">
          {item.pokemon_name}
        </h3>

        {/* Types */}
        <div className="flex justify-center gap-2 mb-3">
          {item.types.map((type) => (
            <span
              key={type}
              className="px-3 py-1 rounded-full text-xs font-semibold text-white capitalize"
              style={{ backgroundColor: typeColors[type.toLowerCase()] || '#777' }}
            >
              {type}
            </span>
          ))}
        </div>

        {/* Stats Preview */}
        <div className="grid grid-cols-2 gap-2 mb-3 text-sm">
          <div className="bg-gray-700 rounded px-2 py-1">
            <span className="text-gray-400">HP:</span>
            <span className="text-green-400 font-bold ml-1">{item.base_hp}</span>
          </div>
          <div className="bg-gray-700 rounded px-2 py-1">
            <span className="text-gray-400">ATK:</span>
            <span className="text-red-400 font-bold ml-1">{item.base_attack}</span>
          </div>
          <div className="bg-gray-700 rounded px-2 py-1">
            <span className="text-gray-400">DEF:</span>
            <span className="text-yellow-400 font-bold ml-1">{item.base_defense}</span>
          </div>
          <div className="bg-gray-700 rounded px-2 py-1">
            <span className="text-gray-400">SPD:</span>
            <span className="text-purple-400 font-bold ml-1">{item.base_speed}</span>
          </div>
        </div>

        {/* Rarity Badge */}
        <div className="flex justify-center mb-3">
          <span className={`${getRarityBadgeClass()} text-xs font-bold px-3 py-1 rounded-full uppercase`}>
            {item.is_legendary ? 'Legendary' : item.is_mythical ? 'Mythical' : item.rarity}
          </span>
        </div>

        {/* Price and Purchase Button */}
        <div className="border-t border-gray-700 pt-3">
          <div className="flex items-center justify-between mb-2">
            <span className="text-gray-400 text-sm">Price:</span>
            <div className="flex items-center gap-2">
              {isDiscounted && (
                <span className="text-gray-500 line-through text-sm">
                  {originalPrice} 🪙
                </span>
              )}
              <span className={`font-bold text-lg ${isDiscounted ? 'text-green-400' : 'text-yellow-400'}`}>
                {item.price} 🪙
              </span>
            </div>
          </div>

          <button
            onClick={() => onPurchase(item)}
            disabled={!canAfford}
            className={`w-full py-2 rounded-lg font-bold transition-colors ${
              canAfford
                ? 'bg-blue-600 hover:bg-blue-700 text-white'
                : 'bg-gray-700 text-gray-500 cursor-not-allowed'
            }`}
          >
            {canAfford ? 'Purchase' : 'Insufficient Coins'}
          </button>
        </div>
      </div>
    </motion.div>
  );
}

ShopCard.propTypes = {
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
  }).isRequired,
  onPurchase: PropTypes.func.isRequired,
  userCoins: PropTypes.number.isRequired,
  isOwned: PropTypes.bool,
  originalPrice: PropTypes.number,
};

ShopCard.defaultProps = {
  isOwned: false,
  originalPrice: null,
};

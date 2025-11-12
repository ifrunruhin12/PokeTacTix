import { motion } from 'framer-motion';
import PropTypes from 'prop-types';

/**
 * PokemonCard Component
 * Displays a Pokemon card with stats, level, XP, and animations
 * Supports active state, face-down state, and rarity-based borders
 */
const PokemonCard = ({ 
  pokemon, 
  isActive = false, 
  isFaceDown = false,
  isKnockedOut = false,
  onSelect,
  className = ''
}) => {
  // Determine border color based on rarity
  const getBorderClass = () => {
    if (isKnockedOut) return 'border-gray-500';
    if (pokemon?.is_legendary) return 'border-yellow-400 shadow-yellow-400/50';
    if (pokemon?.is_mythical) return 'border-purple-500 shadow-purple-500/50';
    if (isActive) return 'border-blue-400 shadow-blue-400/50';
    return 'border-gray-300';
  };

  // Get type colors for display
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

  // If face-down, show Pokeball
  if (isFaceDown) {
    return (
      <motion.div
        className={`relative w-40 h-56 rounded-lg border-4 border-gray-400 bg-gradient-to-br from-gray-100 to-gray-200 shadow-lg flex items-center justify-center ${className}`}
        whileHover={{ scale: 1.05, y: -5 }}
        transition={{ duration: 0.2 }}
      >
        <img 
          src="/assets/pokeball.png" 
          alt="Hidden Pokemon" 
          className="w-20 h-20 opacity-50"
        />
      </motion.div>
    );
  }

  // If no pokemon data, show empty card
  if (!pokemon) {
    return (
      <div className={`w-40 h-56 rounded-lg border-4 border-dashed border-gray-300 bg-gray-50 flex items-center justify-center ${className}`}>
        <span className="text-gray-400 text-sm">Empty Slot</span>
      </div>
    );
  }

  return (
    <motion.div
      className={`relative w-40 h-56 rounded-lg border-4 ${getBorderClass()} bg-gradient-to-br from-white to-gray-50 shadow-lg overflow-hidden cursor-pointer ${className}`}
      whileHover={!isKnockedOut ? { scale: 1.05, y: -5 } : {}}
      animate={isActive ? { 
        boxShadow: '0 0 20px rgba(59, 130, 246, 0.5)',
        scale: 1.02
      } : {}}
      transition={{ duration: 0.2 }}
      onClick={onSelect}
      style={{
        filter: isKnockedOut ? 'grayscale(100%)' : 'none',
        opacity: isKnockedOut ? 0.5 : 1
      }}
    >
      {/* Rarity indicator for legendary/mythical */}
      {pokemon.is_legendary && (
        <div className="absolute top-1 right-1 bg-yellow-400 text-yellow-900 text-xs px-2 py-0.5 rounded-full font-bold z-10">
          ⭐ Legendary
        </div>
      )}
      {pokemon.is_mythical && (
        <div className="absolute top-1 right-1 bg-purple-500 text-white text-xs px-2 py-0.5 rounded-full font-bold z-10">
          ✨ Mythical
        </div>
      )}

      {/* Knocked out overlay */}
      {isKnockedOut && (
        <div className="absolute inset-0 flex items-center justify-center z-20 bg-black/30">
          <div className="text-red-500 text-6xl font-bold">✕</div>
        </div>
      )}

      {/* Card content */}
      <div className="p-2 flex flex-col h-full">
        {/* Pokemon name */}
        <h3 className="text-sm font-bold text-gray-800 text-center truncate">
          {pokemon.name || pokemon.pokemon_name}
        </h3>

        {/* Type badges */}
        <div className="flex gap-1 justify-center my-1">
          {pokemon.types?.map((type, idx) => (
            <span
              key={idx}
              className="text-xs px-2 py-0.5 rounded-full text-white font-semibold"
              style={{ backgroundColor: getTypeColor(type) }}
            >
              {type}
            </span>
          ))}
        </div>

        {/* Pokemon sprite */}
        <div className="flex-1 flex items-center justify-center my-1">
          {pokemon.sprite ? (
            <img 
              src={pokemon.sprite} 
              alt={pokemon.name || pokemon.pokemon_name}
              className="w-20 h-20 object-contain"
            />
          ) : (
            <div className="w-20 h-20 bg-gray-200 rounded-full flex items-center justify-center">
              <span className="text-gray-400 text-xs">No Image</span>
            </div>
          )}
        </div>

        {/* Stats section */}
        <div className="space-y-0.5 text-xs">
          {/* HP */}
          <div className="flex items-center gap-1">
            <span className="text-gray-600 w-8">HP:</span>
            <div className="flex-1 bg-gray-200 rounded-full h-2 overflow-hidden">
              <motion.div
                className="h-full bg-green-500"
                initial={{ width: 0 }}
                animate={{ 
                  width: `${Math.max(0, Math.min(100, (pokemon.hp / pokemon.hp_max) * 100))}%`,
                  backgroundColor: (pokemon.hp / pokemon.hp_max) < 0.3 ? '#ef4444' : '#10b981'
                }}
                transition={{ duration: 0.5 }}
              />
            </div>
            <span className="text-gray-700 w-10 text-right">{pokemon.hp}/{pokemon.hp_max}</span>
          </div>

          {/* Stamina */}
          {pokemon.stamina !== undefined && (
            <div className="flex items-center gap-1">
              <span className="text-gray-600 w-8">STA:</span>
              <div className="flex-1 bg-gray-200 rounded-full h-2 overflow-hidden">
                <motion.div
                  className="h-full bg-blue-500"
                  initial={{ width: 0 }}
                  animate={{ width: `${Math.max(0, Math.min(100, (pokemon.stamina / pokemon.stamina_max) * 100))}%` }}
                  transition={{ duration: 0.5 }}
                />
              </div>
              <span className="text-gray-700 w-10 text-right">{pokemon.stamina}/{pokemon.stamina_max}</span>
            </div>
          )}

          {/* Attack, Defense, Speed */}
          <div className="flex justify-between text-gray-700">
            <span>ATK: {pokemon.attack}</span>
            <span>DEF: {pokemon.defense}</span>
            <span>SPD: {pokemon.speed}</span>
          </div>

          {/* Level and XP */}
          {pokemon.level !== undefined && (
            <div className="text-center">
              <span className="font-semibold text-gray-800">Lv. {pokemon.level}</span>
              {pokemon.xp !== undefined && pokemon.level < 50 && (
                <div className="mt-0.5">
                  <div className="bg-gray-200 rounded-full h-1.5 overflow-hidden">
                    <motion.div
                      className="h-full bg-purple-500"
                      initial={{ width: 0 }}
                      animate={{ width: `${(pokemon.xp / (100 * pokemon.level)) * 100}%` }}
                      transition={{ duration: 0.5 }}
                    />
                  </div>
                  <span className="text-xs text-gray-500">
                    {pokemon.xp}/{100 * pokemon.level} XP
                  </span>
                </div>
              )}
              {pokemon.level === 50 && (
                <span className="text-xs text-purple-600 font-semibold">MAX LEVEL</span>
              )}
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
};

PokemonCard.propTypes = {
  pokemon: PropTypes.shape({
    name: PropTypes.string,
    pokemon_name: PropTypes.string,
    types: PropTypes.arrayOf(PropTypes.string),
    sprite: PropTypes.string,
    hp: PropTypes.number,
    hp_max: PropTypes.number,
    stamina: PropTypes.number,
    stamina_max: PropTypes.number,
    attack: PropTypes.number,
    defense: PropTypes.number,
    speed: PropTypes.number,
    level: PropTypes.number,
    xp: PropTypes.number,
    is_legendary: PropTypes.bool,
    is_mythical: PropTypes.bool
  }),
  isActive: PropTypes.bool,
  isFaceDown: PropTypes.bool,
  isKnockedOut: PropTypes.bool,
  onSelect: PropTypes.func,
  className: PropTypes.string
};

export default PokemonCard;

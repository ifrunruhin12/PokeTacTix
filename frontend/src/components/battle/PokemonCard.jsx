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
  className = '',
  compact = false
}) => {
  // Determine border color based on rarity
  const getBorderClass = () => {
    if (isKnockedOut) return 'border-gray-600';
    if (pokemon?.is_legendary) return 'border-yellow-500 shadow-lg shadow-yellow-500/30';
    if (pokemon?.is_mythical) return 'border-purple-500 shadow-lg shadow-purple-500/30';
    if (isActive) return 'border-blue-500 shadow-lg shadow-blue-500/30';
    return 'border-gray-600';
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
        className={`relative rounded-lg border-4 border-gray-600 bg-gradient-to-br from-gray-800 to-gray-900 shadow-lg flex items-center justify-center ${className}`}
        whileHover={{ scale: 1.05, y: -5 }}
        transition={{ duration: 0.2 }}
        style={{ width: '100%', height: '100%', minWidth: '96px', minHeight: '128px' }}
      >
        {/* Pokeball design */}
        <div className="relative w-16 h-16">
          <div className="absolute inset-0 rounded-full bg-gradient-to-b from-red-500 to-red-600 border-4 border-gray-900"></div>
          <div className="absolute bottom-0 left-0 right-0 h-1/2 rounded-b-full bg-gradient-to-b from-gray-100 to-gray-300 border-4 border-gray-900"></div>
          <div className="absolute top-1/2 left-0 right-0 h-1 bg-gray-900 transform -translate-y-1/2"></div>
          <div className="absolute top-1/2 left-1/2 w-6 h-6 rounded-full bg-gray-100 border-4 border-gray-900 transform -translate-x-1/2 -translate-y-1/2">
            <div className="absolute inset-1 rounded-full bg-gray-300"></div>
          </div>
        </div>
      </motion.div>
    );
  }

  // If no pokemon data, show empty card
  if (!pokemon) {
    return (
      <div 
        className={`rounded-lg border-4 border-dashed border-gray-600 bg-gray-800/50 flex items-center justify-center ${className}`}
        style={{ width: '100%', height: '100%', minWidth: '96px', minHeight: '128px' }}
      >
        <span className="text-gray-500 text-sm">Empty</span>
      </div>
    );
  }

  return (
    <motion.div
      className={`relative rounded-lg border-4 ${getBorderClass()} bg-gradient-to-br from-gray-800 via-gray-850 to-gray-900 shadow-xl overflow-hidden cursor-pointer ${className}`}
      whileHover={!isKnockedOut ? { scale: 1.05, y: -5 } : {}}
      animate={isActive ? { 
        boxShadow: '0 0 25px rgba(59, 130, 246, 0.6)',
        scale: 1.02
      } : {}}
      transition={{ duration: 0.2 }}
      onClick={onSelect}
      style={{
        filter: isKnockedOut ? 'grayscale(100%)' : 'none',
        opacity: isKnockedOut ? 0.4 : 1,
        width: '100%',
        height: '100%',
        minWidth: compact ? '96px' : '128px',
        minHeight: compact ? '128px' : '176px'
      }}
    >
      {/* Rarity indicator for legendary/mythical */}
      {!compact && pokemon.is_legendary && (
        <div className="absolute top-1 right-1 bg-yellow-500 text-yellow-900 text-xs px-2 py-0.5 rounded-full font-bold z-10">
          ⭐ Legendary
        </div>
      )}
      {!compact && pokemon.is_mythical && (
        <div className="absolute top-1 right-1 bg-purple-500 text-white text-xs px-2 py-0.5 rounded-full font-bold z-10">
          ✨ Mythical
        </div>
      )}
      {compact && (pokemon.is_legendary || pokemon.is_mythical) && (
        <div className="absolute top-1 right-1 text-lg z-10">
          {pokemon.is_legendary ? '⭐' : '✨'}
        </div>
      )}

      {/* Knocked out overlay */}
      {isKnockedOut && (
        <div className="absolute inset-0 flex items-center justify-center z-20 bg-black/30">
          <div className="text-red-500 text-6xl font-bold">✕</div>
        </div>
      )}

      {/* Card content */}
      <div className={`${compact ? 'p-1' : 'p-2'} flex flex-col h-full`}>
        {/* Pokemon name - always visible, no truncation */}
        <h3 className={`${compact ? 'text-xs' : 'text-sm'} font-bold text-white text-center break-words leading-tight ${compact ? 'min-h-[1.5rem]' : 'min-h-[2rem]'} flex items-center justify-center`}>
          {pokemon.name || pokemon.pokemon_name}
        </h3>

        {/* Type badges */}
        {!compact && (
          <div className="flex gap-1 justify-center my-1 flex-wrap">
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
        )}

        {/* Pokemon sprite */}
        <div className={`flex-1 flex items-center justify-center ${compact ? 'my-0.5' : 'my-1'}`}>
          {pokemon.sprite ? (
            <img 
              src={pokemon.sprite} 
              alt={pokemon.name || pokemon.pokemon_name}
              className={`${compact ? 'w-12 h-12' : 'w-16 h-16'} object-contain`}
            />
          ) : (
            <div className={`${compact ? 'w-12 h-12' : 'w-16 h-16'} bg-gray-700 rounded-full flex items-center justify-center`}>
              <span className="text-gray-500 text-xs">?</span>
            </div>
          )}
        </div>

        {/* Stats section */}
        <div className={`space-y-0.5 ${compact ? 'text-[10px]' : 'text-xs'}`}>
          {/* HP */}
          <div className="flex items-center gap-1">
            <span className="text-gray-400 w-7 text-xs">HP:</span>
            <div className="flex-1 bg-gray-700 rounded-full h-2 overflow-hidden">
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
            {!compact && <span className="text-gray-300 text-xs whitespace-nowrap">{pokemon.hp}/{pokemon.hp_max}</span>}
          </div>

          {/* Stamina */}
          {!compact && pokemon.stamina !== undefined && (
            <div className="flex items-center gap-1">
              <span className="text-gray-400 w-7 text-xs">STA:</span>
              <div className="flex-1 bg-gray-700 rounded-full h-2 overflow-hidden">
                <motion.div
                  className="h-full bg-blue-500"
                  initial={{ width: 0 }}
                  animate={{ width: `${Math.max(0, Math.min(100, (pokemon.stamina / pokemon.stamina_max) * 100))}%` }}
                  transition={{ duration: 0.5 }}
                />
              </div>
              <span className="text-gray-300 text-xs whitespace-nowrap">{pokemon.stamina}/{pokemon.stamina_max}</span>
            </div>
          )}

          {/* Attack, Defense, Speed */}
          {!compact && (
            <div className="flex justify-between text-gray-300 text-xs">
              <span>ATK: {pokemon.attack}</span>
              <span>DEF: {pokemon.defense}</span>
              <span>SPD: {pokemon.speed}</span>
            </div>
          )}

          {/* Level and XP */}
          {!compact && pokemon.level !== undefined && (
            <div className="text-center">
              <span className="font-semibold text-white text-xs">Lv. {pokemon.level}</span>
              {pokemon.xp !== undefined && pokemon.level < 50 && (
                <div className="mt-0.5">
                  <div className="bg-gray-700 rounded-full h-1.5 overflow-hidden">
                    <motion.div
                      className="h-full bg-purple-500"
                      initial={{ width: 0 }}
                      animate={{ width: `${(pokemon.xp / (100 * pokemon.level)) * 100}%` }}
                      transition={{ duration: 0.5 }}
                    />
                  </div>
                  <span className="text-xs text-gray-400">
                    {pokemon.xp}/{100 * pokemon.level} XP
                  </span>
                </div>
              )}
              {pokemon.level === 50 && (
                <span className="text-xs text-purple-400 font-semibold">MAX LEVEL</span>
              )}
            </div>
          )}
          
          {/* Compact mode: just show level */}
          {compact && pokemon.level !== undefined && (
            <div className="text-center">
              <span className="font-semibold text-white text-xs">Lv. {pokemon.level}</span>
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
  className: PropTypes.string,
  compact: PropTypes.bool
};

export default PokemonCard;

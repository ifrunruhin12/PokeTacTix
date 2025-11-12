import { motion } from 'framer-motion';
import PropTypes from 'prop-types';
import PokemonCard from './PokemonCard';

/**
 * FlippableCard Component
 * Wraps PokemonCard with 3D flip animation
 * Shows Pokeball on back, Pokemon details on front
 * Used for AI Pokemon that switch during battle
 */
const FlippableCard = ({ 
  pokemon, 
  isFlipped = false,
  isActive = false,
  isKnockedOut = false,
  onSelect,
  className = ''
}) => {
  return (
    <div className={`perspective-1000 ${className}`}>
      <motion.div
        className="relative preserve-3d"
        initial={{ rotateY: 0 }}
        animate={{ rotateY: isFlipped ? 180 : 0 }}
        transition={{ 
          duration: 0.6,
          ease: 'easeInOut'
        }}
        style={{
          transformStyle: 'preserve-3d'
        }}
      >
        {/* Front side - Pokemon details */}
        <div 
          className="backface-hidden"
          style={{
            backfaceVisibility: 'hidden',
            WebkitBackfaceVisibility: 'hidden'
          }}
        >
          <PokemonCard
            pokemon={pokemon}
            isActive={isActive}
            isFaceDown={false}
            isKnockedOut={isKnockedOut}
            onSelect={onSelect}
          />
        </div>

        {/* Back side - Pokeball */}
        <div 
          className="absolute top-0 left-0 backface-hidden"
          style={{
            backfaceVisibility: 'hidden',
            WebkitBackfaceVisibility: 'hidden',
            transform: 'rotateY(180deg)'
          }}
        >
          <motion.div
            className="w-40 h-56 rounded-lg border-4 border-gray-400 bg-gradient-to-br from-red-100 to-white shadow-lg flex items-center justify-center"
            whileHover={{ scale: 1.05, y: -5 }}
            transition={{ duration: 0.2 }}
          >
            <div className="text-center">
              <img 
                src="/assets/pokeball.png" 
                alt="Hidden Pokemon" 
                className="w-24 h-24 opacity-70 mx-auto"
              />
              <p className="text-gray-500 text-xs mt-2 font-semibold">???</p>
            </div>
          </motion.div>
        </div>
      </motion.div>
    </div>
  );
};

FlippableCard.propTypes = {
  pokemon: PropTypes.object,
  isFlipped: PropTypes.bool,
  isActive: PropTypes.bool,
  isKnockedOut: PropTypes.bool,
  onSelect: PropTypes.func,
  className: PropTypes.string
};

export default FlippableCard;

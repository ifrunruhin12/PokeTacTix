import { motion } from 'framer-motion';
import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import PokemonCard from './PokemonCard';

/**
 * AnimatedPokemonCard Component
 * Wraps PokemonCard with damage and knockout animations
 * - Shake animation when Pokemon takes damage
 * - Knockout animation (grayscale, fade, scale down)
 * - Red X overlay for knocked out Pokemon
 */
const AnimatedPokemonCard = ({ 
  pokemon, 
  isActive = false,
  isFaceDown = false,
  isKnockedOut = false,
  onDamage,
  onSelect,
  className = ''
}) => {
  const [isDamaged, setIsDamaged] = useState(false);
  const [previousHp, setPreviousHp] = useState(pokemon?.hp);

  // Detect damage and trigger shake animation
  useEffect(() => {
    if (pokemon?.hp !== undefined && previousHp !== undefined) {
      if (pokemon.hp < previousHp) {
        // Pokemon took damage
        setIsDamaged(true);
        if (onDamage) {
          onDamage(previousHp - pokemon.hp);
        }
        
        // Reset shake animation after it completes
        const timer = setTimeout(() => {
          setIsDamaged(false);
        }, 500);
        
        return () => clearTimeout(timer);
      }
    }
    setPreviousHp(pokemon?.hp);
  }, [pokemon?.hp, previousHp, onDamage]);

  // Shake animation variants
  const shakeVariants = {
    normal: { x: 0 },
    shake: {
      x: [-5, 5, -5, 5, -3, 3, 0],
      transition: {
        duration: 0.5,
        ease: 'easeInOut'
      }
    }
  };

  // Knockout animation variants
  const knockoutVariants = {
    normal: {
      scale: 1,
      opacity: 1,
      filter: 'grayscale(0%)',
      y: 0
    },
    knockout: {
      scale: 0.95,
      opacity: 0.6,
      filter: 'grayscale(100%)',
      y: 5,
      transition: {
        duration: 1,
        ease: 'easeOut'
      }
    }
  };

  return (
    <motion.div
      className={`relative ${className}`}
      variants={shakeVariants}
      animate={isDamaged ? 'shake' : 'normal'}
    >
      <motion.div
        variants={knockoutVariants}
        animate={isKnockedOut ? 'knockout' : 'normal'}
      >
        <PokemonCard
          pokemon={pokemon}
          isActive={isActive}
          isFaceDown={isFaceDown}
          isKnockedOut={isKnockedOut}
          onSelect={onSelect}
        />
      </motion.div>

      {/* Damage indicator */}
      {isDamaged && !isKnockedOut && (
        <motion.div
          className="absolute inset-0 pointer-events-none"
          initial={{ opacity: 0 }}
          animate={{ opacity: [0, 0.5, 0] }}
          transition={{ duration: 0.5 }}
        >
          <div className="absolute inset-0 bg-red-500 rounded-lg" />
        </motion.div>
      )}

      {/* Knockout overlay with red X */}
      {isKnockedOut && (
        <motion.div
          className="absolute inset-0 flex items-center justify-center pointer-events-none"
          initial={{ opacity: 0, scale: 0 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ 
            duration: 0.5,
            delay: 0.3,
            type: 'spring',
            stiffness: 200
          }}
        >
          <div className="relative">
            {/* Red X with glow effect */}
            <motion.div
              className="text-red-500 text-7xl font-bold drop-shadow-2xl"
              animate={{
                scale: [1, 1.1, 1],
                rotate: [0, 5, -5, 0]
              }}
              transition={{
                duration: 2,
                repeat: Infinity,
                ease: 'easeInOut'
              }}
            >
              ✕
            </motion.div>
            
            {/* Glow effect */}
            <div className="absolute inset-0 blur-xl bg-red-500 opacity-50 -z-10" />
          </div>
        </motion.div>
      )}
    </motion.div>
  );
};

AnimatedPokemonCard.propTypes = {
  pokemon: PropTypes.object,
  isActive: PropTypes.bool,
  isFaceDown: PropTypes.bool,
  isKnockedOut: PropTypes.bool,
  onDamage: PropTypes.func,
  onSelect: PropTypes.func,
  className: PropTypes.string
};

export default AnimatedPokemonCard;

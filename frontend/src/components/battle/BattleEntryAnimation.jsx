import { motion } from 'framer-motion';
import PropTypes from 'prop-types';

/**
 * BattleEntryAnimation Component
 * Displays dramatic entry animation when battle starts
 * Shows Pokemon sliding in from sides with spring animations
 */
const BattleEntryAnimation = ({ 
  playerPokemon, 
  aiPokemon, 
  onComplete,
  className = '' 
}) => {
  // Container animation
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        when: 'beforeChildren',
        staggerChildren: 0.3
      }
    },
    exit: {
      opacity: 0,
      transition: {
        when: 'afterChildren',
        staggerChildren: 0.1,
        staggerDirection: -1
      }
    }
  };

  // Player Pokemon animation (from left)
  const playerVariants = {
    hidden: { 
      x: -300, 
      opacity: 0,
      scale: 0.5,
      rotate: -20
    },
    visible: {
      x: 0,
      opacity: 1,
      scale: 1,
      rotate: 0,
      transition: {
        type: 'spring',
        stiffness: 100,
        damping: 15,
        duration: 0.8
      }
    },
    exit: {
      x: -300,
      opacity: 0,
      scale: 0.5,
      transition: {
        duration: 0.5
      }
    }
  };

  // AI Pokemon animation (from right)
  const aiVariants = {
    hidden: { 
      x: 300, 
      opacity: 0,
      scale: 0.5,
      rotate: 20
    },
    visible: {
      x: 0,
      opacity: 1,
      scale: 1,
      rotate: 0,
      transition: {
        type: 'spring',
        stiffness: 100,
        damping: 15,
        duration: 0.8
      }
    },
    exit: {
      x: 300,
      opacity: 0,
      scale: 0.5,
      transition: {
        duration: 0.5
      }
    }
  };

  // VS text animation
  const vsVariants = {
    hidden: { 
      scale: 0,
      opacity: 0,
      rotate: -180
    },
    visible: {
      scale: 1,
      opacity: 1,
      rotate: 0,
      transition: {
        type: 'spring',
        stiffness: 200,
        damping: 20,
        delay: 0.6
      }
    },
    exit: {
      scale: 0,
      opacity: 0,
      transition: {
        duration: 0.3
      }
    }
  };

  // Flash effect animation
  const flashVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        duration: 0.2,
        delay: 0.8
      }
    }
  };

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      onAnimationComplete={onComplete}
      className={`fixed inset-0 bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center z-50 ${className}`}
    >
      {/* Flash effect */}
      <motion.div
        variants={flashVariants}
        className="absolute inset-0 bg-white pointer-events-none"
      />

      <div className="relative w-full max-w-6xl px-4">
        <div className="grid grid-cols-3 gap-8 items-center">
          {/* Player Pokemon */}
          <motion.div variants={playerVariants} className="text-center">
            <div className="mb-4">
              <img
                src={playerPokemon?.sprite}
                alt={playerPokemon?.name}
                className="w-48 h-48 mx-auto object-contain drop-shadow-2xl"
              />
            </div>
            <motion.div
              animate={{
                scale: 1.1
              }}
              transition={{
                duration: 1,
                repeat: Infinity,
                repeatType: 'reverse',
                ease: 'easeInOut'
              }}
              className="text-3xl font-bold text-blue-400"
              style={{
                textShadow: '0 0 15px rgba(59, 130, 246, 0.8)'
              }}
            >
              {playerPokemon?.name || playerPokemon?.pokemon_name}
            </motion.div>
            <div className="text-gray-400 text-lg mt-2">Your Champion</div>
          </motion.div>

          {/* VS Text */}
          <motion.div variants={vsVariants} className="text-center">
            <motion.div
              animate={{
                scale: 1.1,
                rotate: 5
              }}
              transition={{
                duration: 1,
                repeat: Infinity,
                repeatType: 'reverse',
                ease: 'easeInOut'
              }}
              className="inline-block bg-gradient-to-r from-red-600 via-yellow-500 to-blue-600 text-white font-black text-7xl px-12 py-6 rounded-2xl shadow-2xl transform"
              style={{
                textShadow: '0 0 30px rgba(255, 255, 255, 0.8)'
              }}
            >
              VS
            </motion.div>
            <motion.div
              initial={{ scaleX: 0 }}
              animate={{ scaleX: 1 }}
              transition={{ delay: 1, duration: 0.5 }}
              className="h-1 bg-gradient-to-r from-blue-500 via-purple-500 to-red-500 mt-6 rounded-full"
            />
          </motion.div>

          {/* AI Pokemon */}
          <motion.div variants={aiVariants} className="text-center">
            <div className="mb-4">
              <img
                src={aiPokemon?.sprite}
                alt={aiPokemon?.name}
                className="w-48 h-48 mx-auto object-contain drop-shadow-2xl"
              />
            </div>
            <motion.div
              animate={{
                scale: 1.1
              }}
              transition={{
                duration: 1,
                repeat: Infinity,
                repeatType: 'reverse',
                ease: 'easeInOut'
              }}
              className="text-3xl font-bold text-red-400"
              style={{
                textShadow: '0 0 15px rgba(239, 68, 68, 0.8)'
              }}
            >
              {aiPokemon?.name || aiPokemon?.pokemon_name}
            </motion.div>
            <div className="text-gray-400 text-lg mt-2">Opponent</div>
          </motion.div>
        </div>

        {/* Battle start text */}
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.2, duration: 0.5 }}
          className="text-center mt-12"
        >
          <div className="text-white text-2xl font-bold">
            Battle Start!
          </div>
        </motion.div>
      </div>

      {/* Particle effects */}
      {[...Array(20)].map((_, i) => (
        <motion.div
          key={i}
          initial={{
            x: Math.random() * (typeof window !== 'undefined' ? window.innerWidth : 1000),
            y: Math.random() * (typeof window !== 'undefined' ? window.innerHeight : 800),
            scale: 0,
            opacity: 0
          }}
          animate={{
            scale: 1,
            opacity: 0.5,
            y: (Math.random() * (typeof window !== 'undefined' ? window.innerHeight : 800)) - 100
          }}
          transition={{
            duration: 2,
            delay: Math.random() * 1,
            repeat: Infinity,
            repeatDelay: Math.random() * 2,
            ease: 'linear'
          }}
          className="absolute w-2 h-2 bg-white rounded-full pointer-events-none"
        />
      ))}
    </motion.div>
  );
};

BattleEntryAnimation.propTypes = {
  playerPokemon: PropTypes.object,
  aiPokemon: PropTypes.object,
  onComplete: PropTypes.func,
  className: PropTypes.string
};

export default BattleEntryAnimation;

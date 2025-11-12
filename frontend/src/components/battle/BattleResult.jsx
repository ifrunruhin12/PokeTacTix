import PropTypes from 'prop-types';
import { motion } from 'framer-motion';
import PokemonCard from './PokemonCard';

/**
 * BattleResult Component
 * Displays victory/defeat/draw message with rewards
 * Shows coins earned, XP gained, level-up notifications
 * For 5v5 victories, shows AI Pokemon selection screen
 */
const BattleResult = ({ 
  result = 'victory', // 'victory', 'defeat', 'draw'
  rewards = {},
  aiPokemon = [],
  onSelectReward,
  onRematch,
  onNewBattle,
  onReturnToMenu,
  className = ''
}) => {
  const { coins_earned = 0, xp_gained = {}, level_ups = [] } = rewards;

  // Get result styling
  const getResultStyle = () => {
    switch (result) {
      case 'victory':
        return {
          bg: 'from-green-600 to-emerald-700',
          text: 'Victory!',
          icon: '🏆',
          message: 'You won the battle!'
        };
      case 'defeat':
        return {
          bg: 'from-red-600 to-rose-700',
          text: 'Defeat',
          icon: '💔',
          message: 'You lost the battle...'
        };
      case 'draw':
        return {
          bg: 'from-gray-600 to-gray-700',
          text: 'Draw',
          icon: '🤝',
          message: 'The battle ended in a draw'
        };
      default:
        return {
          bg: 'from-gray-600 to-gray-700',
          text: 'Battle Over',
          icon: '⚔️',
          message: 'The battle has ended'
        };
    }
  };

  const style = getResultStyle();
  const showRewardSelection = result === 'victory' && aiPokemon && aiPokemon.length > 0;

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ type: 'spring', stiffness: 200, damping: 20 }}
      className={`fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 ${className}`}
    >
      <div className="bg-gray-800 rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto border-2 border-gray-700">
        {/* Header */}
        <div className={`bg-gradient-to-r ${style.bg} p-6 text-center`}>
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1, rotate: [0, 10, -10, 0] }}
            transition={{ delay: 0.2, type: 'spring', stiffness: 200 }}
            className="text-7xl mb-4"
          >
            {style.icon}
          </motion.div>
          <h2 className="text-4xl font-bold text-white mb-2">{style.text}</h2>
          <p className="text-xl text-white/90">{style.message}</p>
        </div>

        {/* Rewards section */}
        <div className="p-6 space-y-6">
          {/* Coins earned */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="bg-gradient-to-r from-yellow-600 to-amber-600 rounded-lg p-4 text-center"
          >
            <div className="text-4xl mb-2">💰</div>
            <div className="text-white text-lg font-semibold">
              Coins Earned: <span className="text-2xl font-bold">{coins_earned}</span>
            </div>
          </motion.div>

          {/* XP gained */}
          {Object.keys(xp_gained).length > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
              className="bg-gray-700 rounded-lg p-4"
            >
              <h3 className="text-xl font-bold text-white mb-3 flex items-center gap-2">
                <span>⭐</span>
                <span>Experience Gained</span>
              </h3>
              <div className="grid grid-cols-2 gap-2">
                {Object.entries(xp_gained).map(([cardId, xp]) => (
                  <div key={cardId} className="bg-gray-800 rounded p-2 text-center">
                    <div className="text-purple-400 font-semibold">+{xp} XP</div>
                  </div>
                ))}
              </div>
            </motion.div>
          )}

          {/* Level ups */}
          {level_ups && level_ups.length > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.5 }}
              className="bg-gradient-to-r from-purple-600 to-pink-600 rounded-lg p-4"
            >
              <h3 className="text-xl font-bold text-white mb-3 flex items-center gap-2">
                <span>🎉</span>
                <span>Level Up!</span>
              </h3>
              <p className="text-white">
                {level_ups.length} Pokémon leveled up!
              </p>
            </motion.div>
          )}

          {/* AI Pokemon selection for 5v5 victories */}
          {showRewardSelection && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.6 }}
              className="bg-gray-700 rounded-lg p-4"
            >
              <h3 className="text-2xl font-bold text-white mb-4 text-center">
                🎁 Choose Your Reward Pokémon!
              </h3>
              <p className="text-gray-300 text-center mb-4">
                Select one Pokémon from your opponent's team to add to your collection
              </p>
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                {aiPokemon.map((pokemon, index) => (
                  <motion.div
                    key={index}
                    whileHover={{ scale: 1.05, y: -5 }}
                    whileTap={{ scale: 0.95 }}
                    onClick={() => onSelectReward && onSelectReward(index)}
                    className="cursor-pointer"
                  >
                    <PokemonCard
                      pokemon={pokemon}
                      isActive={false}
                      isFaceDown={false}
                      isKnockedOut={false}
                    />
                  </motion.div>
                ))}
              </div>
            </motion.div>
          )}

          {/* Action buttons */}
          {!showRewardSelection && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.7 }}
              className="flex flex-wrap gap-3 justify-center"
            >
              {onRematch && (
                <button
                  onClick={onRematch}
                  className="px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white font-bold rounded-lg shadow-lg transition-all"
                >
                  🔄 Rematch
                </button>
              )}
              {onNewBattle && (
                <button
                  onClick={onNewBattle}
                  className="px-6 py-3 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-500 hover:to-green-600 text-white font-bold rounded-lg shadow-lg transition-all"
                >
                  ⚔️ New Battle
                </button>
              )}
              {onReturnToMenu && (
                <button
                  onClick={onReturnToMenu}
                  className="px-6 py-3 bg-gradient-to-r from-gray-600 to-gray-700 hover:from-gray-500 hover:to-gray-600 text-white font-bold rounded-lg shadow-lg transition-all"
                >
                  🏠 Return to Menu
                </button>
              )}
            </motion.div>
          )}
        </div>
      </div>
    </motion.div>
  );
};

BattleResult.propTypes = {
  result: PropTypes.oneOf(['victory', 'defeat', 'draw']),
  rewards: PropTypes.shape({
    coins_earned: PropTypes.number,
    xp_gained: PropTypes.object,
    level_ups: PropTypes.arrayOf(PropTypes.number)
  }),
  aiPokemon: PropTypes.arrayOf(PropTypes.object),
  onSelectReward: PropTypes.func,
  onRematch: PropTypes.func,
  onNewBattle: PropTypes.func,
  onReturnToMenu: PropTypes.func,
  className: PropTypes.string
};

export default BattleResult;

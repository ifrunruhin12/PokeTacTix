import PropTypes from 'prop-types';
import { motion } from 'framer-motion';

/**
 * TurnIndicator Component
 * Displays current turn number, whose turn it is, and battle mode
 * Visual highlight for active player's turn
 */
const TurnIndicator = ({ 
  turnNumber = 1, 
  whoseTurn = 'player', 
  mode = '5v5',
  roundNumber = 1,
  className = '' 
}) => {
  const isPlayerTurn = whoseTurn === 'player';

  return (
    <div className={`flex items-center justify-between gap-4 ${className}`}>
      {/* Battle mode badge */}
      <motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        className="bg-gradient-to-r from-purple-600 to-blue-600 px-4 py-2 rounded-lg shadow-lg"
      >
        <div className="text-white font-bold text-sm">
          {mode === '5v5' ? '⚔️ 5v5 Battle' : '⚔️ 1v1 Battle'}
        </div>
      </motion.div>

      {/* Turn/Round info */}
      <div className="flex items-center gap-4">
        {mode === '5v5' && (
          <div className="bg-gray-800 px-4 py-2 rounded-lg border border-gray-700">
            <div className="text-gray-400 text-xs">Round</div>
            <div className="text-white font-bold text-lg">{roundNumber}</div>
          </div>
        )}

        <div className="bg-gray-800 px-4 py-2 rounded-lg border border-gray-700">
          <div className="text-gray-400 text-xs">Turn</div>
          <div className="text-white font-bold text-lg">{turnNumber}</div>
        </div>
      </div>

      {/* Active turn indicator */}
      <motion.div
        animate={{
          scale: [1, 1.05, 1],
          boxShadow: isPlayerTurn 
            ? ['0 0 0px rgba(59, 130, 246, 0.5)', '0 0 20px rgba(59, 130, 246, 0.8)', '0 0 0px rgba(59, 130, 246, 0.5)']
            : ['0 0 0px rgba(239, 68, 68, 0.5)', '0 0 20px rgba(239, 68, 68, 0.8)', '0 0 0px rgba(239, 68, 68, 0.5)']
        }}
        transition={{
          duration: 2,
          repeat: Infinity,
          ease: 'easeInOut'
        }}
        className={`px-6 py-3 rounded-lg font-bold text-white shadow-lg ${
          isPlayerTurn 
            ? 'bg-gradient-to-r from-blue-600 to-blue-700' 
            : 'bg-gradient-to-r from-red-600 to-red-700'
        }`}
      >
        <div className="flex items-center gap-2">
          <motion.span
            animate={{ rotate: [0, 360] }}
            transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
            className="text-xl"
          >
            {isPlayerTurn ? '👤' : '🤖'}
          </motion.span>
          <div>
            <div className="text-xs opacity-80">Current Turn</div>
            <div className="text-lg font-bold">
              {isPlayerTurn ? 'Your Turn' : 'AI Turn'}
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
};

TurnIndicator.propTypes = {
  turnNumber: PropTypes.number,
  whoseTurn: PropTypes.oneOf(['player', 'ai']),
  mode: PropTypes.oneOf(['1v1', '5v5']),
  roundNumber: PropTypes.number,
  className: PropTypes.string
};

export default TurnIndicator;

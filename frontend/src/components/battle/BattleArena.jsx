import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { motion } from 'framer-motion';
import AnimatedPokemonCard from './AnimatedPokemonCard';
import BattleControls from './BattleControls';
import BattleLog from './BattleLog';
import TurnIndicator from './TurnIndicator';
import MoveSelector from './MoveSelector';
import BattleResult from './BattleResult';

/**
 * BattleArena Component
 * Main battle interface for 5v5 battles
 * - Player deck (5 cards) on left, AI deck (5 cards) on right
 * - Active Pokemon cards with full details in center
 * - Inactive AI Pokemon as face-down Pokeball images
 * - Highlighted active Pokemon with glowing border
 * - Battle log section at bottom
 */
const BattleArena = ({ 
  battleState,
  onMove,
  onSwitchPokemon,
  onSelectReward,
  onRematch,
  onNewBattle,
  onReturnToMenu,
  loading = false,
  error = null,
  className = ''
}) => {
  const [showMoveSelector, setShowMoveSelector] = useState(false);
  const [battleLogs, setBattleLogs] = useState([]);

  // Extract battle state
  const {
    mode = '5v5',
    player_deck = [],
    ai_deck = [],
    player_active_idx = 0,
    ai_active_idx = 0,
    turn_number = 1,
    whose_turn = 'player',
    battle_over = false,
    winner = null,
    round_number = 1,
    log = []
  } = battleState || {};

  const playerActive = player_deck[player_active_idx];
  const aiActive = ai_deck[ai_active_idx];
  const isPlayerTurn = whose_turn === 'player';

  // Update battle logs when log changes
  useEffect(() => {
    if (log && log.length > 0) {
      setBattleLogs(log);
    }
  }, [log]);

  // Handle attack action
  const handleAttack = () => {
    setShowMoveSelector(true);
  };

  // Handle move selection
  const handleSelectMove = (moveIdx) => {
    setShowMoveSelector(false);
    if (onMove) {
      onMove('attack', moveIdx);
    }
  };

  // Handle other actions
  const handleDefend = () => {
    if (onMove) {
      onMove('defend');
    }
  };

  const handlePass = () => {
    if (onMove) {
      onMove('pass');
    }
  };

  const handleSacrifice = () => {
    if (onMove) {
      onMove('sacrifice');
    }
  };

  const handleSurrender = () => {
    if (window.confirm('Are you sure you want to surrender?')) {
      if (onMove) {
        onMove('surrender');
      }
    }
  };

  // Handle Pokemon switching
  const handleSwitchPokemon = (index) => {
    if (onSwitchPokemon && index !== player_active_idx) {
      const pokemon = player_deck[index];
      if (pokemon && !pokemon.is_knocked_out) {
        onSwitchPokemon(index);
      }
    }
  };

  // Entry animations
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 50 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        type: 'spring',
        stiffness: 100,
        damping: 15
      }
    }
  };

  const activeCardVariants = {
    hidden: { opacity: 0, scale: 0.8 },
    visible: {
      opacity: 1,
      scale: 1,
      transition: {
        type: 'spring',
        stiffness: 200,
        damping: 20,
        delay: 0.3
      }
    }
  };

  if (!battleState) {
    return (
      <div className={`min-h-screen bg-gray-900 flex items-center justify-center ${className}`}>
        <div className="text-white text-xl">Loading battle...</div>
      </div>
    );
  }

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className={`min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-4 ${className}`}
    >
      {/* Turn Indicator */}
      <TurnIndicator
        turnNumber={turn_number}
        whoseTurn={whose_turn}
        mode={mode}
        roundNumber={round_number}
        className="mb-6"
      />

      {/* Error Display */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="max-w-2xl mx-auto mb-6"
        >
          <div className="bg-red-900/50 border-2 border-red-500 text-red-200 px-6 py-4 rounded-lg shadow-lg">
            <div className="flex items-center gap-3">
              <span className="text-2xl">⚠️</span>
              <div>
                <div className="font-bold text-lg">Error</div>
                <div className="text-sm">{error}</div>
              </div>
            </div>
          </div>
        </motion.div>
      )}

      {/* Loading Indicator */}
      {loading && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="max-w-2xl mx-auto mb-6"
        >
          <div className="bg-blue-900/50 border-2 border-blue-500 text-blue-200 px-6 py-4 rounded-lg shadow-lg">
            <div className="flex items-center gap-3">
              <motion.span
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                className="text-2xl"
              >
                ⚔️
              </motion.span>
              <div className="font-semibold">Processing move...</div>
            </div>
          </div>
        </motion.div>
      )}

      {/* Main Battle Area */}
      <div className="max-w-7xl mx-auto">
        {/* Active Pokemon Battle Area (Center) */}
        <div className="flex flex-col items-center justify-center mb-8 space-y-6">
          {/* AI Active Pokemon */}
          <motion.div variants={activeCardVariants} className="text-center">
            <h3 className="text-lg font-bold text-red-400 mb-3">Opponent's Active Pokémon</h3>
            {aiActive ? (
              <AnimatedPokemonCard
                pokemon={aiActive}
                isActive={true}
                isFaceDown={false}
                isKnockedOut={aiActive.is_knocked_out}
                className="mx-auto"
              />
            ) : (
              <div className="text-gray-400">No active Pokémon</div>
            )}
          </motion.div>

          {/* VS Indicator */}
          <motion.div
            animate={{
              scale: [1, 1.1, 1],
              rotate: [0, 5, -5, 0]
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: 'easeInOut'
            }}
            className="text-center"
          >
            <div className="inline-block bg-gradient-to-r from-red-600 to-blue-600 text-white font-bold text-2xl px-6 py-2 rounded-full shadow-2xl">
              ⚔️ VS ⚔️
            </div>
          </motion.div>

          {/* Player Active Pokemon */}
          <motion.div variants={activeCardVariants} className="text-center">
            <h3 className="text-lg font-bold text-blue-400 mb-3">Your Active Pokémon</h3>
            {playerActive ? (
              <AnimatedPokemonCard
                pokemon={playerActive}
                isActive={true}
                isFaceDown={false}
                isKnockedOut={playerActive.is_knocked_out}
                className="mx-auto"
              />
            ) : (
              <div className="text-gray-400">No active Pokémon</div>
            )}
          </motion.div>
        </div>

        {/* Team Decks (Side by Side) */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-6">
          {/* Player Deck (Left) */}
          <motion.div variants={cardVariants} className="space-y-3">
            <h3 className="text-lg font-bold text-white text-center">
              👤 Your Team
            </h3>
            <div className="flex flex-wrap justify-center gap-2">
              {player_deck.map((pokemon, index) => (
                <motion.div
                  key={index}
                  whileHover={!pokemon.is_knocked_out && index !== player_active_idx ? { scale: 1.05 } : {}}
                  onClick={() => handleSwitchPokemon(index)}
                  className="cursor-pointer"
                >
                  <AnimatedPokemonCard
                    pokemon={pokemon}
                    isActive={index === player_active_idx}
                    isFaceDown={false}
                    isKnockedOut={pokemon.is_knocked_out}
                    className="w-32 h-44"
                  />
                </motion.div>
              ))}
            </div>
          </motion.div>

          {/* AI Deck (Right) */}
          <motion.div variants={cardVariants} className="space-y-3">
            <h3 className="text-lg font-bold text-white text-center">
              🤖 Opponent's Team
            </h3>
            <div className="flex flex-wrap justify-center gap-2">
              {ai_deck.map((pokemon, index) => (
                <motion.div key={index}>
                  <AnimatedPokemonCard
                    pokemon={pokemon}
                    isActive={index === ai_active_idx}
                    isFaceDown={index !== ai_active_idx && !pokemon.is_knocked_out}
                    isKnockedOut={pokemon.is_knocked_out}
                    className="w-32 h-44"
                  />
                </motion.div>
              ))}
            </div>
          </motion.div>
        </div>

        {/* Battle Log */}
        <motion.div variants={cardVariants}>
          <BattleLog logs={battleLogs} className="mb-6" />
        </motion.div>

        {/* Battle Controls - Always visible when battle is active */}
        {!battle_over && (
          <motion.div variants={cardVariants}>
            <BattleControls
              onAttack={handleAttack}
              onDefend={handleDefend}
              onPass={handlePass}
              onSacrifice={handleSacrifice}
              onSurrender={handleSurrender}
              disabled={!isPlayerTurn || loading}
              currentStamina={playerActive?.stamina || 0}
              maxHp={playerActive?.hp_max || 0}
              className="mb-6"
            />
            
            {/* Turn status message */}
            {!isPlayerTurn && (
              <div className="text-center mt-4">
                <div className="inline-block bg-gray-700/80 text-gray-300 px-6 py-3 rounded-lg">
                  <span className="text-lg">🤖 AI is thinking...</span>
                </div>
              </div>
            )}
          </motion.div>
        )}
      </div>

      {/* Move Selector Modal */}
      <MoveSelector
        isOpen={showMoveSelector}
        moves={playerActive?.moves || []}
        currentStamina={playerActive?.stamina || 0}
        onSelectMove={handleSelectMove}
        onCancel={() => setShowMoveSelector(false)}
      />

      {/* Battle Result Screen */}
      {battle_over && (
        <BattleResult
          result={winner === 'player' ? 'victory' : winner === 'ai' ? 'defeat' : 'draw'}
          rewards={battleState.rewards}
          aiPokemon={winner === 'player' && mode === '5v5' ? ai_deck : []}
          onSelectReward={onSelectReward}
          onRematch={onRematch}
          onNewBattle={onNewBattle}
          onReturnToMenu={onReturnToMenu}
        />
      )}
    </motion.div>
  );
};

BattleArena.propTypes = {
  battleState: PropTypes.shape({
    mode: PropTypes.string,
    player_deck: PropTypes.array,
    ai_deck: PropTypes.array,
    player_active_idx: PropTypes.number,
    ai_active_idx: PropTypes.number,
    turn_number: PropTypes.number,
    whose_turn: PropTypes.string,
    battle_over: PropTypes.bool,
    winner: PropTypes.string,
    round_number: PropTypes.number,
    log: PropTypes.array,
    rewards: PropTypes.object
  }),
  onMove: PropTypes.func,
  onSwitchPokemon: PropTypes.func,
  onSelectReward: PropTypes.func,
  onRematch: PropTypes.func,
  onNewBattle: PropTypes.func,
  onReturnToMenu: PropTypes.func,
  loading: PropTypes.bool,
  error: PropTypes.string,
  className: PropTypes.string
};

export default BattleArena;

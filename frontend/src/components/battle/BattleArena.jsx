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
    // Different confirmation messages based on battle mode
    let confirmMessage = '';
    if (mode === '1v1') {
      confirmMessage = 'Surrendering will end the battle and you will lose. Are you sure?';
    } else if (mode === '5v5') {
      confirmMessage = 'Surrendering will knock out your current Pokemon. Are you sure?';
    } else {
      confirmMessage = 'Are you sure you want to surrender?';
    }
    
    if (window.confirm(confirmMessage)) {
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
      className={`min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-4 flex flex-col ${className}`}
    >
      {/* Error Display */}
      {error && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-4"
        >
          <div className="bg-red-900/50 border-2 border-red-500 text-red-200 px-4 py-2 rounded-lg shadow-lg text-center">
            ⚠️ {error}
          </div>
        </motion.div>
      )}

      {/* Loading Indicator */}
      {loading && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="mb-4"
        >
          <div className="bg-blue-900/50 border-2 border-blue-500 text-blue-200 px-4 py-2 rounded-lg shadow-lg text-center">
            <motion.span
              animate={{ rotate: 360 }}
              transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
              className="inline-block mr-2"
            >
              ⚔️
            </motion.span>
            Processing move...
          </div>
        </motion.div>
      )}

      {/* Main Battle Container - Fixed Height Layout */}
      <div className="flex-1 flex flex-col max-w-7xl mx-auto w-full">
        {/* Top Section: Active Pokemon Cards */}
        <div className="grid grid-cols-2 gap-8 mb-6">
          {/* Player Active Pokemon (Left) */}
          <motion.div variants={activeCardVariants} className="flex flex-col items-center">
            <div className="text-sm font-semibold text-blue-400 mb-2">Player's Active Card</div>
            {playerActive ? (
              <div className="w-64 h-80">
                <AnimatedPokemonCard
                  pokemon={playerActive}
                  isActive={true}
                  isFaceDown={false}
                  isKnockedOut={playerActive.is_knocked_out}
                  className="w-full h-full"
                />
              </div>
            ) : (
              <div className="w-64 h-80 flex items-center justify-center text-gray-400 border-2 border-dashed border-gray-600 rounded-lg">
                No active Pokémon
              </div>
            )}
          </motion.div>

          {/* AI Active Pokemon (Right) */}
          <motion.div variants={activeCardVariants} className="flex flex-col items-center">
            <div className="text-sm font-semibold text-red-400 mb-2">AI's Active Card</div>
            {aiActive ? (
              <div className="w-64 h-80">
                <AnimatedPokemonCard
                  pokemon={aiActive}
                  isActive={true}
                  isFaceDown={false}
                  isKnockedOut={aiActive.is_knocked_out}
                  className="w-full h-full"
                />
              </div>
            ) : (
              <div className="w-64 h-80 flex items-center justify-center text-gray-400 border-2 border-dashed border-gray-600 rounded-lg">
                No active Pokémon
              </div>
            )}
          </motion.div>
        </div>

        {/* Middle Section: Turn Indicator, Battle Log, and Controls */}
        <div className="flex-1 flex flex-col items-center justify-center mb-6 space-y-4">
          {/* Turn Indicator - Hide when battle is over */}
          {!battle_over && (
            <TurnIndicator
              turnNumber={turn_number}
              whoseTurn={whose_turn}
              mode={mode}
              roundNumber={round_number}
            />
          )}

          {/* Battle Over Message */}
          {battle_over && (
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ type: 'spring', stiffness: 200, damping: 20 }}
              className="text-center"
            >
              <div className="bg-gradient-to-r from-purple-600 to-pink-600 text-white px-8 py-4 rounded-xl shadow-lg text-2xl font-bold">
                ⚔️ Battle Complete! ⚔️
              </div>
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 }}
                className="text-gray-400 mt-2 text-sm"
              >
                Calculating rewards...
              </motion.div>
            </motion.div>
          )}

          {/* Battle Log - Compact */}
          <motion.div variants={cardVariants} className="w-full max-w-2xl">
            <BattleLog logs={battleLogs} compact={true} />
          </motion.div>

          {/* Battle Controls - Only show when battle is not over */}
          {!battle_over && (
            <motion.div variants={cardVariants} className="w-full max-w-2xl">
              <BattleControls
                onAttack={handleAttack}
                onDefend={handleDefend}
                onPass={handlePass}
                onSacrifice={handleSacrifice}
                onSurrender={handleSurrender}
                disabled={!isPlayerTurn || loading}
                currentStamina={playerActive?.stamina || 0}
                maxHp={playerActive?.hp_max || 0}
              />
              
              {/* Turn status message */}
              {!isPlayerTurn && (
                <div className="text-center mt-2">
                  <div className="inline-block bg-gray-700/80 text-gray-300 px-4 py-2 rounded-lg text-sm">
                    🤖 AI is thinking...
                  </div>
                </div>
              )}
            </motion.div>
          )}
        </div>

        {/* Bottom Section: Team Decks */}
        <div className="grid grid-cols-2 gap-8">
          {/* Player Deck (Left) */}
          <motion.div variants={cardVariants}>
            <div className="text-sm font-semibold text-white mb-2 text-center">Player's Deck</div>
            <div className="flex justify-center gap-2">
              {player_deck.map((pokemon, index) => (
                <motion.div
                  key={index}
                  whileHover={!pokemon.is_knocked_out && index !== player_active_idx ? { scale: 1.05, y: -5 } : {}}
                  onClick={() => handleSwitchPokemon(index)}
                  className={`cursor-pointer ${index === player_active_idx ? 'opacity-50' : ''}`}
                >
                  <div className="w-24 h-32">
                    <AnimatedPokemonCard
                      pokemon={pokemon}
                      isActive={index === player_active_idx}
                      isFaceDown={false}
                      isKnockedOut={pokemon.is_knocked_out}
                      className="w-full h-full"
                      compact={true}
                    />
                  </div>
                </motion.div>
              ))}
            </div>
          </motion.div>

          {/* AI Deck (Right) */}
          <motion.div variants={cardVariants}>
            <div className="text-sm font-semibold text-white mb-2 text-center">AI's Deck</div>
            <div className="flex justify-center gap-2">
              {ai_deck.map((pokemon, index) => (
                <motion.div key={index}>
                  <div className="w-24 h-32">
                    <AnimatedPokemonCard
                      pokemon={pokemon}
                      isActive={index === ai_active_idx}
                      isFaceDown={index !== ai_active_idx && !pokemon.is_knocked_out}
                      isKnockedOut={pokemon.is_knocked_out}
                      className="w-full h-full"
                      compact={true}
                    />
                  </div>
                </motion.div>
              ))}
            </div>
          </motion.div>
        </div>
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
          aiPokemon={winner === 'player' && mode === '5v5' && !battleState.reward_claimed ? ai_deck : []}
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

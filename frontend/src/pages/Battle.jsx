import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { BattleArena, BattleEntryAnimation } from '../components/battle';
import { startBattle, submitMove, switchPokemon, selectReward } from '../services/battle.service';
import { motion, AnimatePresence } from 'framer-motion';

export default function Battle() {
  const navigate = useNavigate();
  const [battleState, setBattleState] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [battleMode, setBattleMode] = useState(null);
  const [showEntryAnimation, setShowEntryAnimation] = useState(false);

  // Transform backend response to frontend format
  const transformBattleState = (response) => {
    console.log('Transforming battle state:', response);
    if (!response || !response.state) {
      console.error('Invalid response:', response);
      return null;
    }
    
    const { state, session } = response;
    
    // Transform Pokemon data from backend format
    const transformPokemon = (pokemon) => {
      if (!pokemon) return null;
      return {
        name: pokemon.Name,
        pokemon_name: pokemon.Name,
        hp: pokemon.HP,
        hp_max: pokemon.HPMax,
        stamina: pokemon.Stamina,
        defense: pokemon.Defense,
        attack: pokemon.Attack,
        speed: pokemon.Speed,
        moves: pokemon.Moves || [],
        types: pokemon.Types || [],
        sprite: pokemon.Sprite,
        level: pokemon.Level || 1,
        xp: pokemon.XP || 0,
        is_legendary: pokemon.IsLegendary || false,
        is_mythical: pokemon.IsMythical || false
      };
    };
    
    const transformed = {
      id: session,
      mode: state.BattleMode || '1v1',
      player_deck: (state.Player?.Deck || []).map(transformPokemon),
      ai_deck: (state.AI?.Deck || []).map(transformPokemon),
      player_active_idx: state.PlayerActiveIdx || 0,
      ai_active_idx: state.AIActiveIdx || 0,
      turn_number: state.TurnNumber || 1,
      round_number: state.Round || 1,
      whose_turn: state.turn?.WhoseTurn || 'player',
      battle_over: state.BattleOver || false,
      winner: state.BattleOver ? (state.PlayerSurrendered ? 'ai' : 'player') : null,
      log: []
    };
    
    console.log('Transformed state:', transformed);
    return transformed;
  };

  // Start a new battle
  const handleStartBattle = async (mode) => {
    setLoading(true);
    setError(null);
    try {
      const response = await startBattle(mode);
      const transformedState = transformBattleState(response);
      setBattleState(transformedState);
      setBattleMode(mode);
      setShowEntryAnimation(true);
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Failed to start battle');
      console.error('Error starting battle:', err);
    } finally {
      setLoading(false);
    }
  };

  // Handle entry animation complete
  const handleEntryAnimationComplete = () => {
    console.log('Entry animation complete');
    setShowEntryAnimation(false);
  };
  
  // Auto-skip animation after 3 seconds as fallback
  useEffect(() => {
    if (showEntryAnimation) {
      const timer = setTimeout(() => {
        console.log('Auto-skipping animation');
        setShowEntryAnimation(false);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [showEntryAnimation]);

  // Handle move submission
  const handleMove = async (move, moveIdx = null) => {
    if (!battleState?.id) {
      console.error('No battle ID found:', battleState);
      return;
    }
    
    console.log('Submitting move:', { battleId: battleState.id, move, moveIdx });
    setLoading(true);
    try {
      const result = await submitMove(battleState.id, move, moveIdx);
      console.log('Move result:', result);
      const transformedState = transformBattleState({ state: result.state || result, session: battleState.id });
      setBattleState(transformedState);
    } catch (err) {
      setError(err.message || 'Failed to submit move');
      console.error('Error submitting move:', err);
    } finally {
      setLoading(false);
    }
  };

  // Handle Pokemon switching
  const handleSwitchPokemon = async (newIdx) => {
    if (!battleState?.id) return;
    
    setLoading(true);
    try {
      const result = await switchPokemon(battleState.id, newIdx);
      const transformedState = transformBattleState({ state: result.state || result, session: battleState.id });
      setBattleState(transformedState);
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Failed to switch Pokemon');
      console.error('Error switching Pokemon:', err);
    } finally {
      setLoading(false);
    }
  };

  // Handle reward selection (5v5 victory)
  const handleSelectReward = async (pokemonIdx) => {
    if (!battleState?.id) return;
    
    setLoading(true);
    try {
      await selectReward(battleState.id, pokemonIdx);
      // After selecting reward, return to menu
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Failed to select reward');
      console.error('Error selecting reward:', err);
    } finally {
      setLoading(false);
    }
  };

  // Handle rematch
  const handleRematch = () => {
    if (battleMode) {
      handleStartBattle(battleMode);
    }
  };

  // Handle new battle
  const handleNewBattle = () => {
    setBattleState(null);
    setBattleMode(null);
  };

  // Handle return to menu
  const handleReturnToMenu = () => {
    navigate('/dashboard');
  };

  // Battle mode selection screen
  if (!battleState && !loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center p-4">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          className="max-w-2xl w-full"
        >
          <div className="bg-gray-800 rounded-2xl shadow-2xl p-8 border-2 border-gray-700">
            <h1 className="text-4xl font-bold text-white text-center mb-4">
              ⚔️ Battle Arena
            </h1>
            <p className="text-gray-400 text-center mb-8">
              Choose your battle mode and test your skills!
            </p>

            {error && (
              <div className="bg-red-900/50 border border-red-500 text-red-200 px-4 py-3 rounded-lg mb-6">
                {error}
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* 1v1 Battle */}
              <motion.button
                whileHover={{ scale: 1.05, y: -5 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => handleStartBattle('1v1')}
                className="bg-gradient-to-br from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white rounded-xl p-6 shadow-lg transition-all"
              >
                <div className="text-5xl mb-4">⚔️</div>
                <h2 className="text-2xl font-bold mb-2">1v1 Battle</h2>
                <p className="text-blue-100 text-sm mb-4">
                  Quick battle with one Pokemon
                </p>
                <div className="text-yellow-300 font-semibold">
                  💰 Win: 50 coins | Loss: 10 coins
                </div>
              </motion.button>

              {/* 5v5 Battle */}
              <motion.button
                whileHover={{ scale: 1.05, y: -5 }}
                whileTap={{ scale: 0.95 }}
                onClick={() => handleStartBattle('5v5')}
                className="bg-gradient-to-br from-purple-600 to-purple-700 hover:from-purple-500 hover:to-purple-600 text-white rounded-xl p-6 shadow-lg transition-all"
              >
                <div className="text-5xl mb-4">🏆</div>
                <h2 className="text-2xl font-bold mb-2">5v5 Battle</h2>
                <p className="text-purple-100 text-sm mb-4">
                  Epic battle with your full team
                </p>
                <div className="text-yellow-300 font-semibold">
                  💰 Win: 150 coins | Loss: 25 coins
                </div>
                <div className="text-green-300 text-xs mt-2">
                  + Choose 1 opponent Pokemon on victory!
                </div>
              </motion.button>
            </div>

            <div className="mt-8 text-center">
              <button
                onClick={() => navigate('/dashboard')}
                className="text-gray-400 hover:text-white transition-colors"
              >
                ← Back to Dashboard
              </button>
            </div>
          </div>
        </motion.div>
      </div>
    );
  }

  // Loading screen
  if (loading && !battleState) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <motion.div
            animate={{ rotate: 360 }}
            transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
            className="text-6xl mb-4"
          >
            ⚔️
          </motion.div>
          <div className="text-white text-xl">Preparing battle...</div>
        </div>
      </div>
    );
  }

  // Battle screen
  return (
    <>
      <AnimatePresence>
        {showEntryAnimation && battleState && (
          <BattleEntryAnimation
            playerPokemon={battleState.player_deck?.[battleState.player_active_idx]}
            aiPokemon={battleState.ai_deck?.[battleState.ai_active_idx]}
            onComplete={handleEntryAnimationComplete}
          />
        )}
      </AnimatePresence>

      {!showEntryAnimation && (
        <BattleArena
          battleState={battleState}
          onMove={handleMove}
          onSwitchPokemon={handleSwitchPokemon}
          onSelectReward={handleSelectReward}
          onRematch={handleRematch}
          onNewBattle={handleNewBattle}
          onReturnToMenu={handleReturnToMenu}
        />
      )}
    </>
  );
}

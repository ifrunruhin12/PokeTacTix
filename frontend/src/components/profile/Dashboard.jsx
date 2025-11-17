import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { useAuth } from '../../contexts/AuthContext';
import { getPlayerStats } from '../../services/stats.service';
import { getUserCards } from '../../services/card.service';

/**
 * Dashboard Component
 * Displays user info and overall statistics
 * Requirements: 8.1, 8.5, 9.5
 */
export default function Dashboard() {
  const { user } = useAuth();
  const [stats, setStats] = useState(null);
  const [highestLevelPokemon, setHighestLevelPokemon] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Load stats and cards separately with individual error handling
      let statsData = null;
      let cardsData = null;

      try {
        statsData = await getPlayerStats();
      } catch (err) {
        console.error('Failed to load stats:', err);
        // Use default stats if API fails
        statsData = {
          total_battles_1v1: 0,
          wins_1v1: 0,
          losses_1v1: 0,
          draws_1v1: 0,
          total_battles_5v5: 0,
          wins_5v5: 0,
          losses_5v5: 0,
          draws_5v5: 0,
          total_coins_earned: 0,
          highest_level: 1,
        };
      }

      try {
        cardsData = await getUserCards();
      } catch (err) {
        console.error('Failed to load cards:', err);
        cardsData = { cards: [] };
      }

      setStats(statsData);

      // Find highest level Pokemon
      if (cardsData && cardsData.cards && cardsData.cards.length > 0) {
        const highest = cardsData.cards.reduce((max, card) => 
          card.level > max.level ? card : max
        );
        setHighestLevelPokemon(highest);
      }
    } catch (err) {
      console.error('Failed to load dashboard data:', err);
      setError('Failed to load dashboard data. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-900/20 border border-red-500 rounded-lg p-4 text-red-400">
        {error}
      </div>
    );
  }

  if (!stats) {
    return null;
  }

  // Calculate overall stats
  const totalBattles = stats.total_battles_1v1 + stats.total_battles_5v5;
  const totalWins = stats.wins_1v1 + stats.wins_5v5;
  const totalLosses = stats.losses_1v1 + stats.losses_5v5;
  const totalDraws = (stats.draws_1v1 || 0) + (stats.draws_5v5 || 0);
  
  // Win rate excludes draws: wins / (wins + losses)
  const winRate = (totalWins + totalLosses) > 0 
    ? ((totalWins / (totalWins + totalLosses)) * 100).toFixed(1) 
    : 0;

  // Calculate 1v1 win rate (excluding draws)
  const winRate1v1 = (stats.wins_1v1 + stats.losses_1v1) > 0 
    ? ((stats.wins_1v1 / (stats.wins_1v1 + stats.losses_1v1)) * 100).toFixed(1) 
    : 0;

  // Calculate 5v5 win rate (excluding draws)
  const winRate5v5 = (stats.wins_5v5 + stats.losses_5v5) > 0 
    ? ((stats.wins_5v5 / (stats.wins_5v5 + stats.losses_5v5)) * 100).toFixed(1) 
    : 0;

  return (
    <div className="space-y-6">
      {/* User Info Card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-gradient-to-br from-blue-900/40 to-purple-900/40 rounded-xl p-6 border border-blue-500/30"
      >
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold text-white mb-2">{user?.username}</h2>
            <p className="text-gray-400">
              Trainer since {new Date(user?.created_at).toLocaleDateString()}
            </p>
          </div>
          <div className="text-right">
            <div className="flex items-center gap-2 justify-end">
              <span className="text-4xl">💰</span>
              <span className="text-3xl font-bold text-yellow-400">{user?.coins}</span>
            </div>
            <p className="text-gray-400 text-sm mt-1">Coins</p>
          </div>
        </div>
      </motion.div>

      {/* Overall Statistics */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
      >
        <h3 className="text-2xl font-bold text-white mb-4">Overall Statistics</h3>
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
          <StatCard
            label="Total Battles"
            value={totalBattles}
            icon="⚔️"
            color="blue"
          />
          <StatCard
            label="Wins"
            value={totalWins}
            icon="🏆"
            color="green"
          />
          <StatCard
            label="Losses"
            value={totalLosses}
            icon="💔"
            color="red"
          />
          <StatCard
            label="Draws"
            value={totalDraws}
            icon="🤝"
            color="gray"
          />
          <StatCard
            label="Win Rate"
            value={`${winRate}%`}
            icon="📊"
            color="purple"
          />
        </div>
      </motion.div>

      {/* Battle Mode Statistics */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* 1v1 Stats */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
          className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
        >
          <h3 className="text-xl font-bold text-white mb-4 flex items-center gap-2">
            <span>⚡</span>
            1v1 Battles
          </h3>
          <div className="space-y-3">
            <StatRow label="Total Battles" value={stats.total_battles_1v1} />
            <StatRow label="Wins" value={stats.wins_1v1} color="green" />
            <StatRow label="Losses" value={stats.losses_1v1} color="red" />
            <StatRow label="Draws" value={stats.draws_1v1 || 0} color="gray" />
            <StatRow label="Win Rate" value={`${winRate1v1}%`} color="purple" />
          </div>
        </motion.div>

        {/* 5v5 Stats */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
          className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
        >
          <h3 className="text-xl font-bold text-white mb-4 flex items-center gap-2">
            <span>🔥</span>
            5v5 Battles
          </h3>
          <div className="space-y-3">
            <StatRow label="Total Battles" value={stats.total_battles_5v5} />
            <StatRow label="Wins" value={stats.wins_5v5} color="green" />
            <StatRow label="Losses" value={stats.losses_5v5} color="red" />
            <StatRow label="Draws" value={stats.draws_5v5 || 0} color="gray" />
            <StatRow label="Win Rate" value={`${winRate5v5}%`} color="purple" />
          </div>
        </motion.div>
      </div>

      {/* Additional Stats */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
      >
        <h3 className="text-2xl font-bold text-white mb-4">Achievements</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="flex items-center gap-4 p-4 bg-gray-900/50 rounded-lg">
            <span className="text-4xl">💰</span>
            <div>
              <p className="text-gray-400 text-sm">Total Coins Earned</p>
              <p className="text-2xl font-bold text-yellow-400">
                {stats.total_coins_earned}
              </p>
            </div>
          </div>
          {highestLevelPokemon && (
            <div className="flex items-center gap-4 p-4 bg-gray-900/50 rounded-lg">
              <img
                src={highestLevelPokemon.sprite}
                alt={highestLevelPokemon.pokemon_name}
                className="w-16 h-16"
              />
              <div>
                <p className="text-gray-400 text-sm">Highest Level Pokémon</p>
                <p className="text-xl font-bold text-white">
                  {highestLevelPokemon.pokemon_name}
                </p>
                <p className="text-lg text-blue-400">
                  Level {highestLevelPokemon.level}
                </p>
              </div>
            </div>
          )}
        </div>
      </motion.div>
    </div>
  );
}

/**
 * StatCard Component
 * Displays a single statistic in a card format
 */
function StatCard({ label, value, icon, color }) {
  const colorClasses = {
    blue: 'from-blue-900/40 to-blue-800/40 border-blue-500/30',
    green: 'from-green-900/40 to-green-800/40 border-green-500/30',
    red: 'from-red-900/40 to-red-800/40 border-red-500/30',
    gray: 'from-gray-900/40 to-gray-800/40 border-gray-500/30',
    purple: 'from-purple-900/40 to-purple-800/40 border-purple-500/30',
  };

  return (
    <div className={`bg-gradient-to-br ${colorClasses[color]} rounded-lg p-4 border`}>
      <div className="flex items-center gap-2 mb-2">
        <span className="text-2xl">{icon}</span>
      </div>
      <p className="text-3xl font-bold text-white">{value}</p>
      <p className="text-gray-400 text-sm mt-1">{label}</p>
    </div>
  );
}

/**
 * StatRow Component
 * Displays a statistic in a row format
 */
function StatRow({ label, value, color = 'white' }) {
  const colorClasses = {
    white: 'text-white',
    green: 'text-green-400',
    red: 'text-red-400',
    gray: 'text-gray-400',
    purple: 'text-purple-400',
  };

  return (
    <div className="flex justify-between items-center">
      <span className="text-gray-400">{label}</span>
      <span className={`text-xl font-bold ${colorClasses[color]}`}>{value}</span>
    </div>
  );
}

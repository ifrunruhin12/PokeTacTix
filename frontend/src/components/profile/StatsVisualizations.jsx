import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { getPlayerStats, getBattleHistory } from '../../services/stats.service';
import { getUserCards } from '../../services/card.service';

/**
 * StatsVisualizations Component
 * Displays visual charts and graphs for player statistics
 * Requirements: 8.1
 */
export default function StatsVisualizations() {
  const [stats, setStats] = useState(null);
  const [history, setHistory] = useState([]);
  const [cards, setCards] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadVisualizationData();
  }, []);

  const loadVisualizationData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Load each data source separately with fallbacks
      let statsData = null;
      let historyData = null;
      let cardsData = null;

      try {
        statsData = await getPlayerStats();
      } catch (err) {
        console.error('Failed to load stats:', err);
        statsData = {
          total_battles_1v1: 0,
          wins_1v1: 0,
          losses_1v1: 0,
          draws_1v1: 0,
          total_battles_5v5: 0,
          wins_5v5: 0,
          losses_5v5: 0,
          draws_5v5: 0,
        };
      }

      try {
        historyData = await getBattleHistory(50);
      } catch (err) {
        console.error('Failed to load history:', err);
        historyData = { history: [] };
      }

      try {
        cardsData = await getUserCards();
      } catch (err) {
        console.error('Failed to load cards:', err);
        cardsData = { cards: [] };
      }

      setStats(statsData);
      setHistory(historyData.history || []);
      setCards(cardsData.cards || []);
    } catch (err) {
      console.error('Failed to load visualization data:', err);
      setError('Failed to load statistics. Please try again.');
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

  return (
    <div className="space-y-6">
      {/* Win/Loss Pie Chart */}
      <WinLossPieChart stats={stats} />

      {/* Battle History Timeline */}
      <BattleTimeline history={history} />

      {/* Level Distribution */}
      <LevelDistribution cards={cards} />
    </div>
  );
}

/**
 * WinLossPieChart Component
 * Displays win/loss/draw ratio as a visual pie chart
 */
function WinLossPieChart({ stats }) {
  const totalWins = stats.wins_1v1 + stats.wins_5v5;
  const totalLosses = stats.losses_1v1 + stats.losses_5v5;
  const totalDraws = (stats.draws_1v1 || 0) + (stats.draws_5v5 || 0);
  const total = totalWins + totalLosses + totalDraws;

  if (total === 0) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
      >
        <h3 className="text-2xl font-bold text-white mb-4">Win/Loss/Draw Ratio</h3>
        <p className="text-gray-400 text-center py-8">
          No battles yet. Start battling to see your stats!
        </p>
      </motion.div>
    );
  }

  const winPercentage = (totalWins / total) * 100;
  const lossPercentage = (totalLosses / total) * 100;
  const drawPercentage = (totalDraws / total) * 100;

  // Calculate pie chart segments (using conic gradient)
  const winDegrees = (winPercentage / 100) * 360;
  const lossDegrees = winDegrees + ((lossPercentage / 100) * 360);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
    >
      <h3 className="text-2xl font-bold text-white mb-6">Win/Loss/Draw Ratio</h3>
      
      <div className="flex flex-col md:flex-row items-center gap-8">
        {/* Pie Chart */}
        <div className="relative">
          <div
            className="w-48 h-48 rounded-full"
            style={{
              background: `conic-gradient(
                #10b981 0deg ${winDegrees}deg,
                #ef4444 ${winDegrees}deg ${lossDegrees}deg,
                #6b7280 ${lossDegrees}deg 360deg
              )`,
            }}
          />
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-32 h-32 bg-gray-800 rounded-full flex items-center justify-center">
              <div className="text-center">
                <p className="text-3xl font-bold text-white">
                  {winPercentage.toFixed(0)}%
                </p>
                <p className="text-xs text-gray-400">Win Rate</p>
              </div>
            </div>
          </div>
        </div>

        {/* Legend */}
        <div className="flex-1 space-y-4">
          <div className="flex items-center justify-between p-4 bg-green-900/20 rounded-lg border border-green-500/30">
            <div className="flex items-center gap-3">
              <div className="w-4 h-4 bg-green-500 rounded-full"></div>
              <span className="text-white font-semibold">Wins</span>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold text-green-400">{totalWins}</p>
              <p className="text-sm text-gray-400">{winPercentage.toFixed(1)}%</p>
            </div>
          </div>

          <div className="flex items-center justify-between p-4 bg-red-900/20 rounded-lg border border-red-500/30">
            <div className="flex items-center gap-3">
              <div className="w-4 h-4 bg-red-500 rounded-full"></div>
              <span className="text-white font-semibold">Losses</span>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold text-red-400">{totalLosses}</p>
              <p className="text-sm text-gray-400">{lossPercentage.toFixed(1)}%</p>
            </div>
          </div>

          <div className="flex items-center justify-between p-4 bg-gray-900/20 rounded-lg border border-gray-500/30">
            <div className="flex items-center gap-3">
              <div className="w-4 h-4 bg-gray-500 rounded-full"></div>
              <span className="text-white font-semibold">Draws</span>
            </div>
            <div className="text-right">
              <p className="text-2xl font-bold text-gray-400">{totalDraws}</p>
              <p className="text-sm text-gray-400">{drawPercentage.toFixed(1)}%</p>
            </div>
          </div>

          <div className="flex items-center justify-between p-4 bg-gray-900/50 rounded-lg border border-gray-700">
            <span className="text-gray-400">Total Battles</span>
            <span className="text-xl font-bold text-white">{total}</span>
          </div>
        </div>
      </div>
    </motion.div>
  );
}

/**
 * BattleTimeline Component
 * Displays recent battle history as a timeline
 */
function BattleTimeline({ history }) {
  if (history.length === 0) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
      >
        <h3 className="text-2xl font-bold text-white mb-4">Battle Timeline</h3>
        <p className="text-gray-400 text-center py-8">
          No battle history available yet.
        </p>
      </motion.div>
    );
  }

  // Group battles by date
  const groupedByDate = history.reduce((acc, battle) => {
    const date = new Date(battle.created_at).toLocaleDateString();
    if (!acc[date]) {
      acc[date] = [];
    }
    acc[date].push(battle);
    return acc;
  }, {});

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
      className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
    >
      <h3 className="text-2xl font-bold text-white mb-6">Battle Timeline</h3>

      <div className="space-y-6">
        {Object.entries(groupedByDate).slice(0, 7).map(([date, battles], dateIndex) => (
          <div key={date} className="relative">
            {/* Date Header */}
            <div className="flex items-center gap-3 mb-3">
              <span className="text-sm font-semibold text-gray-400">{date}</span>
              <div className="flex-1 h-px bg-gray-700"></div>
            </div>

            {/* Battles for this date */}
            <div className="space-y-2 pl-4 border-l-2 border-gray-700">
              {battles.map((battle, battleIndex) => {
                // Helper functions for result styling
                const getResultIcon = (result) => {
                  switch(result) {
                    case 'win': return '🏆';
                    case 'loss': return '💔';
                    case 'draw': return '🤝';
                    default: return '⚔️';
                  }
                };

                const getResultColor = (result) => {
                  switch(result) {
                    case 'win': return 'bg-green-900/10 border-l-4 border-green-500';
                    case 'loss': return 'bg-red-900/10 border-l-4 border-red-500';
                    case 'draw': return 'bg-gray-900/10 border-l-4 border-gray-500';
                    default: return 'bg-gray-900/10 border-l-4 border-gray-700';
                  }
                };

                const getDotColor = (result) => {
                  switch(result) {
                    case 'win': return 'bg-green-500';
                    case 'loss': return 'bg-red-500';
                    case 'draw': return 'bg-gray-500';
                    default: return 'bg-gray-700';
                  }
                };

                return (
                  <motion.div
                    key={battle.id}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: (dateIndex * 0.1) + (battleIndex * 0.05) }}
                    className={`relative pl-6 py-2 rounded-lg ${getResultColor(battle.result)}`}
                  >
                    {/* Timeline dot */}
                    <div
                      className={`absolute left-0 top-1/2 -translate-y-1/2 -translate-x-[21px] w-3 h-3 rounded-full ${getDotColor(battle.result)}`}
                    />

                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <span className="text-lg">
                          {getResultIcon(battle.result)}
                        </span>
                        <div>
                          <p className="text-white font-semibold">
                            {battle.mode.toUpperCase()} Battle
                          </p>
                          <p className="text-xs text-gray-400">
                            {new Date(battle.created_at).toLocaleTimeString()}
                          </p>
                        </div>
                      </div>
                      <span className="text-yellow-400 font-semibold">
                        +{battle.coins_earned} 💰
                      </span>
                    </div>
                  </motion.div>
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </motion.div>
  );
}

/**
 * LevelDistribution Component
 * Shows distribution of Pokemon levels in collection
 */
function LevelDistribution({ cards }) {
  if (cards.length === 0) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
        className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
      >
        <h3 className="text-2xl font-bold text-white mb-4">Level Distribution</h3>
        <p className="text-gray-400 text-center py-8">
          No Pokémon in your collection yet.
        </p>
      </motion.div>
    );
  }

  // Group cards by level ranges
  const levelRanges = {
    '1-10': 0,
    '11-20': 0,
    '21-30': 0,
    '31-40': 0,
    '41-50': 0,
  };

  cards.forEach(card => {
    const level = card.level;
    if (level <= 10) levelRanges['1-10']++;
    else if (level <= 20) levelRanges['11-20']++;
    else if (level <= 30) levelRanges['21-30']++;
    else if (level <= 40) levelRanges['31-40']++;
    else levelRanges['41-50']++;
  });

  const maxCount = Math.max(...Object.values(levelRanges));
  const avgLevel = (cards.reduce((sum, card) => sum + card.level, 0) / cards.length).toFixed(1);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
      className="bg-gray-800/50 rounded-xl p-6 border border-gray-700"
    >
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-2xl font-bold text-white">Level Distribution</h3>
        <div className="text-right">
          <p className="text-sm text-gray-400">Average Level</p>
          <p className="text-2xl font-bold text-blue-400">{avgLevel}</p>
        </div>
      </div>

      <div className="space-y-4">
        {Object.entries(levelRanges).map(([range, count], index) => {
          const percentage = maxCount > 0 ? (count / maxCount) * 100 : 0;
          
          return (
            <motion.div
              key={range}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.3 + (index * 0.05) }}
            >
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-semibold text-gray-300">
                  Level {range}
                </span>
                <span className="text-sm text-gray-400">
                  {count} Pokémon
                </span>
              </div>
              <div className="relative w-full h-8 bg-gray-900/50 rounded-lg overflow-hidden">
                <motion.div
                  initial={{ width: 0 }}
                  animate={{ width: `${percentage}%` }}
                  transition={{ duration: 0.8, delay: 0.3 + (index * 0.05) }}
                  className="absolute top-0 left-0 h-full bg-gradient-to-r from-blue-500 to-purple-500"
                />
                <div className="absolute inset-0 flex items-center justify-center text-sm font-bold text-white">
                  {count > 0 && count}
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>

      {/* Collection Stats */}
      <div className="mt-6 pt-6 border-t border-gray-700 grid grid-cols-3 gap-4">
        <div className="text-center">
          <p className="text-2xl font-bold text-white">{cards.length}</p>
          <p className="text-xs text-gray-400">Total Pokémon</p>
        </div>
        <div className="text-center">
          <p className="text-2xl font-bold text-green-400">
            {Math.max(...cards.map(c => c.level))}
          </p>
          <p className="text-xs text-gray-400">Highest Level</p>
        </div>
        <div className="text-center">
          <p className="text-2xl font-bold text-blue-400">
            {Math.min(...cards.map(c => c.level))}
          </p>
          <p className="text-xs text-gray-400">Lowest Level</p>
        </div>
      </div>
    </motion.div>
  );
}

import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { getBattleHistory } from '../../services/stats.service';

/**
 * BattleHistory Component
 * Displays the last 20 battles in a table format
 * Requirements: 8.3
 */
export default function BattleHistory() {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [limit, setLimit] = useState(20);
  const [hasMore, setHasMore] = useState(false);

  useEffect(() => {
    loadBattleHistory();
  }, [limit]);

  const loadBattleHistory = async () => {
    try {
      setLoading(true);
      setError(null);

      const data = await getBattleHistory(limit);
      setHistory(data.history || []);
      
      // Check if there might be more battles
      setHasMore(data.count === limit);
    } catch (err) {
      console.error('Failed to load battle history:', err);
      // Don't show error for empty history, just show empty state
      setHistory([]);
      setError(null);
    } finally {
      setLoading(false);
    }
  };

  const loadMore = () => {
    setLimit(prev => Math.min(prev + 20, 100)); // Cap at 100
  };

  if (loading && history.length === 0) {
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

  if (history.length === 0) {
    return (
      <div className="bg-gray-800/50 rounded-xl p-8 border border-gray-700 text-center">
        <p className="text-gray-400 text-lg">No battle history yet.</p>
        <p className="text-gray-500 text-sm mt-2">
          Start battling to see your history here!
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-gray-800/50 rounded-xl border border-gray-700 overflow-hidden"
      >
        <div className="p-6 border-b border-gray-700">
          <h3 className="text-2xl font-bold text-white">Battle History</h3>
          <p className="text-gray-400 text-sm mt-1">
            Showing {history.length} most recent battles
          </p>
        </div>

        {/* Desktop Table View */}
        <div className="hidden md:block overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-900/50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                  Date
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                  Mode
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                  Result
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                  Coins Earned
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                  Duration
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {history.map((battle, index) => (
                <motion.tr
                  key={battle.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                  className="hover:bg-gray-700/30 transition-colors"
                >
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {formatDate(battle.created_at)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <BattleModeTag mode={battle.mode} />
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <ResultBadge result={battle.result} />
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="text-yellow-400 font-semibold">
                      +{battle.coins_earned} 💰
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400">
                    {formatDuration(battle.duration)}
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Mobile Card View */}
        <div className="md:hidden divide-y divide-gray-700">
          {history.map((battle, index) => (
            <motion.div
              key={battle.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              className="p-4 hover:bg-gray-700/30 transition-colors"
            >
              <div className="flex justify-between items-start mb-2">
                <div>
                  <BattleModeTag mode={battle.mode} />
                  <p className="text-xs text-gray-400 mt-1">
                    {formatDate(battle.created_at)}
                  </p>
                </div>
                <ResultBadge result={battle.result} />
              </div>
              <div className="flex justify-between items-center mt-3">
                <span className="text-yellow-400 font-semibold">
                  +{battle.coins_earned} 💰
                </span>
                <span className="text-sm text-gray-400">
                  {formatDuration(battle.duration)}
                </span>
              </div>
            </motion.div>
          ))}
        </div>

        {/* Load More Button */}
        {hasMore && limit < 100 && (
          <div className="p-4 border-t border-gray-700 text-center">
            <button
              onClick={loadMore}
              disabled={loading}
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 
                       text-white rounded-lg transition-colors font-semibold"
            >
              {loading ? 'Loading...' : 'Load More'}
            </button>
          </div>
        )}
      </motion.div>
    </div>
  );
}

/**
 * BattleModeTag Component
 * Displays the battle mode with appropriate styling
 */
function BattleModeTag({ mode }) {
  const modeConfig = {
    '1v1': {
      label: '1v1',
      icon: '⚡',
      bgColor: 'bg-blue-900/40',
      borderColor: 'border-blue-500/50',
      textColor: 'text-blue-400',
    },
    '5v5': {
      label: '5v5',
      icon: '🔥',
      bgColor: 'bg-purple-900/40',
      borderColor: 'border-purple-500/50',
      textColor: 'text-purple-400',
    },
  };

  const config = modeConfig[mode] || modeConfig['1v1'];

  return (
    <span
      className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold
                  ${config.bgColor} ${config.borderColor} ${config.textColor} border`}
    >
      <span>{config.icon}</span>
      <span>{config.label}</span>
    </span>
  );
}

/**
 * ResultBadge Component
 * Displays the battle result with appropriate styling
 */
function ResultBadge({ result }) {
  const resultConfig = {
    win: {
      label: 'Victory',
      icon: '🏆',
      bgColor: 'bg-green-900/40',
      borderColor: 'border-green-500/50',
      textColor: 'text-green-400',
    },
    loss: {
      label: 'Defeat',
      icon: '💔',
      bgColor: 'bg-red-900/40',
      borderColor: 'border-red-500/50',
      textColor: 'text-red-400',
    },
    draw: {
      label: 'Draw',
      icon: '🤝',
      bgColor: 'bg-gray-900/40',
      borderColor: 'border-gray-500/50',
      textColor: 'text-gray-400',
    },
  };

  const config = resultConfig[result] || resultConfig['draw'];

  return (
    <span
      className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold
                  ${config.bgColor} ${config.borderColor} ${config.textColor} border`}
    >
      <span>{config.icon}</span>
      <span>{config.label}</span>
    </span>
  );
}

/**
 * Format date to readable string
 */
function formatDate(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now - date;
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return 'Just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  });
}

/**
 * Format duration in seconds to readable string
 */
function formatDuration(seconds) {
  if (!seconds) return 'N/A';
  
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  
  if (mins === 0) return `${secs}s`;
  return `${mins}m ${secs}s`;
}

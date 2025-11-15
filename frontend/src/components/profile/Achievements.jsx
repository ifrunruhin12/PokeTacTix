import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { getAchievements, checkAchievements } from '../../services/stats.service';

/**
 * Achievements Component
 * Displays all achievements with locked/unlocked status
 * Requirements: 8.4
 */
export default function Achievements() {
  const [achievements, setAchievements] = useState([]);
  const [stats, setStats] = useState({ total: 0, unlocked: 0, locked: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [newlyUnlocked, setNewlyUnlocked] = useState([]);
  const [checking, setChecking] = useState(false);

  useEffect(() => {
    loadAchievements();
  }, []);

  const loadAchievements = async () => {
    try {
      setLoading(true);
      setError(null);

      const data = await getAchievements();
      setAchievements(data.achievements || []);
      setStats({
        total: data.total || 0,
        unlocked: data.unlocked || 0,
        locked: data.locked || 0,
      });
    } catch (err) {
      console.error('Failed to load achievements:', err);
      // Show empty achievements instead of error
      setAchievements([]);
      setStats({ total: 0, unlocked: 0, locked: 0 });
      setError(null);
    } finally {
      setLoading(false);
    }
  };

  const handleCheckAchievements = async () => {
    try {
      setChecking(true);
      const data = await checkAchievements();
      
      if (data.newly_unlocked && data.newly_unlocked.length > 0) {
        setNewlyUnlocked(data.newly_unlocked);
        // Reload achievements to get updated status
        await loadAchievements();
        
        // Clear newly unlocked after 5 seconds
        setTimeout(() => {
          setNewlyUnlocked([]);
        }, 5000);
      }
    } catch (err) {
      console.error('Failed to check achievements:', err);
    } finally {
      setChecking(false);
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

  const progressPercentage = stats.total > 0 
    ? ((stats.unlocked / stats.total) * 100).toFixed(0) 
    : 0;

  return (
    <div className="space-y-6">
      {/* Achievement Stats Header */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-gradient-to-br from-yellow-900/40 to-orange-900/40 rounded-xl p-6 border border-yellow-500/30"
      >
        <div className="flex items-center justify-between mb-4">
          <div>
            <h3 className="text-2xl font-bold text-white">Achievements</h3>
            <p className="text-gray-300 mt-1">
              {stats.unlocked} of {stats.total} unlocked
            </p>
          </div>
          <button
            onClick={handleCheckAchievements}
            disabled={checking}
            className="px-4 py-2 bg-yellow-600 hover:bg-yellow-700 disabled:bg-gray-600 
                     text-white rounded-lg transition-colors font-semibold flex items-center gap-2"
          >
            {checking ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                Checking...
              </>
            ) : (
              <>
                <span>🔍</span>
                Check Progress
              </>
            )}
          </button>
        </div>

        {/* Progress Bar */}
        <div className="relative w-full h-4 bg-gray-900/50 rounded-full overflow-hidden">
          <motion.div
            initial={{ width: 0 }}
            animate={{ width: `${progressPercentage}%` }}
            transition={{ duration: 1, ease: 'easeOut' }}
            className="absolute top-0 left-0 h-full bg-gradient-to-r from-yellow-500 to-orange-500"
          />
          <div className="absolute inset-0 flex items-center justify-center text-xs font-bold text-white">
            {progressPercentage}%
          </div>
        </div>
      </motion.div>

      {/* Newly Unlocked Achievements Notification */}
      <AnimatePresence>
        {newlyUnlocked.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="bg-gradient-to-r from-yellow-900/60 to-orange-900/60 rounded-xl p-4 border-2 border-yellow-500"
          >
            <div className="flex items-center gap-3">
              <span className="text-4xl">🎉</span>
              <div>
                <h4 className="text-lg font-bold text-white">
                  New Achievement{newlyUnlocked.length > 1 ? 's' : ''} Unlocked!
                </h4>
                <p className="text-yellow-200">
                  {newlyUnlocked.map(a => a.name).join(', ')}
                </p>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Achievements Grid */}
      {achievements.length === 0 ? (
        <div className="bg-gray-800/50 rounded-xl p-8 border border-gray-700 text-center">
          <p className="text-gray-400 text-lg">No achievements available yet.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {achievements.map((achievement, index) => (
            <AchievementCard
              key={achievement.id}
              achievement={achievement}
              index={index}
              isNewlyUnlocked={newlyUnlocked.some(a => a.id === achievement.id)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

/**
 * AchievementCard Component
 * Displays a single achievement with its status
 */
function AchievementCard({ achievement, index, isNewlyUnlocked }) {
  const isUnlocked = achievement.unlocked;

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ delay: index * 0.05 }}
      className={`relative rounded-xl p-6 border transition-all ${
        isUnlocked
          ? 'bg-gradient-to-br from-yellow-900/40 to-orange-900/40 border-yellow-500/50'
          : 'bg-gray-800/50 border-gray-700 opacity-60'
      } ${isNewlyUnlocked ? 'ring-2 ring-yellow-500 animate-pulse' : ''}`}
    >
      {/* Lock/Unlock Icon */}
      <div className="absolute top-4 right-4">
        {isUnlocked ? (
          <span className="text-2xl">🔓</span>
        ) : (
          <span className="text-2xl opacity-50">🔒</span>
        )}
      </div>

      {/* Achievement Icon */}
      <div className="mb-4">
        <span className={`text-5xl ${!isUnlocked && 'grayscale opacity-50'}`}>
          {achievement.icon || '🏆'}
        </span>
      </div>

      {/* Achievement Name */}
      <h4 className={`text-xl font-bold mb-2 ${
        isUnlocked ? 'text-white' : 'text-gray-500'
      }`}>
        {achievement.name}
      </h4>

      {/* Achievement Description */}
      <p className={`text-sm ${
        isUnlocked ? 'text-gray-300' : 'text-gray-600'
      }`}>
        {achievement.description}
      </p>

      {/* Progress Bar (if applicable) */}
      {!isUnlocked && achievement.requirement_value && (
        <div className="mt-4">
          <div className="flex justify-between text-xs text-gray-500 mb-1">
            <span>Progress</span>
            <span>0 / {achievement.requirement_value}</span>
          </div>
          <div className="w-full h-2 bg-gray-900/50 rounded-full overflow-hidden">
            <div className="h-full bg-gray-600 w-0" />
          </div>
        </div>
      )}

      {/* Unlocked Date */}
      {isUnlocked && achievement.unlocked_at && (
        <p className="text-xs text-yellow-400 mt-3">
          Unlocked {formatDate(achievement.unlocked_at)}
        </p>
      )}

      {/* Newly Unlocked Badge */}
      {isNewlyUnlocked && (
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          className="absolute -top-2 -right-2 bg-yellow-500 text-gray-900 text-xs font-bold 
                     px-2 py-1 rounded-full"
        >
          NEW!
        </motion.div>
      )}
    </motion.div>
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

  if (diffMins < 1) return 'just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  });
}

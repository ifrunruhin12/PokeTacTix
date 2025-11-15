import { useState } from 'react';
import { motion } from 'framer-motion';
import Dashboard from '../components/profile/Dashboard';
import BattleHistory from '../components/profile/BattleHistory';
import Achievements from '../components/profile/Achievements';
import StatsVisualizations from '../components/profile/StatsVisualizations';

/**
 * Profile Page
 * Main profile page with tabs for different sections
 * Requirements: 8.1, 8.3, 8.4, 8.5, 9.5
 */
export default function Profile() {
  const [activeTab, setActiveTab] = useState('dashboard');

  const tabs = [
    { id: 'dashboard', label: 'Dashboard', icon: '📊' },
    { id: 'history', label: 'Battle History', icon: '⚔️' },
    { id: 'achievements', label: 'Achievements', icon: '🏆' },
    { id: 'stats', label: 'Statistics', icon: '📈' },
  ];

  return (
    <div className="min-h-screen bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-4xl font-bold text-white mb-2">Profile</h1>
          <p className="text-gray-400">View your stats, achievements, and battle history</p>
        </motion.div>

        {/* Tabs */}
        <div className="mb-8">
          <div className="flex flex-wrap gap-2 border-b border-gray-700">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-6 py-3 font-semibold transition-all ${
                  activeTab === tab.id
                    ? 'text-white border-b-2 border-blue-500 bg-blue-900/20'
                    : 'text-gray-400 hover:text-white hover:bg-gray-800/50'
                }`}
              >
                <span>{tab.icon}</span>
                <span>{tab.label}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Tab Content */}
        <motion.div
          key={activeTab}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          {activeTab === 'dashboard' && <Dashboard />}
          {activeTab === 'history' && <BattleHistory />}
          {activeTab === 'achievements' && <Achievements />}
          {activeTab === 'stats' && <StatsVisualizations />}
        </motion.div>
      </div>
    </div>
  );
}

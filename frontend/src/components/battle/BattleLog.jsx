import { useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { motion, AnimatePresence } from 'framer-motion';

/**
 * BattleLog Component
 * Displays turn-by-turn battle events with auto-scroll
 * Color codes different event types (damage, healing, status changes)
 */
const BattleLog = ({ logs = [], className = '' }) => {
  const logEndRef = useRef(null);
  const containerRef = useRef(null);

  // Auto-scroll to latest log entry
  useEffect(() => {
    if (logEndRef.current) {
      logEndRef.current.scrollIntoView({ behavior: 'smooth', block: 'end' });
    }
  }, [logs]);

  // Determine log entry color based on content
  const getLogColor = (log) => {
    const logLower = log.toLowerCase();
    
    if (logLower.includes('victory') || logLower.includes('won') || logLower.includes('wins')) {
      return 'text-green-400';
    }
    if (logLower.includes('defeat') || logLower.includes('lost') || logLower.includes('loses')) {
      return 'text-red-400';
    }
    if (logLower.includes('damage') || logLower.includes('dealt') || logLower.includes('hit')) {
      return 'text-orange-400';
    }
    if (logLower.includes('heal') || logLower.includes('recover')) {
      return 'text-green-300';
    }
    if (logLower.includes('defend') || logLower.includes('block')) {
      return 'text-blue-400';
    }
    if (logLower.includes('knocked out') || logLower.includes('fainted')) {
      return 'text-red-500 font-bold';
    }
    if (logLower.includes('switch') || logLower.includes('sent out')) {
      return 'text-purple-400';
    }
    if (logLower.includes('turn') || logLower.includes('round')) {
      return 'text-yellow-400 font-semibold';
    }
    
    return 'text-gray-300';
  };

  // Get icon based on log content
  const getLogIcon = (log) => {
    const logLower = log.toLowerCase();
    
    if (logLower.includes('victory') || logLower.includes('won')) return '🏆';
    if (logLower.includes('defeat') || logLower.includes('lost')) return '💔';
    if (logLower.includes('damage') || logLower.includes('dealt')) return '⚔️';
    if (logLower.includes('heal')) return '💚';
    if (logLower.includes('defend')) return '🛡️';
    if (logLower.includes('knocked out') || logLower.includes('fainted')) return '💀';
    if (logLower.includes('switch')) return '🔄';
    if (logLower.includes('turn') || logLower.includes('round')) return '▶️';
    
    return '•';
  };

  if (!logs || logs.length === 0) {
    return (
      <div className={`bg-gray-800 rounded-lg p-4 ${className}`}>
        <h3 className="text-lg font-bold text-gray-400 mb-2">Battle Log</h3>
        <p className="text-gray-500 text-sm italic">Waiting for battle to start...</p>
      </div>
    );
  }

  return (
    <div className={`bg-gray-800 rounded-lg p-4 ${className}`}>
      <h3 className="text-lg font-bold text-white mb-3 flex items-center gap-2">
        <span>📜</span>
        <span>Battle Log</span>
      </h3>
      
      <div 
        ref={containerRef}
        className="space-y-1 max-h-48 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-gray-700"
      >
        <AnimatePresence initial={false}>
          {logs.map((log, index) => (
            <motion.div
              key={index}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
              transition={{ duration: 0.3 }}
              className={`text-sm ${getLogColor(log)} flex items-start gap-2 py-1`}
            >
              <span className="flex-shrink-0 mt-0.5">{getLogIcon(log)}</span>
              <span className="flex-1">{log}</span>
            </motion.div>
          ))}
        </AnimatePresence>
        <div ref={logEndRef} />
      </div>
    </div>
  );
};

BattleLog.propTypes = {
  logs: PropTypes.arrayOf(PropTypes.string),
  className: PropTypes.string
};

export default BattleLog;

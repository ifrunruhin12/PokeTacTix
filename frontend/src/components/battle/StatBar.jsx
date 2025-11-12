import { motion } from 'framer-motion';
import PropTypes from 'prop-types';

/**
 * StatBar Component
 * Displays an animated stat bar with color coding
 * Colors: green (HP), blue (stamina), red (attack), yellow (defense), purple (speed)
 */
const StatBar = ({ 
  label, 
  current, 
  max, 
  type = 'default',
  showValues = true,
  height = 'h-3',
  className = ''
}) => {
  // Calculate percentage
  const percentage = Math.max(0, Math.min(100, (current / max) * 100));

  // Get color based on stat type
  const getColor = () => {
    switch (type.toLowerCase()) {
      case 'hp':
        // HP changes color based on percentage
        if (percentage < 30) return '#ef4444'; // red
        if (percentage < 60) return '#f59e0b'; // orange
        return '#10b981'; // green
      case 'stamina':
        return '#3b82f6'; // blue
      case 'attack':
        return '#ef4444'; // red
      case 'defense':
        return '#eab308'; // yellow
      case 'speed':
        return '#a855f7'; // purple
      case 'xp':
        return '#8b5cf6'; // violet
      default:
        return '#6b7280'; // gray
    }
  };

  // Get background color (lighter version)
  const getBgColor = () => {
    switch (type.toLowerCase()) {
      case 'hp':
        return 'bg-green-100';
      case 'stamina':
        return 'bg-blue-100';
      case 'attack':
        return 'bg-red-100';
      case 'defense':
        return 'bg-yellow-100';
      case 'speed':
        return 'bg-purple-100';
      case 'xp':
        return 'bg-violet-100';
      default:
        return 'bg-gray-200';
    }
  };

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      {/* Label */}
      {label && (
        <span className="text-sm font-medium text-gray-700 min-w-[60px]">
          {label}:
        </span>
      )}

      {/* Bar container */}
      <div className={`flex-1 ${getBgColor()} rounded-full ${height} overflow-hidden relative`}>
        {/* Animated fill */}
        <motion.div
          className="h-full rounded-full"
          initial={{ width: 0 }}
          animate={{ 
            width: `${percentage}%`,
            backgroundColor: getColor()
          }}
          transition={{ 
            duration: 0.5,
            ease: 'easeOut'
          }}
        />

        {/* Optional percentage text overlay */}
        {showValues && percentage > 20 && (
          <span className="absolute inset-0 flex items-center justify-center text-xs font-semibold text-white drop-shadow-md">
            {Math.round(percentage)}%
          </span>
        )}
      </div>

      {/* Current/Max values */}
      {showValues && (
        <span className="text-sm font-medium text-gray-700 min-w-[50px] text-right">
          {current}/{max}
        </span>
      )}
    </div>
  );
};

StatBar.propTypes = {
  label: PropTypes.string,
  current: PropTypes.number.isRequired,
  max: PropTypes.number.isRequired,
  type: PropTypes.oneOf(['hp', 'stamina', 'attack', 'defense', 'speed', 'xp', 'default']),
  showValues: PropTypes.bool,
  height: PropTypes.string,
  className: PropTypes.string
};

export default StatBar;

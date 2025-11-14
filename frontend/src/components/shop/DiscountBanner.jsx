import { motion } from 'framer-motion';
import PropTypes from 'prop-types';

/**
 * DiscountBanner Component
 * Displays a banner when a discount event is active
 */
export default function DiscountBanner({ discountPercent, refreshTime }) {
  if (!discountPercent || discountPercent === 0) return null;

  // Calculate time remaining until refresh
  const getTimeRemaining = () => {
    const now = new Date();
    const refresh = new Date(refreshTime);
    const diff = refresh - now;

    if (diff <= 0) return 'Ending soon';

    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

    if (hours > 0) {
      return `${hours}h ${minutes}m remaining`;
    }
    return `${minutes}m remaining`;
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-gradient-to-r from-purple-600 via-pink-600 to-red-600 rounded-lg p-4 mb-6 shadow-lg"
    >
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-3">
          <motion.div
            animate={{ rotate: [0, 10, -10, 10, 0] }}
            transition={{ duration: 0.5, repeat: Infinity, repeatDelay: 2 }}
            className="text-4xl"
          >
            🎉
          </motion.div>
          <div>
            <h3 className="text-xl font-bold text-white mb-1">
              Special Discount Event Active!
            </h3>
            <p className="text-white/90 text-sm">
              <span className="font-bold">{discountPercent}%</span> off Legendary & Mythical Pokemon
            </p>
          </div>
        </div>
        <div className="bg-white/20 backdrop-blur-sm rounded-lg px-4 py-2">
          <p className="text-white font-semibold text-sm">
            ⏰ {getTimeRemaining()}
          </p>
        </div>
      </div>
    </motion.div>
  );
}

DiscountBanner.propTypes = {
  discountPercent: PropTypes.number,
  refreshTime: PropTypes.string,
};

DiscountBanner.defaultProps = {
  discountPercent: 0,
  refreshTime: null,
};

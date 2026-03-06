import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import PropTypes from 'prop-types';
import tokenService from '../../services/token.service';

/**
 * TokenPurchaseSection Component
 * Displays game token purchase interface in the shop
 */
export default function TokenPurchaseSection({ 
  userCoins, 
  onPurchase, 
  isProcessing,
  error,
  successMessage,
}) {
  const [tokenData, setTokenData] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadTokenData();
  }, []);

  const loadTokenData = async () => {
    try {
      setLoading(true);
      const data = await tokenService.getTokenBalance();
      setTokenData(data);
    } catch (err) {
      console.error('Failed to load token data:', err);
    } finally {
      setLoading(false);
    }
  };

  // Refresh token data after successful purchase
  useEffect(() => {
    if (successMessage) {
      loadTokenData();
    }
  }, [successMessage]);

  if (loading) {
    return (
      <div className="bg-gray-800 rounded-lg border-2 border-gray-700 p-6 mb-8">
        <div className="flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
      </div>
    );
  }

  const tokensPurchasedToday = tokenData?.tokens_purchased_today || 0;
  const dailyLimit = tokenData?.daily_purchase_limit || 10;
  const remainingPurchases = dailyLimit - tokensPurchasedToday;
  const isLimitReached = remainingPurchases <= 0;
  
  const pricePerToken = 100;
  const totalCost = quantity * pricePerToken;
  const canAfford = userCoins >= totalCost;
  const maxQuantity = Math.min(remainingPurchases, 10);

  const handleQuantityChange = (newQuantity) => {
    const validQuantity = Math.max(1, Math.min(newQuantity, maxQuantity));
    setQuantity(validQuantity);
  };

  const handlePurchase = () => {
    if (canAfford && !isLimitReached && !isProcessing) {
      onPurchase(quantity);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-lg border-2 border-blue-500 shadow-lg shadow-blue-500/30 p-6 mb-8"
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-white flex items-center gap-2">
            🎫 Game Tokens
          </h2>
          <p className="text-gray-400 text-sm mt-1">
            Purchase additional tokens to play more battles
          </p>
        </div>
        <div className="text-right">
          <div className="text-sm text-gray-400">Current Tokens</div>
          <div className="text-2xl font-bold text-blue-400">
            {tokenData?.game_tokens || 0}
            {tokenData?.daily_token_limit && (
              <span className="text-sm text-gray-400">/{tokenData.daily_token_limit}</span>
            )}
          </div>
        </div>
      </div>

      {/* Success Message */}
      {successMessage && (
        <div className="bg-green-900/50 border border-green-500 rounded-lg p-3 mb-4">
          <p className="text-green-200 text-sm text-center font-semibold">✓ {successMessage}</p>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="bg-red-900/50 border border-red-500 rounded-lg p-3 mb-4">
          <p className="text-red-200 text-sm">{error}</p>
        </div>
      )}

      {/* Limit Reached Warning */}
      {isLimitReached && (
        <div className="bg-yellow-900/50 border border-yellow-500 rounded-lg p-4 mb-4">
          <p className="text-yellow-200 text-center font-semibold">
            ⚠️ Daily purchase limit reached ({tokensPurchasedToday}/{dailyLimit})
          </p>
          <p className="text-yellow-300 text-sm text-center mt-1">
            Resets at {tokenData?.reset_time || '00:00 UTC'} ({tokenData?.time_until_reset || 'calculating...'})
          </p>
        </div>
      )}

      {/* Purchase Interface */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Left Column - Info */}
        <div className="space-y-4">
          <div className="bg-gray-900 rounded-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-gray-400">Price per token:</span>
              <span className="text-yellow-400 font-bold">100 coins 🪙</span>
            </div>
            <div className="flex justify-between items-center mb-2">
              <span className="text-gray-400">Daily purchase limit:</span>
              <span className="text-white font-bold">
                {tokensPurchasedToday}/{dailyLimit} purchased today
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-400">Resets at:</span>
              <span className="text-blue-400 font-bold">
                {tokenData?.reset_time || '00:00 UTC'}
              </span>
            </div>
          </div>

          <div className="bg-blue-900/30 border border-blue-500/50 rounded-lg p-4">
            <p className="text-blue-200 text-sm mb-2">
              💡 <span className="font-semibold">Tip:</span> Each battle costs tokens:
            </p>
            <ul className="text-blue-200 text-sm space-y-1 ml-4">
              <li>• 1v1 Battle: {tokenData?.token_cost_1v1 || 1} token{(tokenData?.token_cost_1v1 || 1) !== 1 ? 's' : ''}</li>
              <li>• 5v5 Battle: {tokenData?.token_cost_5v5 || 2} token{(tokenData?.token_cost_5v5 || 2) !== 1 ? 's' : ''}</li>
            </ul>
            <p className="text-blue-200 text-sm mt-2">
              Purchase tokens to play more battles after using your daily free tokens!
            </p>
          </div>
        </div>

        {/* Right Column - Purchase Controls */}
        <div className="space-y-4">
          {/* Quantity Selector */}
          <div className="bg-gray-900 rounded-lg p-4">
            <label className="block text-gray-400 text-sm mb-2">
              Quantity (max {maxQuantity} remaining today)
            </label>
            <div className="flex items-center gap-3">
              <button
                onClick={() => handleQuantityChange(quantity - 1)}
                disabled={quantity <= 1 || isLimitReached || isProcessing}
                className="bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:text-gray-600 text-white font-bold w-10 h-10 rounded-lg transition-colors disabled:cursor-not-allowed"
              >
                -
              </button>
              <input
                type="number"
                min="1"
                max={maxQuantity}
                value={quantity}
                onChange={(e) => handleQuantityChange(parseInt(e.target.value) || 1)}
                disabled={isLimitReached || isProcessing}
                className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-center text-white font-bold text-xl disabled:opacity-50 disabled:cursor-not-allowed"
              />
              <button
                onClick={() => handleQuantityChange(quantity + 1)}
                disabled={quantity >= maxQuantity || isLimitReached || isProcessing}
                className="bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:text-gray-600 text-white font-bold w-10 h-10 rounded-lg transition-colors disabled:cursor-not-allowed"
              >
                +
              </button>
            </div>
            <div className="flex gap-2 mt-2">
              {[1, 5, maxQuantity].filter((val, idx, arr) => arr.indexOf(val) === idx && val > 0).map((preset) => (
                <button
                  key={preset}
                  onClick={() => handleQuantityChange(preset)}
                  disabled={isLimitReached || isProcessing}
                  className="flex-1 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:text-gray-600 text-white text-sm py-1 rounded transition-colors disabled:cursor-not-allowed"
                >
                  {preset}
                </button>
              ))}
            </div>
          </div>

          {/* Cost Calculator */}
          <div className="bg-gray-900 rounded-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-gray-400">Total cost:</span>
              <span className="text-yellow-400 font-bold text-xl">
                {totalCost} 🪙
              </span>
            </div>
            <div className="flex justify-between items-center mb-2">
              <span className="text-gray-400">Your coins:</span>
              <span className="text-white font-bold">{userCoins} 🪙</span>
            </div>
            <div className="border-t border-gray-700 pt-2 mt-2">
              <div className="flex justify-between items-center">
                <span className="text-gray-400">After purchase:</span>
                <span className={`font-bold ${canAfford ? 'text-green-400' : 'text-red-400'}`}>
                  {canAfford ? `${userCoins - totalCost} 🪙` : 'Insufficient Coins'}
                </span>
              </div>
            </div>
          </div>

          {/* Purchase Button */}
          <button
            onClick={handlePurchase}
            disabled={!canAfford || isLimitReached || isProcessing}
            className={`w-full py-3 rounded-lg font-bold transition-colors ${
              canAfford && !isLimitReached && !isProcessing
                ? 'bg-blue-600 hover:bg-blue-700 text-white'
                : 'bg-gray-700 text-gray-500 cursor-not-allowed'
            }`}
          >
            {isProcessing ? (
              <span className="flex items-center justify-center">
                <svg className="animate-spin h-5 w-5 mr-2" viewBox="0 0 24 24">
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                    fill="none"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                Processing...
              </span>
            ) : isLimitReached ? (
              'Daily Limit Reached'
            ) : !canAfford ? (
              'Insufficient Coins'
            ) : (
              `Purchase ${quantity} Token${quantity > 1 ? 's' : ''}`
            )}
          </button>
        </div>
      </div>
    </motion.div>
  );
}

TokenPurchaseSection.propTypes = {
  userCoins: PropTypes.number.isRequired,
  onPurchase: PropTypes.func.isRequired,
  isProcessing: PropTypes.bool,
  error: PropTypes.string,
  successMessage: PropTypes.string,
};

TokenPurchaseSection.defaultProps = {
  isProcessing: false,
  error: null,
  successMessage: null,
};

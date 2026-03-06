import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useState, useEffect } from 'react';
import tokenService from '../../services/token.service';

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [tokenData, setTokenData] = useState(null);
  const [showTooltip, setShowTooltip] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  // Fetch token balance when user is authenticated
  useEffect(() => {
    if (isAuthenticated) {
      const fetchTokenBalance = async () => {
        try {
          const data = await tokenService.getTokenBalance();
          setTokenData(data);
        } catch (error) {
          console.error('Failed to fetch token balance:', error);
        }
      };

      fetchTokenBalance();
      
      // Refresh token data every minute to keep time until reset accurate
      const interval = setInterval(fetchTokenBalance, 60000);
      
      return () => clearInterval(interval);
    }
  }, [isAuthenticated]);

  // Format token display based on count
  const formatTokenDisplay = () => {
    if (!tokenData) return '🎫 -';
    
    const tokens = tokenData.game_tokens || 0;
    const dailyLimit = tokenData.daily_token_limit || 5; // Use API value, fallback to 5
    
    if (tokens > dailyLimit) {
      return `🎫 ${tokens}`;
    }
    return `🎫 ${tokens}/${dailyLimit}`;
  };

  return (
    <nav className="bg-gray-800 border-b border-gray-700">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          <Link to="/" className="text-2xl font-bold text-white">
            PokeTacTix
          </Link>

          <div className="flex items-center gap-6">
            {isAuthenticated ? (
              <>
                <Link to="/dashboard" className="text-gray-300 hover:text-white transition-colors">
                  Dashboard
                </Link>
                <Link to="/battle" className="text-gray-300 hover:text-white transition-colors">
                  Battle
                </Link>
                <Link to="/shop" className="text-gray-300 hover:text-white transition-colors">
                  Shop
                </Link>
                <Link to="/deck" className="text-gray-300 hover:text-white transition-colors">
                  Deck
                </Link>
                <Link to="/profile" className="text-gray-300 hover:text-white transition-colors">
                  Profile
                </Link>
                
                <div className="flex items-center gap-4 ml-4 pl-4 border-l border-gray-700">
                  <div 
                    className="relative"
                    onMouseEnter={() => setShowTooltip(true)}
                    onMouseLeave={() => setShowTooltip(false)}
                  >
                    <span className="text-blue-400 font-semibold cursor-help">
                      {formatTokenDisplay()}
                    </span>
                    
                    {showTooltip && tokenData && (
                      <div className="absolute top-full mt-2 left-1/2 transform -translate-x-1/2 bg-gray-900 text-white text-sm rounded-lg px-3 py-2 whitespace-nowrap z-50 shadow-lg border border-gray-700">
                        <div className="text-center">
                          <div>Resets at {tokenData.reset_time}</div>
                          <div className="text-gray-400">({tokenData.time_until_reset})</div>
                        </div>
                        <div className="absolute -top-1 left-1/2 transform -translate-x-1/2 w-2 h-2 bg-gray-900 border-l border-t border-gray-700 rotate-45"></div>
                      </div>
                    )}
                  </div>
                  
                  <span className="text-yellow-400 font-semibold">
                    {user?.coins || 0} 🪙
                  </span>
                  <span className="text-gray-300">{user?.username}</span>
                  <button
                    onClick={handleLogout}
                    className="text-red-400 hover:text-red-300 transition-colors"
                  >
                    Logout
                  </button>
                </div>
              </>
            ) : (
              <Link to="/auth" className="btn-primary">
                Login / Register
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}

import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
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

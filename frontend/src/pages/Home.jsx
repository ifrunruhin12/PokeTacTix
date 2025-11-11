import { Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export default function Home() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-b from-gray-900 via-blue-900 to-gray-900">
      <div className="text-center space-y-8 px-4">
        <h1 className="text-6xl font-bold text-white mb-4">
          PokeTacTix
        </h1>
        <p className="text-xl text-gray-300 max-w-2xl mx-auto">
          Battle, collect, and level up your Pokemon cards in this strategic turn-based game
        </p>
        
        <div className="flex gap-4 justify-center mt-8">
          {isAuthenticated ? (
            <>
              <Link to="/dashboard" className="btn-primary">
                Go to Dashboard
              </Link>
              <Link to="/battle" className="btn-secondary">
                Start Battle
              </Link>
            </>
          ) : (
            <>
              <Link to="/auth" className="btn-primary">
                Get Started
              </Link>
              <a href="#features" className="btn-secondary">
                Learn More
              </a>
            </>
          )}
        </div>

        <div id="features" className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-4xl mx-auto">
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl font-bold mb-2">5v5 Battles</h3>
            <p className="text-gray-400">
              Engage in strategic 5v5 battles with your Pokemon deck
            </p>
          </div>
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl font-bold mb-2">Level Up</h3>
            <p className="text-gray-400">
              Earn XP and level up your Pokemon to increase their stats
            </p>
          </div>
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl font-bold mb-2">Collect Cards</h3>
            <p className="text-gray-400">
              Purchase new Pokemon from the shop and build your collection
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

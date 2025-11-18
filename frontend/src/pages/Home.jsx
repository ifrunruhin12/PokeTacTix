import { Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export default function Home() {
  const { isAuthenticated } = useAuth();

  return (
    <div 
      className="min-h-screen flex flex-col items-center justify-center relative"
      style={{
        backgroundImage: 'url(/assets/wallpaper.jpg)',
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat'
      }}
    >
      {/* Dark overlay for better text readability */}
      <div className="absolute inset-0 bg-black/50"></div>
      
      <div className="text-center space-y-8 px-4 relative z-10">
        <div className="flex items-center justify-center gap-4 mb-4">
          <img 
            src="/assets/pokeball.png" 
            alt="Pokeball" 
            className="w-16 h-16 animate-spin-slow"
            style={{ animationDuration: '3s' }}
          />
          <h1 className="text-6xl font-bold text-white">
            PokeTacTix
          </h1>
          <img 
            src="/assets/pokeball.png" 
            alt="Pokeball" 
            className="w-16 h-16 animate-spin-slow"
            style={{ animationDuration: '3s', animationDirection: 'reverse' }}
          />
        </div>
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
          <div className="bg-gray-800/80 backdrop-blur-sm p-6 rounded-lg border border-gray-700">
            <h3 className="text-xl font-bold mb-2 text-white">5v5 Battles</h3>
            <p className="text-gray-300">
              Engage in strategic 5v5 battles with your Pokemon deck
            </p>
          </div>
          <div className="bg-gray-800/80 backdrop-blur-sm p-6 rounded-lg border border-gray-700">
            <h3 className="text-xl font-bold mb-2 text-white">Level Up</h3>
            <p className="text-gray-300">
              Earn XP and level up your Pokemon to increase their stats
            </p>
          </div>
          <div className="bg-gray-800/80 backdrop-blur-sm p-6 rounded-lg border border-gray-700">
            <h3 className="text-xl font-bold mb-2 text-white">Collect Cards</h3>
            <p className="text-gray-300">
              Purchase new Pokemon from the shop and build your collection
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

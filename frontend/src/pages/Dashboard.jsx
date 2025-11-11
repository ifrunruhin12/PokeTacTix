import { useAuth } from '../hooks/useAuth';
import { Link } from 'react-router-dom';

export default function Dashboard() {
  const { user } = useAuth();

  return (
    <div className="min-h-screen bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4">
        <h1 className="text-4xl font-bold mb-8">Welcome, {user?.username}!</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl font-bold mb-2">Coins</h3>
            <p className="text-3xl text-yellow-400">{user?.coins || 0}</p>
          </div>
          
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl font-bold mb-2">Total Battles</h3>
            <p className="text-3xl text-blue-400">0</p>
          </div>
          
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl font-bold mb-2">Win Rate</h3>
            <p className="text-3xl text-green-400">0%</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Link to="/battle" className="bg-gradient-to-r from-red-600 to-red-700 p-8 rounded-lg hover:from-red-700 hover:to-red-800 transition-all">
            <h2 className="text-2xl font-bold mb-2">Battle Arena</h2>
            <p className="text-gray-200">Start a 1v1 or 5v5 battle</p>
          </Link>

          <Link to="/shop" className="bg-gradient-to-r from-blue-600 to-blue-700 p-8 rounded-lg hover:from-blue-700 hover:to-blue-800 transition-all">
            <h2 className="text-2xl font-bold mb-2">Shop</h2>
            <p className="text-gray-200">Purchase new Pokemon cards</p>
          </Link>

          <Link to="/deck" className="bg-gradient-to-r from-green-600 to-green-700 p-8 rounded-lg hover:from-green-700 hover:to-green-800 transition-all">
            <h2 className="text-2xl font-bold mb-2">Deck Manager</h2>
            <p className="text-gray-200">Customize your battle deck</p>
          </Link>

          <Link to="/profile" className="bg-gradient-to-r from-purple-600 to-purple-700 p-8 rounded-lg hover:from-purple-700 hover:to-purple-800 transition-all">
            <h2 className="text-2xl font-bold mb-2">Profile</h2>
            <p className="text-gray-200">View stats and achievements</p>
          </Link>
        </div>
      </div>
    </div>
  );
}

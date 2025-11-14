import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import shopService from '../services/shop.service';
import ShopGrid from '../components/shop/ShopGrid';
import ShopFilters from '../components/shop/ShopFilters';
import PurchaseModal from '../components/shop/PurchaseModal';
import DiscountBanner from '../components/shop/DiscountBanner';
import api from '../services/api';

export default function Shop() {
  const { user, updateUser } = useAuth();
  const [inventory, setInventory] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [ownedPokemon, setOwnedPokemon] = useState([]);
  
  // Filter and sort state
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedRarity, setSelectedRarity] = useState('all');
  const [sortBy, setSortBy] = useState('rarity');
  
  // Purchase modal state
  const [selectedItem, setSelectedItem] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [purchaseError, setPurchaseError] = useState(null);
  
  // Success message state
  const [successMessage, setSuccessMessage] = useState(null);

  // Store original prices for discount display
  const [originalPrices, setOriginalPrices] = useState({});

  // Load shop inventory
  useEffect(() => {
    loadInventory();
    loadOwnedPokemon();
  }, []);

  const loadInventory = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await shopService.getInventory();
      setInventory(data);
      
      // Calculate original prices if discount is active
      if (data.discount_active && data.discount_percent > 0) {
        const prices = {};
        data.items.forEach(item => {
          if (item.is_legendary || item.is_mythical) {
            // Reverse calculate original price
            const discountMultiplier = item.is_legendary ? 0.6 : 0.7;
            prices[item.pokemon_name] = Math.round(item.price / discountMultiplier);
          }
        });
        setOriginalPrices(prices);
      }
    } catch (err) {
      setError(err.message || 'Failed to load shop inventory');
    } finally {
      setLoading(false);
    }
  };

  const loadOwnedPokemon = async () => {
    try {
      const response = await api.get('/api/cards');
      const cards = response.data.cards || response.data || [];
      const owned = cards.map(card => card.pokemon_name.toLowerCase());
      setOwnedPokemon(owned);
    } catch (err) {
      console.error('Failed to load owned Pokemon:', err);
    }
  };

  const handlePurchaseClick = (item) => {
    setSelectedItem(item);
    setPurchaseError(null);
    setIsModalOpen(true);
  };

  const handleConfirmPurchase = async () => {
    if (!selectedItem) return;

    try {
      setIsProcessing(true);
      setPurchaseError(null);
      
      const result = await shopService.purchaseCard(selectedItem.pokemon_name);
      
      // Update user coins
      updateUser({ coins: result.remaining_coins });
      
      // Add to owned Pokemon
      setOwnedPokemon(prev => [...prev, selectedItem.pokemon_name.toLowerCase()]);
      
      // Close modal
      setIsModalOpen(false);
      setSelectedItem(null);
      
      // Show success message
      setSuccessMessage(`Successfully purchased ${selectedItem.pokemon_name}!`);
      setTimeout(() => setSuccessMessage(null), 5000);
      
    } catch (err) {
      setPurchaseError(err.message || 'Failed to complete purchase');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleCloseModal = () => {
    if (!isProcessing) {
      setIsModalOpen(false);
      setSelectedItem(null);
      setPurchaseError(null);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 py-8">
        <div className="max-w-7xl mx-auto px-4">
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
              <p className="text-gray-400">Loading shop inventory...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-900 py-8">
        <div className="max-w-7xl mx-auto px-4">
          <div className="bg-red-900/50 border border-red-500 rounded-lg p-6 text-center">
            <p className="text-red-200 mb-4">{error}</p>
            <button
              onClick={loadInventory}
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg transition-colors"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Pokemon Shop</h1>
          <div className="flex items-center justify-between flex-wrap gap-4">
            <p className="text-gray-400">
              Browse and purchase Pokemon cards to expand your collection
            </p>
            <div className="bg-gray-800 rounded-lg px-4 py-2 border border-gray-700">
              <span className="text-gray-400 mr-2">Your Coins:</span>
              <span className="text-yellow-400 font-bold text-xl">{user?.coins || 0} 🪙</span>
            </div>
          </div>
        </div>

        {/* Success Message */}
        {successMessage && (
          <div className="bg-green-900/50 border border-green-500 rounded-lg p-4 mb-6">
            <p className="text-green-200 text-center font-semibold">✓ {successMessage}</p>
          </div>
        )}

        {/* Discount Banner */}
        {inventory?.discount_active && (
          <DiscountBanner
            discountPercent={inventory.discount_percent}
            refreshTime={inventory.refresh_time}
          />
        )}

        {/* Filters */}
        <ShopFilters
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          selectedRarity={selectedRarity}
          onRarityChange={setSelectedRarity}
          sortBy={sortBy}
          onSortChange={setSortBy}
        />

        {/* Shop Grid */}
        {inventory?.items && (
          <ShopGrid
            items={inventory.items}
            onPurchase={handlePurchaseClick}
            userCoins={user?.coins || 0}
            ownedPokemon={ownedPokemon}
            searchQuery={searchQuery}
            selectedRarity={selectedRarity}
            sortBy={sortBy}
            originalPrices={originalPrices}
          />
        )}

        {/* Purchase Modal */}
        <PurchaseModal
          isOpen={isModalOpen}
          onClose={handleCloseModal}
          item={selectedItem}
          userCoins={user?.coins || 0}
          onConfirm={handleConfirmPurchase}
          isProcessing={isProcessing}
          error={purchaseError}
        />
      </div>
    </div>
  );
}

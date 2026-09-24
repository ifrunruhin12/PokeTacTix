import { useState, useEffect, useRef } from 'react';
import { useAuth } from '../contexts/AuthContext';
import shopService from '../services/shop.service';
import itemService from '../services/item.service';
import ShopGrid from '../components/shop/ShopGrid';
import ShopFilters from '../components/shop/ShopFilters';
import PurchaseModal from '../components/shop/PurchaseModal';
import DiscountBanner from '../components/shop/DiscountBanner';
import TokenPurchaseSection from '../components/shop/TokenPurchaseSection';
import ItemsGrid from '../components/shop/ItemsGrid';
import api from '../services/api';

// Shop categories
const CATEGORIES = [
  { id: 'pokemon', label: '👤 Pokémon' },
  { id: 'tokens', label: '🎟️ Game Tokens' },
  { id: 'items', label: '💎 Items' },
];

export default function Shop() {
  const { user, updateUser } = useAuth();
  const [inventory, setInventory] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [ownedPokemon, setOwnedPokemon] = useState([]);
  const [itemInventory, setItemInventory] = useState([]);
  const [itemLoadingError, setItemLoadingError] = useState(null);
  const [requiresItemReload, setRequiresItemReload] = useState(false);
  const [activeCategory, setActiveCategory] = useState('pokemon');
  
  // Filter and sort state
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedRarity, setSelectedRarity] = useState('all');
  const [sortBy, setSortBy] = useState('rarity');
  
  // Purchase modal state
  const [selectedItem, setSelectedItem] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [purchaseError, setPurchaseError] = useState(null);
  
  // Item purchase state
  const [itemPurchaseError, setItemPurchaseError] = useState(null);
  const itemPurchaseKeys = useRef(new Map());
  const [usingItem, setUsingItem] = useState(null);
  const [activeBoosts, setActiveBoosts] = useState([]);
  
  // Token purchase state
  const [isTokenProcessing, setIsTokenProcessing] = useState(false);
  const [tokenPurchaseError, setTokenPurchaseError] = useState(null);
  const [tokenSuccessMessage, setTokenSuccessMessage] = useState(null);
  
  // Success message state
  const [successMessage, setSuccessMessage] = useState(null);

  // Store original prices for discount display
  const [originalPrices, setOriginalPrices] = useState({});

  // Load shop inventory
  useEffect(() => {
    loadInventory();
    loadOwnedPokemon();
    loadItemInventory();
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

  const loadItemInventory = async () => {
    try {
      const { inventory, activeBoosts } = await itemService.getInventory();
      setItemInventory(inventory);
      setActiveBoosts(activeBoosts);
      setItemLoadingError(null);
      setRequiresItemReload(false);
    } catch (err) {
      console.error('Failed to load item inventory:', err);
      setItemLoadingError('Failed to load item inventory. Reload inventory before using or purchasing items.');
    }
  };

  const handleUseItem = async (item) => {
    try {
      setUsingItem(item.id);
      setItemPurchaseError(null);

      const result = await itemService.useItem(item.id);

      setSuccessMessage(result.message || `${item.name} activated!`);
      setTimeout(() => setSuccessMessage(null), 5000);

      // Reload inventory (quantity dropped) and boosts (new active boost)
      await loadItemInventory();
    } catch (err) {
      setItemPurchaseError(err.message || 'Failed to use item');
    } finally {
      setUsingItem(null);
    }
  };

  const handleItemPurchase = async (item) => {
    if (isProcessing || itemLoadingError || requiresItemReload) return;
    const key = itemPurchaseKeys.current.get(item.id) || crypto.randomUUID();
    itemPurchaseKeys.current.set(item.id, key);
    try {
      setIsProcessing(true);
      setItemPurchaseError(null);

      const result = await shopService.purchaseItem(item.id, 1, key);
      itemPurchaseKeys.current.delete(item.id);

      // Update user coins and local inventory view
      updateUser({ coins: result.remaining_coins });
      setItemInventory(prev => {
        const existing = prev.find(entry => entry.id === item.id);
        if (existing) {
          return prev.map(entry =>
            entry.id === item.id
              ? { ...entry, quantity: result.new_quantity }
              : entry
          );
        }
        return [...prev, { ...item, quantity: result.new_quantity }];
      });

      setSuccessMessage(`Successfully purchased ${item.name}! You now own ${result.new_quantity}.`);
      setTimeout(() => setSuccessMessage(null), 5000);

    } catch (err) {
      if (err.code === 'ECONNABORTED'
        || err.message?.includes('timeout')
        || err.message === 'No response from server') {
        setItemPurchaseError('Purchase result unknown. Reload inventory before retrying; the purchase may already have completed.');
        setRequiresItemReload(true);
        console.error('Item purchase timed out (result uncertain):', err);
        return;
      }
      itemPurchaseKeys.current.delete(item.id);
      setItemPurchaseError(err.message || 'Failed to purchase item');
    } finally {
      setIsProcessing(false);
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
      // Timeouts / unreachable servers land here without a status. The
      // purchase may still complete server-side (coins deducted, card
      // created), so warn instead of showing a plain failure — a blind
      // retry can double-charge.
      if (err.code === 'ECONNABORTED'
        || err.message?.includes('timeout')
        || err.message === 'No response from server') {
        setPurchaseError('Purchase timed out. It may still have completed — reload the page before trying again to avoid buying twice.');
        console.error('Purchase timed out (result uncertain):', err);
        return;
      }

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

  const handleTokenPurchase = async (quantity) => {
    try {
      setIsTokenProcessing(true);
      setTokenPurchaseError(null);
      setTokenSuccessMessage(null);
      
      const result = await shopService.purchaseTokens(quantity);
      
      // Update user coins
      updateUser({ coins: result.remaining_coins });
      
      // Show success message
      setTokenSuccessMessage(
        `Successfully purchased ${result.tokens_added} token${result.tokens_added > 1 ? 's' : ''} for ${result.coins_spent} coins!`
      );
      setTimeout(() => setTokenSuccessMessage(null), 5000);
      
    } catch (err) {
      setTokenPurchaseError(err.response?.data?.error || err.message || 'Failed to purchase tokens');
    } finally {
      setIsTokenProcessing(false);
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

        {/* Category Tabs */}
        <div className="mb-6 flex gap-2 flex-wrap">
          {CATEGORIES.map(category => (
            <button
              key={category.id}
              onClick={() => setActiveCategory(category.id)}
              className={`px-4 py-2 rounded-lg font-semibold transition-colors ${
                activeCategory === category.id
                  ? 'bg-blue-600 hover:bg-blue-700 text-white'
                  : 'bg-gray-800 hover:bg-gray-700 text-gray-300 border border-gray-700'
              }`}
            >
              {category.label}
            </button>
          ))}
        </div>

        {/* Pokemon Card Category */}
        {activeCategory === 'pokemon' && (
          <>
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
          </>
        )}

        {/* Game Token Category */}
        {activeCategory === 'tokens' && (
          <TokenPurchaseSection
            userCoins={user?.coins || 0}
            onPurchase={handleTokenPurchase}
            isProcessing={isTokenProcessing}
            error={tokenPurchaseError}
            successMessage={tokenSuccessMessage}
          />
        )}

        {/* Items Category */}
        {activeCategory === 'items' && (
          <ItemsGrid
            items={inventory?.game_items || []}
            onPurchase={handleItemPurchase}
            onUse={handleUseItem}
            usingItem={usingItem}
            isProcessing={isProcessing || Boolean(itemLoadingError) || requiresItemReload}
            userCoins={user?.coins || 0}
            inventory={itemInventory}
            activeBoosts={activeBoosts}
          />
        )}

        {(itemLoadingError || requiresItemReload) && activeCategory === 'items' && (
          <div className="mt-4 bg-red-900/50 border border-red-500 rounded-lg p-4 text-center">
            <p className="text-red-200 mb-3">{itemLoadingError || 'Reload inventory before retrying your purchase.'}</p>
            <button onClick={loadItemInventory} className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg">Reload inventory</button>
          </div>
        )}

        {/* Item purchase error */}
        {itemPurchaseError && activeCategory === 'items' && (
          <div className="mt-4 bg-red-900/50 border border-red-500 rounded-lg p-4">
            <p className="text-red-200 text-center font-semibold">{itemPurchaseError}</p>
          </div>
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

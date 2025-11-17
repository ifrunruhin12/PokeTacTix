import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import PokemonCard from '../components/battle/PokemonCard';
import { getUserCards, getUserDeck, updateDeck, transformCardData } from '../services/card.service';

/**
 * DeckManager Component
 * Allows users to view their collection and manage their 5-card deck
 * Requirements: 12.1, 12.2, 12.3
 */
export default function DeckManager() {
  const [collection, setCollection] = useState([]);
  const [deck, setDeck] = useState([]);
  const [selectedCards, setSelectedCards] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  
  // Filter and sort states
  const [searchTerm, setSearchTerm] = useState('');
  const [typeFilter, setTypeFilter] = useState('all');
  const [rarityFilter, setRarityFilter] = useState('all');
  const [sortBy, setSortBy] = useState('level');
  const [sortOrder, setSortOrder] = useState('desc');

  // Load collection and deck on mount
  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [cardsData, deckData] = await Promise.all([
        getUserCards(),
        getUserDeck()
      ]);
      
      // Transform card data to include current stats
      const transformedCards = cardsData.map(transformCardData);
      const transformedDeck = deckData.map(transformCardData);
      
      setCollection(transformedCards);
      setDeck(transformedDeck);
      setSelectedCards(transformedDeck.map(card => card.id));
    } catch (err) {
      setError(err.message || 'Failed to load cards');
    } finally {
      setLoading(false);
    }
  };

  // Toggle card selection for deck
  const toggleCardSelection = (cardId) => {
    setSelectedCards(prev => {
      if (prev.includes(cardId)) {
        // Remove from selection
        return prev.filter(id => id !== cardId);
      } else {
        // Add to selection (max 5)
        if (prev.length >= 5) {
          setError('Deck can only contain 5 Pokemon');
          setTimeout(() => setError(null), 3000);
          return prev;
        }
        return [...prev, cardId];
      }
    });
    setSuccess(null);
  };

  // Save deck configuration
  const handleSaveDeck = async () => {
    // Validate deck has 1-5 cards
    if (selectedCards.length < 1) {
      setError('Deck must contain at least 1 Pokemon');
      return;
    }
    if (selectedCards.length > 5) {
      setError('Deck can only contain 5 Pokemon');
      return;
    }

    try {
      setSaving(true);
      setError(null);
      setSuccess(null);
      
      await updateDeck(selectedCards);
      
      // Reload data to get updated deck
      await loadData();
      
      setSuccess('Deck saved successfully!');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError(err.message || 'Failed to save deck');
    } finally {
      setSaving(false);
    }
  };

  // Reset to current deck
  const handleReset = () => {
    setSelectedCards(deck.map(card => card.id));
    setError(null);
    setSuccess(null);
  };

  // Get filtered and sorted collection
  const getFilteredCollection = () => {
    let filtered = [...collection];

    // Search filter
    if (searchTerm) {
      filtered = filtered.filter(card =>
        card.pokemon_name.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Type filter
    if (typeFilter !== 'all') {
      filtered = filtered.filter(card =>
        card.types.some(type => type.toLowerCase() === typeFilter.toLowerCase())
      );
    }

    // Rarity filter
    if (rarityFilter === 'legendary') {
      filtered = filtered.filter(card => card.is_legendary);
    } else if (rarityFilter === 'mythical') {
      filtered = filtered.filter(card => card.is_mythical);
    } else if (rarityFilter === 'common') {
      filtered = filtered.filter(card => !card.is_legendary && !card.is_mythical);
    }

    // Sort
    filtered.sort((a, b) => {
      let aVal, bVal;
      
      switch (sortBy) {
        case 'name':
          aVal = a.pokemon_name;
          bVal = b.pokemon_name;
          break;
        case 'level':
          aVal = a.level;
          bVal = b.level;
          break;
        case 'hp':
          aVal = a.hp_max;
          bVal = b.hp_max;
          break;
        case 'attack':
          aVal = a.attack;
          bVal = b.attack;
          break;
        default:
          return 0;
      }

      if (sortBy === 'name') {
        return sortOrder === 'asc' 
          ? aVal.localeCompare(bVal)
          : bVal.localeCompare(aVal);
      } else {
        return sortOrder === 'asc' ? aVal - bVal : bVal - aVal;
      }
    });

    return filtered;
  };

  // Get unique types from collection
  const getAvailableTypes = () => {
    const types = new Set();
    collection.forEach(card => {
      card.types.forEach(type => types.add(type));
    });
    return Array.from(types).sort();
  };

  const filteredCollection = getFilteredCollection();
  // Preserve selection order by mapping selectedCards to actual card objects
  const selectedCardsData = selectedCards.map(id => collection.find(card => card.id === id)).filter(Boolean);
  const hasChanges = JSON.stringify(selectedCards.sort()) !== JSON.stringify(deck.map(c => c.id).sort());

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 py-8 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading your collection...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Deck Manager</h1>
          <p className="text-gray-400">
            Select 1-5 Pokemon for your battle deck. You have {collection.length} Pokemon in your collection.
          </p>
        </div>

        {/* Error/Success Messages */}
        <AnimatePresence>
          {error && (
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="mb-6 bg-red-500/20 border border-red-500 text-red-200 px-4 py-3 rounded-lg"
            >
              {error}
            </motion.div>
          )}
          {success && (
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="mb-6 bg-green-500/20 border border-green-500 text-green-200 px-4 py-3 rounded-lg"
            >
              {success}
            </motion.div>
          )}
        </AnimatePresence>

        {/* Current Deck Section */}
        <div className="mb-8 bg-gray-800 rounded-lg p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-2xl font-bold">
              Current Deck ({selectedCards.length}/5)
            </h2>
            <div className="flex gap-2">
              {hasChanges && (
                <button
                  onClick={handleReset}
                  className="px-4 py-2 bg-gray-600 hover:bg-gray-700 rounded-lg transition-colors"
                  disabled={saving}
                >
                  Reset
                </button>
              )}
              <button
                onClick={handleSaveDeck}
                disabled={selectedCards.length < 1 || saving || !hasChanges}
                className={`px-6 py-2 rounded-lg font-semibold transition-colors ${
                  selectedCards.length >= 1 && hasChanges
                    ? 'bg-blue-500 hover:bg-blue-600 text-white'
                    : 'bg-gray-600 text-gray-400 cursor-not-allowed'
                }`}
              >
                {saving ? 'Saving...' : 'Save Deck'}
              </button>
            </div>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-4">
            {[0, 1, 2, 3, 4].map(index => {
              const card = selectedCardsData[index];
              return (
                <div key={index} className="aspect-[3/4]">
                  {card ? (
                    <PokemonCard
                      pokemon={card}
                      isActive={true}
                      onSelect={() => toggleCardSelection(card.id)}
                      className="h-full"
                    />
                  ) : (
                    <div className="h-full rounded-lg border-4 border-dashed border-gray-600 bg-gray-700/50 flex items-center justify-center">
                      <span className="text-gray-500 text-sm">Empty Slot</span>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        {/* Filters and Search */}
        <div className="mb-6 bg-gray-800 rounded-lg p-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Search */}
            <div>
              <label className="block text-sm font-medium mb-2">Search</label>
              <input
                type="text"
                placeholder="Pokemon name..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:border-blue-500"
              />
            </div>

            {/* Type Filter */}
            <div>
              <label className="block text-sm font-medium mb-2">Type</label>
              <select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:border-blue-500"
              >
                <option value="all">All Types</option>
                {getAvailableTypes().map(type => (
                  <option key={type} value={type}>{type}</option>
                ))}
              </select>
            </div>

            {/* Rarity Filter */}
            <div>
              <label className="block text-sm font-medium mb-2">Rarity</label>
              <select
                value={rarityFilter}
                onChange={(e) => setRarityFilter(e.target.value)}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:border-blue-500"
              >
                <option value="all">All Rarities</option>
                <option value="common">Common</option>
                <option value="legendary">Legendary</option>
                <option value="mythical">Mythical</option>
              </select>
            </div>

            {/* Sort */}
            <div>
              <label className="block text-sm font-medium mb-2">Sort By</label>
              <div className="flex gap-2">
                <select
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value)}
                  className="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg focus:outline-none focus:border-blue-500"
                >
                  <option value="level">Level</option>
                  <option value="name">Name</option>
                  <option value="hp">HP</option>
                  <option value="attack">Attack</option>
                </select>
                <button
                  onClick={() => setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc')}
                  className="px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg hover:bg-gray-600 transition-colors"
                  title={sortOrder === 'asc' ? 'Ascending' : 'Descending'}
                >
                  {sortOrder === 'asc' ? '↑' : '↓'}
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Collection Grid */}
        <div>
          <h2 className="text-2xl font-bold mb-4">
            Your Collection ({filteredCollection.length} Pokemon)
          </h2>
          
          {filteredCollection.length === 0 ? (
            <div className="text-center py-12 bg-gray-800 rounded-lg">
              <p className="text-gray-400 text-lg">No Pokemon found matching your filters</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
              {filteredCollection.map(card => (
                <motion.div
                  key={card.id}
                  className="aspect-[3/4] relative"
                  whileHover={{ scale: 1.02 }}
                  transition={{ duration: 0.2 }}
                >
                  <PokemonCard
                    pokemon={card}
                    isActive={selectedCards.includes(card.id)}
                    onSelect={() => toggleCardSelection(card.id)}
                    className="h-full"
                  />
                  {selectedCards.includes(card.id) && (
                    <div className="absolute top-2 left-2 bg-blue-500 text-white text-xs px-2 py-1 rounded-full font-bold z-10">
                      ✓ In Deck
                    </div>
                  )}
                </motion.div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

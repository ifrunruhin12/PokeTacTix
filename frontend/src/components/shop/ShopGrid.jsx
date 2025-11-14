import { useState, useMemo } from 'react';
import PropTypes from 'prop-types';
import ShopCard from './ShopCard';

/**
 * ShopGrid Component
 * Displays a grid of Pokemon cards available in the shop with filtering and sorting
 */
export default function ShopGrid({ 
  items, 
  onPurchase, 
  userCoins, 
  ownedPokemon,
  searchQuery,
  selectedRarity,
  sortBy,
  originalPrices,
}) {
  // Filter and sort items
  const filteredAndSortedItems = useMemo(() => {
    let result = [...items];

    // Apply search filter
    if (searchQuery) {
      result = result.filter(item =>
        item.pokemon_name.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    // Apply rarity filter
    if (selectedRarity && selectedRarity !== 'all') {
      result = result.filter(item => {
        if (selectedRarity === 'legendary') return item.is_legendary;
        if (selectedRarity === 'mythical') return item.is_mythical;
        return item.rarity === selectedRarity;
      });
    }

    // Apply sorting
    result.sort((a, b) => {
      switch (sortBy) {
        case 'price-low':
          return a.price - b.price;
        case 'price-high':
          return b.price - a.price;
        case 'name':
          return a.pokemon_name.localeCompare(b.pokemon_name);
        case 'rarity':
          const rarityOrder = { common: 1, uncommon: 2, rare: 3, legendary: 4, mythical: 5 };
          const aRarity = a.is_legendary ? 4 : a.is_mythical ? 5 : rarityOrder[a.rarity] || 0;
          const bRarity = b.is_legendary ? 4 : b.is_mythical ? 5 : rarityOrder[b.rarity] || 0;
          return bRarity - aRarity;
        default:
          return 0;
      }
    });

    return result;
  }, [items, searchQuery, selectedRarity, sortBy]);

  if (filteredAndSortedItems.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-400 text-lg">No Pokemon found matching your criteria.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {filteredAndSortedItems.map((item, index) => (
        <ShopCard
          key={`${item.pokemon_name}-${index}`}
          item={item}
          onPurchase={onPurchase}
          userCoins={userCoins}
          isOwned={ownedPokemon.includes(item.pokemon_name.toLowerCase())}
          originalPrice={originalPrices?.[item.pokemon_name]}
        />
      ))}
    </div>
  );
}

ShopGrid.propTypes = {
  items: PropTypes.arrayOf(
    PropTypes.shape({
      pokemon_name: PropTypes.string.isRequired,
      price: PropTypes.number.isRequired,
      rarity: PropTypes.string.isRequired,
      is_legendary: PropTypes.bool.isRequired,
      is_mythical: PropTypes.bool.isRequired,
    })
  ).isRequired,
  onPurchase: PropTypes.func.isRequired,
  userCoins: PropTypes.number.isRequired,
  ownedPokemon: PropTypes.arrayOf(PropTypes.string),
  searchQuery: PropTypes.string,
  selectedRarity: PropTypes.string,
  sortBy: PropTypes.string,
  originalPrices: PropTypes.object,
};

ShopGrid.defaultProps = {
  ownedPokemon: [],
  searchQuery: '',
  selectedRarity: 'all',
  sortBy: 'rarity',
  originalPrices: {},
};

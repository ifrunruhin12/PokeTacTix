import PropTypes from 'prop-types';

/**
 * ShopFilters Component
 * Provides search, filter, and sort controls for the shop
 */
export default function ShopFilters({ 
  searchQuery, 
  onSearchChange, 
  selectedRarity, 
  onRarityChange, 
  sortBy, 
  onSortChange 
}) {
  const rarityOptions = [
    { value: 'all', label: 'All Rarities', icon: '🎴' },
    { value: 'common', label: 'Common', icon: '⚪' },
    { value: 'uncommon', label: 'Uncommon', icon: '🟢' },
    { value: 'rare', label: 'Rare', icon: '🔵' },
    { value: 'legendary', label: 'Legendary', icon: '⭐' },
    { value: 'mythical', label: 'Mythical', icon: '✨' },
  ];

  const sortOptions = [
    { value: 'rarity', label: 'Rarity (High to Low)' },
    { value: 'price-low', label: 'Price (Low to High)' },
    { value: 'price-high', label: 'Price (High to Low)' },
    { value: 'name', label: 'Name (A-Z)' },
  ];

  return (
    <div className="bg-gray-800 rounded-lg p-4 mb-6 border border-gray-700">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Search */}
        <div>
          <label htmlFor="search" className="block text-sm font-medium text-gray-400 mb-2">
            Search Pokemon
          </label>
          <div className="relative">
            <input
              id="search"
              type="text"
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              placeholder="Search by name..."
              className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 pl-10 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <svg
              className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>
        </div>

        {/* Rarity Filter */}
        <div>
          <label htmlFor="rarity" className="block text-sm font-medium text-gray-400 mb-2">
            Filter by Rarity
          </label>
          <select
            id="rarity"
            value={selectedRarity}
            onChange={(e) => onRarityChange(e.target.value)}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {rarityOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.icon} {option.label}
              </option>
            ))}
          </select>
        </div>

        {/* Sort */}
        <div>
          <label htmlFor="sort" className="block text-sm font-medium text-gray-400 mb-2">
            Sort By
          </label>
          <select
            id="sort"
            value={sortBy}
            onChange={(e) => onSortChange(e.target.value)}
            className="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {sortOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Active Filters Display */}
      {(searchQuery || selectedRarity !== 'all') && (
        <div className="mt-4 flex items-center gap-2 flex-wrap">
          <span className="text-sm text-gray-400">Active filters:</span>
          {searchQuery && (
            <span className="bg-blue-600 text-white text-xs px-3 py-1 rounded-full flex items-center gap-1">
              Search: "{searchQuery}"
              <button
                onClick={() => onSearchChange('')}
                className="hover:text-gray-300"
                aria-label="Clear search"
              >
                ✕
              </button>
            </span>
          )}
          {selectedRarity !== 'all' && (
            <span className="bg-purple-600 text-white text-xs px-3 py-1 rounded-full flex items-center gap-1 capitalize">
              {selectedRarity}
              <button
                onClick={() => onRarityChange('all')}
                className="hover:text-gray-300"
                aria-label="Clear rarity filter"
              >
                ✕
              </button>
            </span>
          )}
        </div>
      )}
    </div>
  );
}

ShopFilters.propTypes = {
  searchQuery: PropTypes.string.isRequired,
  onSearchChange: PropTypes.func.isRequired,
  selectedRarity: PropTypes.string.isRequired,
  onRarityChange: PropTypes.func.isRequired,
  sortBy: PropTypes.string.isRequired,
  onSortChange: PropTypes.func.isRequired,
};

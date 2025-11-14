# Shop Components

This directory contains all components related to the Pokemon Shop feature.

## Components

### ShopGrid
Main grid component that displays Pokemon cards available for purchase.

**Props:**
- `items` - Array of shop items
- `onPurchase` - Callback when user clicks purchase
- `userCoins` - Current user's coin balance
- `ownedPokemon` - Array of owned Pokemon names
- `searchQuery` - Current search query
- `selectedRarity` - Selected rarity filter
- `sortBy` - Current sort option
- `originalPrices` - Object mapping Pokemon names to original prices (for discount display)

**Features:**
- Responsive grid layout (1-4 columns based on screen size)
- Filtering by search query and rarity
- Sorting by price, name, or rarity
- Shows "Owned" indicator for Pokemon already in collection

### ShopCard
Individual Pokemon card display in the shop.

**Props:**
- `item` - Shop item object with Pokemon details
- `onPurchase` - Callback when purchase button clicked
- `userCoins` - Current user's coin balance
- `isOwned` - Whether user already owns this Pokemon
- `originalPrice` - Original price before discount (optional)

**Features:**
- Rarity-based border colors (gold for legendary, rainbow for mythical)
- Type badges with official Pokemon colors
- Stats preview (HP, ATK, DEF, SPD)
- "Rare Appearance" badge for legendary/mythical
- "Owned" indicator
- Discount price display with strikethrough
- Hover animation with lift effect
- Disabled state when insufficient coins

### ShopFilters
Filter and search controls for the shop.

**Props:**
- `searchQuery` - Current search query
- `onSearchChange` - Callback for search input changes
- `selectedRarity` - Current rarity filter
- `onRarityChange` - Callback for rarity filter changes
- `sortBy` - Current sort option
- `onSortChange` - Callback for sort option changes

**Features:**
- Search by Pokemon name
- Filter by rarity (all, common, uncommon, rare, legendary, mythical)
- Sort by rarity, price (low/high), or name
- Active filters display with clear buttons

### PurchaseModal
Confirmation modal for Pokemon purchases.

**Props:**
- `isOpen` - Whether modal is visible
- `onClose` - Callback to close modal
- `item` - Shop item being purchased
- `userCoins` - Current user's coin balance
- `onConfirm` - Callback when purchase confirmed
- `isProcessing` - Whether purchase is in progress
- `error` - Error message to display

**Features:**
- Full Pokemon details display
- Price breakdown with remaining coins calculation
- Insufficient coins warning
- Loading state during purchase
- Error message display
- Backdrop blur effect
- Smooth animations with Framer Motion

### DiscountBanner
Banner displayed when a discount event is active.

**Props:**
- `discountPercent` - Discount percentage (0-100)
- `refreshTime` - ISO timestamp when discount ends

**Features:**
- Gradient background with animation
- Countdown timer showing time remaining
- Auto-hides when no discount active
- Animated emoji icon

## Usage Example

```jsx
import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import shopService from '../services/shop.service';
import {
  ShopGrid,
  ShopFilters,
  PurchaseModal,
  DiscountBanner
} from '../components/shop';

export default function Shop() {
  const { user, updateUser } = useAuth();
  const [inventory, setInventory] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedRarity, setSelectedRarity] = useState('all');
  const [sortBy, setSortBy] = useState('rarity');
  
  useEffect(() => {
    loadInventory();
  }, []);

  const loadInventory = async () => {
    const data = await shopService.getInventory();
    setInventory(data);
  };

  const handlePurchase = async (item) => {
    const result = await shopService.purchaseCard(item.pokemon_name);
    updateUser({ coins: result.remaining_coins });
  };

  return (
    <div>
      <DiscountBanner
        discountPercent={inventory?.discount_percent}
        refreshTime={inventory?.refresh_time}
      />
      
      <ShopFilters
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        selectedRarity={selectedRarity}
        onRarityChange={setSelectedRarity}
        sortBy={sortBy}
        onSortChange={setSortBy}
      />
      
      <ShopGrid
        items={inventory?.items || []}
        onPurchase={handlePurchase}
        userCoins={user?.coins || 0}
        searchQuery={searchQuery}
        selectedRarity={selectedRarity}
        sortBy={sortBy}
      />
    </div>
  );
}
```

## Styling

All components use Tailwind CSS for styling with the following color scheme:
- Background: gray-800, gray-900
- Borders: gray-700
- Text: white, gray-400
- Accent: blue-600, purple-600
- Success: green-400, green-600
- Warning: yellow-400
- Error: red-400, red-600

Rarity colors:
- Common: gray-600
- Uncommon: green-400
- Rare: blue-400
- Legendary: yellow-400 (gold)
- Mythical: purple-500 (rainbow gradient)

## API Integration

The shop components integrate with the backend API through `shop.service.js`:

- `GET /api/shop/inventory` - Fetch current shop inventory
- `POST /api/shop/purchase` - Purchase a Pokemon card

The service handles authentication tokens automatically through the API client.

## State Management

The Shop page manages the following state:
- Inventory data (items, discount info, refresh time)
- User's owned Pokemon (for "Owned" indicators)
- Filter/search/sort preferences
- Purchase modal state
- Loading and error states

User coin balance is managed through the AuthContext and updated after successful purchases.

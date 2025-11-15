# Deck Manager Component

## Overview

The Deck Manager is a comprehensive UI component that allows users to view their Pokemon collection and manage their 5-card battle deck. It provides filtering, sorting, and search capabilities to help users find and select their preferred Pokemon.

## Features

### 1. Deck Management (Requirement 12.1, 12.2, 12.3)
- **Current Deck Display**: Shows the 5 Pokemon currently in the user's deck at the top
- **Card Selection**: Click on any Pokemon card to add/remove it from the deck
- **Visual Feedback**: Selected cards show a "✓ In Deck" badge and active border
- **Validation**: Enforces exactly 5 Pokemon in the deck before saving
- **Save/Reset**: Save changes to the deck or reset to the current saved deck

### 2. Collection Display (Requirement 12.1)
- **Grid Layout**: Responsive grid showing all Pokemon in the user's collection
- **Card Details**: Each card displays:
  - Pokemon name and sprite
  - Type badges with official Pokemon colors
  - Current stats (HP, Stamina, Attack, Defense, Speed)
  - Level and XP progress
  - Rarity indicators (Legendary/Mythical badges)

### 3. Filtering (Requirement 12.1)
- **Search**: Filter by Pokemon name (case-insensitive)
- **Type Filter**: Filter by Pokemon type (Fire, Water, Grass, etc.)
- **Rarity Filter**: Filter by rarity (Common, Legendary, Mythical)

### 4. Sorting (Requirement 12.1)
- **Sort Options**: Sort by Level, Name, HP, or Attack
- **Sort Order**: Toggle between ascending and descending order
- **Visual Indicator**: Arrow icon shows current sort direction

## Component Structure

```
DeckManager (Page Component)
├── Current Deck Section
│   ├── Deck Cards (5 slots)
│   └── Save/Reset Buttons
├── Filters Section
│   ├── Search Input
│   ├── Type Filter
│   ├── Rarity Filter
│   └── Sort Controls
└── Collection Grid
    └── Pokemon Cards (clickable)
```

## API Integration

### Endpoints Used
- `GET /api/cards` - Fetch user's card collection
- `GET /api/cards/deck` - Fetch current deck
- `PUT /api/cards/deck` - Update deck configuration

### Data Flow
1. On mount, fetch collection and deck data
2. Transform backend data to include current stats based on level
3. User selects/deselects cards (max 5)
4. Validate deck has exactly 5 cards
5. Save deck to backend
6. Reload data to confirm changes

## State Management

### Local State
- `collection`: All Pokemon cards owned by the user
- `deck`: Current saved deck (5 cards)
- `selectedCards`: Array of card IDs currently selected for the deck
- `loading`: Loading state for initial data fetch
- `saving`: Loading state for save operation
- `error`: Error message to display
- `success`: Success message to display

### Filter State
- `searchTerm`: Search query string
- `typeFilter`: Selected type filter
- `rarityFilter`: Selected rarity filter
- `sortBy`: Sort field (level, name, hp, attack)
- `sortOrder`: Sort direction (asc, desc)

## User Experience

### Visual Feedback
- **Active Cards**: Blue border and glow effect for cards in the deck
- **Hover Effects**: Cards lift and scale on hover
- **Loading States**: Spinner during data fetch and save operations
- **Error Messages**: Red banner for errors (auto-dismiss after 3s)
- **Success Messages**: Green banner for successful saves (auto-dismiss after 3s)

### Validation
- **Deck Size**: Must have exactly 5 Pokemon
- **Save Button**: Disabled when deck is invalid or no changes made
- **Error Display**: Clear error message when trying to add more than 5 cards

### Responsive Design
- **Mobile**: 2 columns for cards
- **Tablet**: 3-4 columns for cards
- **Desktop**: 5 columns for cards
- **Filters**: Stack vertically on mobile, horizontal on desktop

## Card Data Transformation

The component transforms backend card data to include current stats:

```javascript
{
  id: 1,
  name: "Pikachu",
  level: 15,
  xp: 450,
  base_hp: 35,
  base_attack: 55,
  base_defense: 40,
  base_speed: 90,
  // Calculated current stats based on level
  hp: 50,
  hp_max: 50,
  attack: 71,
  defense: 52,
  speed: 103,
  stamina: 206,
  stamina_max: 206
}
```

## Performance Considerations

- **Memoization**: Filter and sort operations are optimized
- **Lazy Loading**: Cards render efficiently with Framer Motion
- **Debouncing**: Search input could be debounced for large collections
- **Batch Updates**: State updates are batched for better performance

## Future Enhancements

- Drag-and-drop deck building
- Deck presets/templates
- Card comparison view
- Deck statistics (average level, type distribution)
- Export/import deck configurations
- Deck recommendations based on type coverage

## Testing

### Manual Testing Checklist
- [ ] Load collection and deck on mount
- [ ] Select/deselect cards (max 5)
- [ ] Save deck with exactly 5 cards
- [ ] Try to save with < 5 or > 5 cards (should show error)
- [ ] Reset to current deck
- [ ] Search by Pokemon name
- [ ] Filter by type
- [ ] Filter by rarity
- [ ] Sort by different fields
- [ ] Toggle sort order
- [ ] Verify visual feedback (borders, badges, hover effects)
- [ ] Test responsive layout on different screen sizes

## Related Components

- `PokemonCard`: Reusable card component from battle system
- `card.service.js`: API service for card operations
- `api.js`: Base API client with authentication

## Requirements Mapping

- **Requirement 12.1**: Display all Pokemon in user's collection ✓
- **Requirement 12.2**: Show current deck (5 selected Pokemon) at top ✓
- **Requirement 12.3**: Update deck (must have exactly 5 cards) ✓
- **Requirement 12.5**: Display success message on save ✓

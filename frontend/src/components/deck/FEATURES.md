# Deck Manager Features Guide

## Visual Layout

```
┌─────────────────────────────────────────────────────────────┐
│                      DECK MANAGER                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Current Deck (3/5)                    [Reset] [Save Deck]  │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐                       │
│  │ 🔵 │ │ 🔵 │ │ 🔵 │ │ -- │ │ -- │                       │
│  │Card│ │Card│ │Card│ │Empt│ │Empt│                       │
│  │ 1  │ │ 2  │ │ 3  │ │ y  │ │ y  │                       │
│  └────┘ └────┘ └────┘ └────┘ └────┘                       │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│  Filters                                                     │
│  [Search: ____] [Type: All ▼] [Rarity: All ▼] [Sort: Lv ▼↓]│
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Your Collection (25 Pokemon)                               │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐                       │
│  │✓In │ │    │ │✓In │ │    │ │    │                       │
│  │Card│ │Card│ │Card│ │Card│ │Card│                       │
│  │ 1  │ │ 4  │ │ 2  │ │ 5  │ │ 6  │                       │
│  └────┘ └────┘ └────┘ └────┘ └────┘                       │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐                       │
│  │    │ │✓In │ │    │ │    │ │    │                       │
│  │Card│ │Card│ │Card│ │Card│ │Card│                       │
│  │ 7  │ │ 3  │ │ 8  │ │ 9  │ │ 10 │                       │
│  └────┘ └────┘ └────┘ └────┘ └────┘                       │
│  ... (more cards)                                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Feature Breakdown

### 1. Current Deck Section
**Location**: Top of page
**Purpose**: Shows the 5 Pokemon currently selected for battle

**Visual Elements**:
- 5 card slots in a row
- Selected cards show full Pokemon details
- Empty slots show dashed border with "Empty Slot" text
- Counter shows "X/5" cards selected
- Blue glow effect on selected cards

**Interactions**:
- Click a card in the deck to remove it
- Cards animate when added/removed

**Buttons**:
- **Reset**: Reverts to last saved deck (only shown if changes exist)
- **Save Deck**: Saves current selection (disabled if not exactly 5 cards)

### 2. Filters Section
**Location**: Below deck, above collection
**Purpose**: Help users find specific Pokemon

**Filter Options**:

#### Search Input
- **Type**: Text input
- **Placeholder**: "Pokemon name..."
- **Behavior**: Case-insensitive, real-time filtering
- **Example**: Type "pika" to find Pikachu

#### Type Filter
- **Type**: Dropdown
- **Options**: All Types, Fire, Water, Grass, Electric, etc.
- **Behavior**: Shows only Pokemon with selected type
- **Dynamic**: Only shows types present in collection

#### Rarity Filter
- **Type**: Dropdown
- **Options**: All Rarities, Common, Legendary, Mythical
- **Behavior**: 
  - Common = non-legendary and non-mythical
  - Legendary = gold border Pokemon
  - Mythical = rainbow border Pokemon

#### Sort Controls
- **Type**: Dropdown + Toggle button
- **Sort By Options**: Level, Name, HP, Attack
- **Sort Order**: Ascending (↑) or Descending (↓)
- **Default**: Level descending (highest first)

### 3. Collection Grid
**Location**: Bottom section
**Purpose**: Display all Pokemon with filtering applied

**Layout**:
- Responsive grid (2-5 columns based on screen size)
- Cards maintain aspect ratio
- Smooth animations on hover

**Card Display**:
Each card shows:
- Pokemon name (always visible, no truncation)
- Type badges with official colors
- Pokemon sprite/image
- HP bar (green, turns red when low)
- Stamina bar (blue)
- Attack, Defense, Speed stats
- Level indicator
- XP progress bar (if not max level)
- "MAX LEVEL" badge (if level 50)
- Rarity badge (⭐ Legendary or ✨ Mythical)

**Selection Indicator**:
- "✓ In Deck" badge on selected cards
- Blue border and glow effect
- Badge positioned at top-left corner

**Interactions**:
- Click any card to add/remove from deck
- Hover effect: card lifts and scales slightly
- Error message if trying to add 6th card

## User Workflows

### Workflow 1: Building a New Deck
1. User navigates to Deck Manager
2. Current deck loads (may be empty or have previous selection)
3. User clicks on 5 Pokemon from collection
4. Selected cards appear in deck slots at top
5. "✓ In Deck" badges appear on selected cards
6. Save button becomes enabled when exactly 5 selected
7. User clicks "Save Deck"
8. Success message appears: "Deck saved successfully!"
9. Deck is now ready for battles

### Workflow 2: Modifying Existing Deck
1. User navigates to Deck Manager
2. Current deck loads with 5 Pokemon
3. User clicks a card in the deck to remove it
4. Card is removed, slot becomes empty
5. User clicks a different card from collection
6. New card appears in empty slot
7. Reset button appears (can undo changes)
8. User clicks "Save Deck"
9. Deck is updated

### Workflow 3: Finding Specific Pokemon
1. User wants to find all Fire-type Pokemon
2. User selects "Fire" from Type filter
3. Collection updates to show only Fire types
4. User wants highest level Fire Pokemon
5. User selects "Level" sort and "↓" descending
6. Fire Pokemon sorted by level (highest first)
7. User selects desired Pokemon for deck

### Workflow 4: Searching by Name
1. User wants to add Pikachu to deck
2. User types "pika" in search box
3. Collection filters to show only Pikachu
4. User clicks Pikachu card
5. Pikachu added to deck
6. User clears search to see all Pokemon again

## Visual Feedback

### Success States
- ✅ Green banner: "Deck saved successfully!"
- ✓ Badge on selected cards
- Blue glow on active cards
- Smooth animations

### Error States
- ❌ Red banner: "Deck must contain exactly 5 Pokemon"
- Disabled save button (gray)
- Error message when trying to add 6th card

### Loading States
- Spinner during initial load
- "Saving..." text on save button
- Disabled buttons during operations

### Empty States
- "Empty Slot" in deck slots
- "No Pokemon found matching your filters" in collection

## Responsive Behavior

### Desktop (1024px+)
- 5 columns in collection grid
- Horizontal filter layout
- All elements visible

### Tablet (768px - 1023px)
- 3-4 columns in collection grid
- Horizontal filter layout
- Slightly smaller cards

### Mobile (< 768px)
- 2 columns in collection grid
- Filters may stack vertically
- Smaller cards but still readable
- Touch-friendly tap targets

## Accessibility

- Keyboard navigation support (inherited from PokemonCard)
- Clear visual indicators for selected state
- High contrast colors
- Readable text sizes
- Error messages announced to screen readers
- Focus states on interactive elements

## Performance

- Efficient filtering (no unnecessary re-renders)
- Optimized sorting algorithms
- Smooth animations (60fps)
- Fast load times
- Minimal API calls (only on load and save)

## Integration with Other Features

### Battle System
- Deck saved here is used in battles
- Must have exactly 5 Pokemon to start battle
- Deck can be updated between battles

### Shop System
- Newly purchased Pokemon appear in collection
- Can immediately add to deck
- Collection count updates automatically

### Profile/Stats
- Deck composition may affect battle strategies
- Type diversity can be tracked
- Level distribution visible

## Tips for Users

1. **Build Balanced Decks**: Include different types for type coverage
2. **Level Up Pokemon**: Higher level = better stats
3. **Use Filters**: Find Pokemon quickly with filters
4. **Check Stats**: Compare Pokemon stats before selecting
5. **Save Often**: Changes aren't saved until you click "Save Deck"
6. **Reset if Needed**: Use Reset button to undo unwanted changes
7. **Legendary Pokemon**: Powerful but rare, use wisely
8. **Type Advantage**: Consider type matchups when building deck

## Common Questions

**Q: Can I have duplicate Pokemon in my deck?**
A: Yes! If you own multiple copies, you can add them all to your deck.

**Q: What happens if I don't save my changes?**
A: Changes are lost when you navigate away. Always click "Save Deck".

**Q: Can I have less than 5 Pokemon in my deck?**
A: No, you must have exactly 5 Pokemon to save and battle.

**Q: How do I remove a Pokemon from my deck?**
A: Click on the card in the deck section or in the collection (if it has the "✓ In Deck" badge).

**Q: Why can't I add a 6th Pokemon?**
A: Decks are limited to 5 Pokemon for balanced gameplay.

**Q: Do filters affect my deck?**
A: No, filters only affect what you see in the collection. Your deck remains unchanged.

**Q: Can I sort my deck?**
A: The deck shows cards in the order you selected them. Sorting only affects the collection view.

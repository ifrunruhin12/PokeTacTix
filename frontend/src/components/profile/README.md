# Profile Components

This directory contains all components related to the user profile and statistics display.

## Components

### Dashboard.jsx
Main dashboard component that displays:
- User information (username, coins, join date)
- Overall battle statistics (total battles, wins, losses, win rate)
- Separate statistics for 1v1 and 5v5 battles
- Highest level Pokémon display
- Total coins earned

**Requirements**: 8.1, 8.5, 9.5

### BattleHistory.jsx
Displays the user's battle history:
- Last 20 battles by default (expandable to 100)
- Table view on desktop, card view on mobile
- Shows date, mode, result, coins earned, and duration for each battle
- Load more functionality with pagination
- Responsive design

**Requirements**: 8.3

### Achievements.jsx
Achievement tracking and display:
- Shows all achievements with locked/unlocked status
- Progress bar showing overall completion
- Visual indicators for newly unlocked achievements
- Check progress button to evaluate and unlock new achievements
- Achievement cards with icons, names, and descriptions
- Displays unlock dates for completed achievements

**Requirements**: 8.4

### StatsVisualizations.jsx
Visual statistics and charts:
- **Win/Loss Pie Chart**: Visual representation of win rate with color-coded segments
- **Battle Timeline**: Chronological display of recent battles grouped by date
- **Level Distribution**: Bar chart showing distribution of Pokémon levels in collection
- Collection statistics (total Pokémon, highest/lowest levels, average level)

**Requirements**: 8.1

## Services

### stats.service.js
API client for statistics endpoints:
- `getPlayerStats()` - Fetch player statistics
- `getBattleHistory(limit)` - Fetch battle history with optional limit
- `getAchievements()` - Fetch all achievements with unlock status
- `checkAchievements()` - Check and unlock new achievements

## Usage

The Profile page (`pages/Profile.jsx`) uses a tabbed interface to display all profile components:

```jsx
import { Dashboard, BattleHistory, Achievements, StatsVisualizations } from '../components/profile';

// In Profile page
<Dashboard />           // Tab 1: Overview
<BattleHistory />       // Tab 2: Battle history
<Achievements />        // Tab 3: Achievements
<StatsVisualizations /> // Tab 4: Visual stats
```

## Features

### Responsive Design
- Desktop: Full table layouts and side-by-side displays
- Mobile: Card-based layouts and stacked views
- Adaptive typography and spacing

### Animations
- Framer Motion for smooth transitions
- Staggered animations for lists
- Progress bar animations
- Card hover effects
- Tab switching animations

### Data Loading
- Loading states with spinners
- Error handling with user-friendly messages
- Parallel data fetching for performance
- Empty states for new users

### Visual Design
- Color-coded statistics (green for wins, red for losses, etc.)
- Gradient backgrounds for emphasis
- Type-based color schemes
- Icon-based visual language
- Consistent spacing and borders

## API Endpoints Used

- `GET /api/profile/stats` - Player statistics
- `GET /api/profile/history?limit=20` - Battle history
- `GET /api/profile/achievements` - All achievements
- `POST /api/profile/achievements/check` - Check for new achievements
- `GET /api/cards` - User's card collection (for level distribution)

## Future Enhancements

Potential improvements:
- Export battle history to CSV
- Share achievements on social media
- Compare stats with friends
- More detailed battle analytics
- Achievement progress tracking
- Custom date range filters for history
- Interactive charts with tooltips

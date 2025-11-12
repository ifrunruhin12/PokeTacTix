# Battle Components

This directory contains all Pokemon card components for the battle system.

## Components

### PokemonCard
The base Pokemon card component that displays all Pokemon information including stats, level, XP, and type badges.

**Props:**
- `pokemon` (object): Pokemon data including name, types, sprite, stats, level, XP
- `isActive` (bool): Whether this Pokemon is currently active in battle
- `isFaceDown` (bool): Whether to show the card face-down (Pokeball image)
- `isKnockedOut` (bool): Whether the Pokemon is knocked out
- `onSelect` (func): Callback when card is clicked
- `className` (string): Additional CSS classes

**Features:**
- Rarity-based borders (gold for legendary, rainbow for mythical)
- Animated stat bars with color coding
- Hover effects with lift and shadow
- Level and XP progress display
- Type badges with official Pokemon colors

**Example:**
```jsx
import { PokemonCard } from './components/battle';

<PokemonCard
  pokemon={{
    name: 'Pikachu',
    types: ['electric'],
    sprite: 'https://...',
    hp: 80,
    hp_max: 100,
    stamina: 50,
    stamina_max: 60,
    attack: 55,
    defense: 40,
    speed: 90,
    level: 15,
    xp: 450,
    is_legendary: false,
    is_mythical: false
  }}
  isActive={true}
  onSelect={() => console.log('Card selected')}
/>
```

### StatBar
A reusable animated stat bar component with color coding.

**Props:**
- `label` (string): Label text (e.g., "HP", "Stamina")
- `current` (number): Current value
- `max` (number): Maximum value
- `type` (string): Stat type for color coding ('hp', 'stamina', 'attack', 'defense', 'speed', 'xp')
- `showValues` (bool): Whether to show current/max values
- `height` (string): Tailwind height class (default: 'h-3')
- `className` (string): Additional CSS classes

**Color Coding:**
- HP: Green (changes to orange/red when low)
- Stamina: Blue
- Attack: Red
- Defense: Yellow
- Speed: Purple
- XP: Violet

**Example:**
```jsx
import { StatBar } from './components/battle';

<StatBar
  label="HP"
  current={80}
  max={100}
  type="hp"
  showValues={true}
/>
```

### FlippableCard
Wraps PokemonCard with 3D flip animation for AI Pokemon that switch during battle.

**Props:**
- `pokemon` (object): Pokemon data
- `isFlipped` (bool): Whether the card is flipped (showing back)
- `isActive` (bool): Whether this Pokemon is active
- `isKnockedOut` (bool): Whether the Pokemon is knocked out
- `onSelect` (func): Callback when card is clicked
- `className` (string): Additional CSS classes

**Features:**
- 3D card flip animation (0.6s duration)
- Shows Pokeball on back side
- Shows Pokemon details on front side
- Smooth easeInOut transition

**Example:**
```jsx
import { FlippableCard } from './components/battle';

<FlippableCard
  pokemon={aiPokemon}
  isFlipped={!isActive}
  isActive={isActive}
  onSelect={() => console.log('AI card selected')}
/>
```

### AnimatedPokemonCard
Wraps PokemonCard with damage and knockout animations.

**Props:**
- `pokemon` (object): Pokemon data
- `isActive` (bool): Whether this Pokemon is active
- `isFaceDown` (bool): Whether to show face-down
- `isKnockedOut` (bool): Whether the Pokemon is knocked out
- `onDamage` (func): Callback when damage is detected (receives damage amount)
- `onSelect` (func): Callback when card is clicked
- `className` (string): Additional CSS classes

**Features:**
- Automatic shake animation when HP decreases
- Red flash overlay on damage
- Knockout animation (grayscale, fade, scale down)
- Animated red X overlay for knocked out Pokemon
- Pulsing glow effect on knockout X

**Example:**
```jsx
import { AnimatedPokemonCard } from './components/battle';

<AnimatedPokemonCard
  pokemon={playerPokemon}
  isActive={true}
  isKnockedOut={playerPokemon.hp <= 0}
  onDamage={(damage) => console.log(`Took ${damage} damage!`)}
  onSelect={() => console.log('Card selected')}
/>
```

## Usage in Battle Arena

Here's how to use these components together in a battle scene:

```jsx
import { AnimatedPokemonCard, FlippableCard } from './components/battle';

function BattleArena({ playerDeck, aiDeck, activePlayerIdx, activeAiIdx }) {
  return (
    <div className="flex justify-between p-8">
      {/* Player's active Pokemon */}
      <AnimatedPokemonCard
        pokemon={playerDeck[activePlayerIdx]}
        isActive={true}
        isKnockedOut={playerDeck[activePlayerIdx].hp <= 0}
        onDamage={(damage) => console.log(`Player took ${damage} damage`)}
      />

      {/* AI's active Pokemon with flip animation */}
      <FlippableCard
        pokemon={aiDeck[activeAiIdx]}
        isFlipped={false}
        isActive={true}
        isKnockedOut={aiDeck[activeAiIdx].hp <= 0}
      />

      {/* Inactive AI Pokemon (face-down) */}
      {aiDeck.map((pokemon, idx) => (
        idx !== activeAiIdx && (
          <FlippableCard
            key={idx}
            pokemon={pokemon}
            isFlipped={true}
            isActive={false}
          />
        )
      ))}
    </div>
  );
}
```

## Type Colors Reference

The components use official Pokemon type colors:

- Normal: #A8A878
- Fire: #F08030
- Water: #6890F0
- Electric: #F8D030
- Grass: #78C850
- Ice: #98D8D8
- Fighting: #C03028
- Poison: #A040A0
- Ground: #E0C068
- Flying: #A890F0
- Psychic: #F85888
- Bug: #A8B820
- Rock: #B8A038
- Ghost: #705898
- Dragon: #7038F8
- Dark: #705848
- Steel: #B8B8D0
- Fairy: #EE99AC

## Animation Details

All animations use Framer Motion for smooth, performant transitions:

- **HP Bar**: 0.5s easeOut transition, color changes based on percentage
- **Damage Shake**: 0.5s shake with 7 keyframes
- **Card Flip**: 0.6s easeInOut 3D rotation
- **Knockout**: 1s easeOut with grayscale, fade, and scale
- **Hover Effects**: 0.2s lift and shadow on hover
- **XP Bar**: 0.5s animated width transition

## Requirements Covered

This implementation satisfies the following requirements:

- **17.1**: Card-first design with holographic/gradient backgrounds
- **17.2**: Color-coded stat bars (green HP, blue stamina, red attack, yellow defense, purple speed)
- **17.3**: Official Pokemon type colors and icons
- **17.4**: Face-down state with Pokeball image and subtle animation
- **5.3**: Card visibility mechanic for AI Pokemon
- **6.4**: Knockout and damage animations

# PokeTacTix CLI - Battle Guide

## Starting a Battle

To start a battle in the CLI, use the `battle` command (or shortcut `b`):

```bash
> battle
```

or

```bash
> b
```

## Battle Modes

When you start a battle, you'll be prompted to choose a mode:

### 1v1 Battle

- Quick battle with one random Pokemon from your deck
- Rewards: 50 coins on victory, 20 XP
- Faster gameplay for quick sessions

### 5v5 Battle

- Full battle with all 5 Pokemon in your deck
- Rewards: 150 coins on victory, 15 XP per Pokemon
- Victory bonus: Choose one Pokemon from AI's team to add to your collection
- More strategic gameplay with Pokemon switching

## Battle Actions

During battle, you can choose from these actions:

1. **Attack** - Select a move to attack the opponent
   - Each move has different power and stamina cost
   - Type effectiveness applies
   - Must have enough stamina to use the move

2. **Defend** - Reduce incoming damage
   - Costs stamina based on your Pokemon's max HP
   - Reduces damage by your Defense stat
   - Good for conserving HP when low on stamina

3. **Pass** - Skip your turn
   - No stamina cost
   - Useful when you want to regenerate stamina
   - Warning: 3 consecutive passes by both players = draw

4. **Sacrifice** - Trade HP for stamina
   - 1st sacrifice: -10 HP, +50% max stamina
   - 2nd sacrifice: -15 HP, +25% max stamina
   - 3rd sacrifice: -20 HP, +15% max stamina
   - Maximum 3 sacrifices per Pokemon
   - Can only sacrifice when stamina is below 50%

5. **Surrender** - Give up the battle
   - In 1v1: Ends the battle immediately (loss)
   - In 5v5: Knocks out current Pokemon, must switch to another

## Pokemon Switching (5v5 Only)

When your active Pokemon is knocked out in 5v5 mode:

- You'll be prompted to select another Pokemon
- Cannot select knocked out Pokemon
- View HP, stamina, types, and stats before choosing
- New Pokemon enters at current HP/stamina (no restoration)

## Battle Rewards

### Victory Rewards

- **1v1**: 50 coins, 20 XP to your Pokemon
- **5v5**: 150 coins, 15 XP to each Pokemon in your deck
- **5v5 Bonus**: Choose one Pokemon from AI's team to add to your collection

### Loss Rewards

- **1v1**: 10 coins
- **5v5**: 25 coins
- No XP awarded

### Draw Rewards

- **1v1**: 25 coins, 10 XP
- **5v5**: 75 coins, 8 XP per Pokemon

## Level Up System

- Pokemon gain XP after battles
- 100 XP required per level
- Stats increase with each level:
  - HP: +3% per level
  - Attack: +2% per level
  - Defense: +2% per level
  - Speed: +1% per level
  - Stamina: Speed × 2

## Tips

1. **Manage Stamina**: Keep an eye on stamina costs for moves
2. **Use Sacrifice Wisely**: Only when stamina is critically low
3. **Type Advantage**: Moves are color-coded by type for easy identification
4. **Defend Strategically**: Use defend when you need to recover stamina
5. **5v5 Strategy**: Save your strongest Pokemon for later rounds
6. **Post-Battle Selection**: In 5v5 victories, choose Pokemon that complement your deck

## Requirements

Before starting a battle:

- You must have exactly 5 Pokemon in your deck
- Use `deck` command to view your current deck
- Use `collection` command to see all available Pokemon

## Example Session

```bash
> battle
[Battle mode selection menu appears]
Enter your choice (1-3): 1

Starting 1v1 battle with Pikachu!
Press Enter to begin...

[Battle screen displays]
[Choose actions each turn]
[Battle continues until victory/defeat/draw]

🎉 VICTORY! 🎉
Coins earned: +50 (Total: 150)
Experience gained:
  Pikachu: +20 XP → LEVEL UP! 1 → 2
  New stats: HP: 103, ATK: 56, DEF: 41, SPD: 91

✓ Game saved successfully
Press Enter to continue...
```

## Troubleshooting

**"You don't have any Pokemon in your deck"**

- Your deck is empty. This shouldn't happen after setup, but if it does, contact support.

**"Your deck must have exactly 5 Pokemon"**

- Deck editing is not yet implemented. You should have 5 Pokemon from the starter deck.

**"Failed to load player deck"**

- Your save file may be corrupted. Try the `reset` command to start fresh.

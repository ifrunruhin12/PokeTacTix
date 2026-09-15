package battle

import (
	"context"

	"pokemon-cli/internal/items"
)

// BoostProvider is the item-domain surface the battle system needs: the
// player's active deck buffs, and the per-battle tick that expires them.
// Implemented by the items service.
type BoostProvider interface {
	ActiveBoosts(ctx context.Context, userID int) ([]items.ActiveBoost, error)
	ConsumeBattleBoosts(ctx context.Context, userID int) error
}

// ApplyBoosts mutates the player's deck in place, adding every active boost's
// flat bonus to the boosted stat of each Pokemon. HP boosts raise both current
// and maximum HP; attack boosts raise attack. Multiple boosts stack.
func ApplyBoosts(deck []BattleCard, boosts []items.ActiveBoost) {
	for i := range deck {
		card := &deck[i]
		for _, b := range boosts {
			switch b.Stat {
			case items.StatAttack:
				card.Attack += b.Bonus
			case items.StatHP:
				card.HP += b.Bonus
				card.HPMax += b.Bonus
			}
		}
	}
}

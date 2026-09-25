package battle

import (
	"testing"

	"pokemon-cli/internal/items"

	"github.com/stretchr/testify/assert"
)

func TestApplyBoosts(t *testing.T) {
	deck := []BattleCard{
		{Name: "pikachu", HP: 30, HPMax: 30, Attack: 55},
		{Name: "snorlax", HP: 90, HPMax: 90, Attack: 70},
	}

	t.Run("attack and hp boosts apply to every card and stack", func(t *testing.T) {
		boosted := append([]BattleCard(nil), deck...)
		ApplyBoosts(boosted, []items.ActiveBoost{
			{Stat: items.StatAttack, Bonus: 5, BattlesRemaining: 3},
			{Stat: items.StatHP, Bonus: 10, BattlesRemaining: 4},
			{Stat: items.StatAttack, Bonus: 3, BattlesRemaining: 1}, // stacking
		})

		for i, c := range boosted {
			assert.Equal(t, deck[i].Attack+8, c.Attack, "base + 5 + 3")
			assert.Equal(t, deck[i].HP+10, c.HP, "current HP rises with the boost")
			assert.Equal(t, deck[i].HPMax+10, c.HPMax)
		}
	})

	t.Run("no boosts leaves the deck unchanged", func(t *testing.T) {
		unchanged := append([]BattleCard(nil), deck...)
		ApplyBoosts(unchanged, nil)
		assert.Equal(t, deck, unchanged)
	})

	t.Run("unknown stats are ignored", func(t *testing.T) {
		// Defensive: a bad stat in the database must not corrupt the deck.
		boosted := append([]BattleCard(nil), deck...)
		ApplyBoosts(boosted, []items.ActiveBoost{{Stat: "speed", Bonus: 99, BattlesRemaining: 1}})
		assert.Equal(t, deck, boosted)
	})
}

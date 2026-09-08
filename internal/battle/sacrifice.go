package battle

import "fmt"

// sacrificeResult holds the outcome of a sacrifice action
type sacrificeResult struct {
	HPLost       int
	StaminaGained int
}

// sacrificeCost returns the HP cost and stamina gain fraction for the Nth sacrifice.
// Both player and AI use the same table.
// count is the number of sacrifices already used (0 = first sacrifice).
func sacrificeCost(count int) (hpCost int, staminaGainFrac float64, err error) {
	switch count {
	case 0:
		return 10, 0.50, nil
	case 1:
		return 15, 0.25, nil
	case 2:
		return 20, 0.15, nil
	default:
		return 0, 0, fmt.Errorf("maximum sacrifices reached for this Pokemon (3 per battle)")
	}
}

// applySacrifice executes the sacrifice mechanic on a BattleCard.
// It is the single authoritative implementation used by both player and AI,
// in both 1v1 and 5v5 modes.
//
// Rules:
//   - Max 3 sacrifices per Pokemon per battle
//   - Pokemon must have strictly more HP than the cost
//   - Stamina must be below 50% of maxStamina (Speed * 2)
//   - Stamina is capped at maxStamina after the gain
func applySacrifice(card *BattleCard, sacrificeCount int) (sacrificeResult, error) {
	maxStamina := card.Speed * 2

	hpCost, gainFrac, err := sacrificeCost(sacrificeCount)
	if err != nil {
		return sacrificeResult{}, err
	}

	if card.HP <= hpCost {
		return sacrificeResult{}, fmt.Errorf("insufficient HP to sacrifice (need more than %d HP, have %d)", hpCost, card.HP)
	}

	if card.Stamina >= maxStamina/2 {
		return sacrificeResult{}, fmt.Errorf("stamina must be below 50%% of max to sacrifice (current: %d, max: %d)", card.Stamina, maxStamina)
	}

	oldHP := card.HP
	oldStamina := card.Stamina

	card.HP -= hpCost
	gain := int(float64(maxStamina) * gainFrac)
	card.Stamina += gain
	if card.Stamina > maxStamina {
		card.Stamina = maxStamina
	}

	return sacrificeResult{
		HPLost:        oldHP - card.HP,
		StaminaGained: card.Stamina - oldStamina,
	}, nil
}

// canSacrifice returns true if the given card is eligible to sacrifice
// (used by AI decision logic to avoid choosing sacrifice when it would fail)
func canSacrifice(card *BattleCard, sacrificeCount int) bool {
	if sacrificeCount >= 3 {
		return false
	}
	maxStamina := card.Speed * 2
	hpCost, _, _ := sacrificeCost(sacrificeCount)
	return card.HP > hpCost && card.Stamina < maxStamina/2
}

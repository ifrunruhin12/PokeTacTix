package game

import (
	"math/rand"
	"pokemon-cli/pokemon"
	"slices"
	"strings"
)

// Legendary and mythical Pokémon names (same as in fetch.go)
var legendaryNames = []string{
	"articuno", "zapdos", "moltres", "mewtwo", "raikou", "entei", "suicune", "lugia", "ho-oh", "regirock", "regice", "registeel", "latias", "latios", "kyogre", "groudon", "rayquaza", "uxie", "mesprit", "azelf", "dialga", "palkia", "heatran", "regigigas", "giratina", "cresselia", "cobalion", "terrakion", "virizion", "tornadus", "thundurus", "reshiram", "zekrom", "landorus", "kyurem", "xerneas", "yveltal", "zygarde", "tapu-koko", "tapu-lele", "tapu-bulu", "tapu-fini", "cosmog", "cosmoem", "solgaleo", "lunala", "necrozma", "zamacenta", "zacian", "eternatus", "kubfu", "urshifu", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus", "ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

var mythicalNames = []string{
	"mew", "celebi", "jirachi", "deoxys", "phione", "manaphy", "darkrai", "shaymin", "arceus", "victini", "keldeo", "meloetta", "genesect", "diancie", "hoopa", "volcanion", "magearna", "marshadow", "zeraora", "meltan", "melmetal", "zarude", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus", "ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

// Type effectiveness chart - maps attacking type to defending types and their multipliers
var typeChart = map[string]map[string]float64{
	"normal": {
		"rock": 0.5, "ghost": 0.0, "steel": 0.5,
	},
	"fire": {
		"fire": 0.5, "water": 0.5, "grass": 2.0, "ice": 2.0, "bug": 2.0, "rock": 0.5, "dragon": 0.5, "steel": 2.0,
	},
	"water": {
		"fire": 2.0, "water": 0.5, "grass": 0.5, "ground": 2.0, "rock": 2.0, "dragon": 0.5,
	},
	"electric": {
		"water": 2.0, "electric": 0.5, "grass": 0.5, "ground": 0.0, "flying": 2.0, "dragon": 0.5,
	},
	"grass": {
		"fire": 0.5, "water": 2.0, "grass": 0.5, "poison": 0.5, "ground": 2.0, "flying": 0.5, "bug": 0.5, "rock": 2.0, "dragon": 0.5, "steel": 0.5,
	},
	"ice": {
		"fire": 0.5, "water": 0.5, "grass": 2.0, "ice": 0.5, "ground": 2.0, "flying": 2.0, "dragon": 2.0, "steel": 0.5,
	},
	"fighting": {
		"normal": 2.0, "ice": 2.0, "poison": 0.5, "flying": 0.5, "psychic": 0.5, "bug": 0.5, "rock": 2.0, "ghost": 0.0, "dark": 2.0, "steel": 2.0, "fairy": 0.5,
	},
	"poison": {
		"grass": 2.0, "poison": 0.5, "ground": 0.5, "rock": 0.5, "ghost": 0.5, "steel": 0.0, "fairy": 2.0,
	},
	"ground": {
		"fire": 2.0, "electric": 2.0, "grass": 0.5, "poison": 2.0, "flying": 0.0, "bug": 0.5, "rock": 2.0, "steel": 2.0,
	},
	"flying": {
		"electric": 0.5, "grass": 2.0, "fighting": 2.0, "bug": 2.0, "rock": 0.5, "steel": 0.5,
	},
	"psychic": {
		"fighting": 2.0, "poison": 2.0, "psychic": 0.5, "dark": 0.0, "steel": 0.5,
	},
	"bug": {
		"fire": 0.5, "grass": 2.0, "fighting": 0.5, "poison": 0.5, "flying": 0.5, "psychic": 2.0, "ghost": 0.5, "dark": 2.0, "steel": 0.5, "fairy": 0.5,
	},
	"rock": {
		"fire": 2.0, "ice": 2.0, "fighting": 0.5, "ground": 0.5, "flying": 2.0, "bug": 2.0, "steel": 0.5,
	},
	"ghost": {
		"normal": 0.0, "psychic": 2.0, "ghost": 2.0, "dark": 0.5,
	},
	"dragon": {
		"dragon": 2.0, "steel": 0.5, "fairy": 0.0,
	},
	"dark": {
		"fighting": 0.5, "psychic": 2.0, "ghost": 2.0, "dark": 0.5, "fairy": 0.5,
	},
	"steel": {
		"fire": 0.5, "water": 0.5, "electric": 0.5, "ice": 2.0, "rock": 2.0, "steel": 0.5, "fairy": 2.0,
	},
	"fairy": {
		"fire": 0.5, "fighting": 2.0, "poison": 0.5, "dragon": 2.0, "dark": 2.0, "steel": 0.5,
	},
}


// Check if a Pokémon is legendary or mythical
func isLegendaryOrMythical(pokemonName string) bool {
	name := strings.ToLower(pokemonName)

	//how to simplify the use of loop using slices.Contains()
	//slices.Contains(legendaryNames, name)
	//slices.Contains(mythicalNames, name)
	
	if slices.Contains(legendaryNames, name) {
		return true
	}
	
	if slices.Contains(mythicalNames, name) {
		return true
	}
	
	return false
}

// TypeMultiplier returns the type effectiveness multiplier for a move type against defender types.
// Also applies 2x multiplier for legendary/mythical Pokémon attacks.
func TypeMultiplier(moveType string, defenderTypes []string, attackerName string) float64 {
	multiplier := 1.0
	
	// Apply type effectiveness
	moveType = strings.ToLower(moveType)
	if effectiveness, exists := typeChart[moveType]; exists {
		for _, defenderType := range defenderTypes {
			defenderType = strings.ToLower(defenderType)
			if typeMultiplier, typeExists := effectiveness[defenderType]; typeExists {
				multiplier *= typeMultiplier
			}
		}
	}
	
	// Apply legendary/mythical bonus (2x damage for their attacks)
	if isLegendaryOrMythical(attackerName) {
		multiplier *= 2.0
	}
	
	return multiplier
}



// Roll damage percent based on attack stat
func rollDamagePercent(attackStat int) float64 {
	// Three tables: low (<=30), high (70), super (>=120)
	low := []struct{ pct, prob float64 }{
		{0.10, 0.07}, {0.20, 0.13}, {0.30, 0.35}, {0.40, 0.25}, {0.60, 0.10}, {0.80, 0.07}, {1.00, 0.03},
	}
	high := []struct{ pct, prob float64 }{
		{0.10, 0.01}, {0.20, 0.04}, {0.30, 0.10}, {0.40, 0.15}, {0.60, 0.25}, {0.80, 0.25}, {1.00, 0.20},
	}
	super := []struct{ pct, prob float64 }{
		{0.10, 0.00}, {0.20, 0.01}, {0.30, 0.04}, {0.40, 0.10}, {0.60, 0.15}, {0.80, 0.30}, {1.00, 0.40},
	}

	var table []struct{ pct, prob float64 }
	if attackStat <= 30 {
		table = low
	} else if attackStat >= 120 {
		table = super
	} else if attackStat >= 70 {
		// Interpolate between high and super
		frac := float64(attackStat-70) / 50.0
		table = make([]struct{ pct, prob float64 }, len(high))
		for i := range high {
			table[i].pct = high[i].pct
			table[i].prob = high[i].prob*(1-frac) + super[i].prob*frac
		}
	} else {
		// Interpolate between low and high
		frac := float64(attackStat-30) / 40.0
		table = make([]struct{ pct, prob float64 }, len(low))
		for i := range low {
			table[i].pct = low[i].pct
			table[i].prob = low[i].prob*(1-frac) + high[i].prob*frac
		}
	}
	// Roll
	roll := rand.Float64()
	cum := 0.0
	for _, entry := range table {
		cum += entry.prob
		if roll < cum {
			return entry.pct
		}
	}
	return 0.10 // fallback
}

// FetchRandomDeck returns a slice of 5 random Pokémon cards, with legendary/mythical odds handled in FetchRandomPokemonCard
func FetchRandomDeck() []pokemon.Card {
	deck := make([]pokemon.Card, 0, 5)
	for range [5]int{} {
		card := pokemon.FetchRandomPokemonCard(false)
		deck = append(deck, card)
	}
	return deck
}

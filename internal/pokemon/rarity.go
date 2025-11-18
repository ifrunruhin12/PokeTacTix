package pokemon

import (
	"slices"
	"strings"
)

var legendaryNames = []string{
	"articuno", "zapdos", "moltres", "mewtwo", "raikou", "entei", "suicune", "lugia", "ho-oh",
	"regirock", "regice", "registeel", "latias", "latios", "kyogre", "groudon", "rayquaza",
	"uxie", "mesprit", "azelf", "dialga", "palkia", "heatran", "regigigas", "giratina", "cresselia",
	"cobalion", "terrakion", "virizion", "tornadus", "thundurus", "reshiram", "zekrom", "landorus", "kyurem",
	"xerneas", "yveltal", "zygarde", "tapu-koko", "tapu-lele", "tapu-bulu", "tapu-fini",
	"cosmog", "cosmoem", "solgaleo", "lunala", "necrozma", "zamazenta", "zacian", "eternatus",
	"kubfu", "urshifu", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus",
	"ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

var mythicalNames = []string{
	"mew", "celebi", "jirachi", "deoxys", "phione", "manaphy", "darkrai", "shaymin", "arceus",
	"victini", "keldeo", "meloetta", "genesect", "diancie", "hoopa", "volcanion",
	"magearna", "marshadow", "zeraora", "meltan", "melmetal", "zarude",
}

// IsLegendaryOrMythical checks if a Pokemon is legendary or mythical
func IsLegendaryOrMythical(name string) (isLegendary bool, isMythical bool) {
	nameLower := strings.ToLower(name)

	// Check if mythical
	if slices.Contains(mythicalNames, nameLower) {
		return false, true
	}

	// Check if legendary
	if slices.Contains(legendaryNames, nameLower) {
		return true, false
	}

	return false, false
}

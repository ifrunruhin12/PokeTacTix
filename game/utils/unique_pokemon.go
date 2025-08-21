package utils

import (
	"slices"
	"strings"
)

// Legendary and mythical Pokémon names (same as in fetch.go)
var legendaryNames = []string{
	"articuno", "zapdos", "moltres", "mewtwo", "raikou", "entei", "suicune", "lugia", "ho-oh", "regirock", "regice", "registeel", "latias", "latios", "kyogre", "groudon", "rayquaza", "uxie", "mesprit", "azelf", "dialga", "palkia", "heatran", "regigigas", "giratina", "cresselia", "cobalion", "terrakion", "virizion", "tornadus", "thundurus", "reshiram", "zekrom", "landorus", "kyurem", "xerneas", "yveltal", "zygarde", "tapu-koko", "tapu-lele", "tapu-bulu", "tapu-fini", "cosmog", "cosmoem", "solgaleo", "lunala", "necrozma", "zamazenta", "zacian", "eternatus", "kubfu", "urshifu", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus", "ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

var mythicalNames = []string{
	"mew", "celebi", "jirachi", "deoxys", "phione", "manaphy", "darkrai", "shaymin", "arceus", "victini", "keldeo", "meloetta", "genesect", "diancie", "hoopa", "volcanion", "magearna", "marshadow", "zeraora", "meltan", "melmetal", "zarude", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus", "ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

// IsLegendaryOrMythical checks if a Pokémon is legendary or mythical
func IsLegendaryOrMythical(pokemonName string) bool {
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

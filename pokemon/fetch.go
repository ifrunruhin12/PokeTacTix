package pokemon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

var legendaryNames = []string{
	"articuno", "zapdos", "moltres", "mewtwo", "raikou", "entei", "suicune", "lugia", "ho-oh", "regirock", "regice", "registeel", "latias", "latios", "kyogre", "groudon", "rayquaza", "uxie", "mesprit", "azelf", "dialga", "palkia", "heatran", "regigigas", "giratina", "cresselia", "cobalion", "terrakion", "virizion", "tornadus", "thundurus", "reshiram", "zekrom", "landorus", "kyurem", "xerneas", "yveltal", "zygarde", "tapu-koko", "tapu-lele", "tapu-bulu", "tapu-fini", "cosmog", "cosmoem", "solgaleo", "lunala", "necrozma", "zamazenta", "zacian", "eternatus", "kubfu", "urshifu", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus", "ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

var mythicalNames = []string{
	"mew", "celebi", "jirachi", "deoxys", "phione", "manaphy", "darkrai", "shaymin", "arceus", "victini", "keldeo", "meloetta", "genesect", "diancie", "hoopa", "volcanion", "magearna", "marshadow", "zeraora", "meltan", "melmetal", "zarude", "regieleki", "regidrago", "glastrier", "spectrier", "calyrex", "enamorus", "ting-lu", "chien-pao", "wo-chien", "chi-yu", "koraidon", "miraidon", "ogerpon",
}

// IsLegendaryOrMythical checks if a Pokemon is legendary or mythical
func IsLegendaryOrMythical(name string) (isLegendary bool, isMythical bool) {
	nameLower := strings.ToLower(name)
	
	// Check if mythical
	for _, mythical := range mythicalNames {
		if nameLower == mythical {
			return false, true
		}
	}
	
	// Check if legendary
	for _, legendary := range legendaryNames {
		if nameLower == legendary {
			return true, false
		}
	}
	
	return false, false
}

func GetMoves(rawMoves []RawMove) []Move {
	const maxMoves = 4
	perm := rand.Perm(len(rawMoves))
	var gameMoves []Move

	for _, i := range perm {
		moveURL := rawMoves[i].Move.URL
		resp, err := http.Get(moveURL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var data struct {
			Name  string `json:"name"`
			Power int    `json:"power"`
			Type  struct {
				Name string `json:"name"`
			} `json:"type"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			continue
		}

		if data.Power <= 0 {
			continue
		}

		gameMoves = append(gameMoves, Move{
			Name:        data.Name,
			Power:       data.Power,
			StaminaCost: data.Power / 3, // will go to game logic later
			Type:        data.Type.Name,
		})

		if len(gameMoves) == maxMoves {
			break
		}
	}

	return gameMoves
}

func FetchPokemon(name string) (Pokemon, []Move, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(name)
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, nil, fmt.Errorf("Pokemon \"%s\" not found. Please check the name and try again", name)
	}

	defer resp.Body.Close()

	var poke Pokemon

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&poke); err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to decode pokemon data: %w", err)
	}

	pokeMoves := GetMoves(poke.Moves)
	return poke, pokeMoves, nil
}

// FetchRandomPokemonCard returns a random Card. There is a 0.01% chance for a mythical, 0.01% for a legendary, otherwise normal.
func FetchRandomPokemonCard(_ bool) Card {
	mythicalOdds := 0.0001  // 0.01%
	legendaryOdds := 0.0001 // 0.01%
	maxRetries := 3
	for range maxRetries {
		roll := rand.Float64()
		var name string
		if roll < mythicalOdds {
			// Mythical
			name = mythicalNames[rand.Intn(len(mythicalNames))]
		} else if roll < mythicalOdds+legendaryOdds {
			// Legendary
			name = legendaryNames[rand.Intn(len(legendaryNames))]
		} else {
			// Normal
			id := rand.Intn(898) + 1 // Gen 1-8, skip forms
			name = fmt.Sprintf("%d", id)
		}
		poke, moves, err := FetchPokemon(name)
		if err != nil {
			continue
		}
		return BuildCardFromPokemon(poke, moves)
	}
	// Fallback dummy card
	return Card{Name: "MissingNo", HP: 33, Types: []string{"???"}}
}

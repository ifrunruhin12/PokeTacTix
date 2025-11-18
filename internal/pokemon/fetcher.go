package pokemon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"time"
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

// GetMoves fetches move details from the API
func GetMoves(rawMoves []RawMove) []Move {
	const maxMoves = 4
	perm := rand.Perm(len(rawMoves))
	var gameMoves []Move

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, i := range perm {
		moveURL := rawMoves[i].Move.URL
		resp, err := client.Get(moveURL)
		if err != nil {
			continue
		}

		var data struct {
			Name  string `json:"name"`
			Power int    `json:"power"`
			Type  struct {
				Name string `json:"name"`
			} `json:"type"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if data.Power <= 0 {
			continue
		}

		gameMoves = append(gameMoves, Move{
			Name:        data.Name,
			Power:       data.Power,
			StaminaCost: data.Power / 3,
			Type:        data.Type.Name,
		})

		if len(gameMoves) == maxMoves {
			break
		}
	}

	return gameMoves
}

// FetchPokemon fetches Pokemon data from the PokeAPI
func FetchPokemon(name string) (Pokemon, []Move, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(name)

	// Create HTTP client with timeout to prevent hanging
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
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

// FetchRandomPokemonCard returns a random Card
// There is a 0.01% chance for a mythical, 0.01% for a legendary, otherwise normal
func FetchRandomPokemonCard(_ bool) Card {
	mythicalOdds := 0.0001  // 0.01%
	legendaryOdds := 0.0001 // 0.01%
	maxRetries := 5         // Increased retries for better reliability

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
			// Normal - use Gen 1-5 for better reliability (fewer edge cases)
			id := rand.Intn(649) + 1 // Gen 1-5
			name = fmt.Sprintf("%d", id)
		}
		poke, moves, err := FetchPokemon(name)
		if err != nil {
			continue
		}
		card := BuildCardFromPokemon(poke, moves)
		// Ensure card has valid moves
		if len(card.Moves) == 0 {
			// Add a default move if none were found
			card.Moves = []Move{
				{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
			}
		}
		return card
	}

	// Fallback to a reliable starter Pokemon if all retries fail
	poke, moves, err := FetchPokemon("pikachu")
	if err == nil {
		card := BuildCardFromPokemon(poke, moves)
		if len(card.Moves) == 0 {
			card.Moves = []Move{
				{Name: "thunderbolt", Power: 90, StaminaCost: 30, Type: "electric"},
			}
		}
		return card
	}

	// Ultimate fallback - a valid dummy card with proper stats
	return Card{
		Name:    "Pikachu",
		HP:      100,
		HPMax:   100,
		Stamina: 100,
		Attack:  55,
		Defense: 40,
		Speed:   90,
		Types:   []string{"electric"},
		Moves: []Move{
			{Name: "thunderbolt", Power: 90, StaminaCost: 30, Type: "electric"},
			{Name: "quick-attack", Power: 40, StaminaCost: 13, Type: "normal"},
		},
		Sprite: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
	}
}

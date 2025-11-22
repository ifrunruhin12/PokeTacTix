package pokemon

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
)

// Embedded Pokemon data
//go:embed data/pokemon_data.json
var pokemonDataJSON []byte

// PokemonEntry represents a Pokemon in the offline database
type PokemonEntry struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	HP          int      `json:"hp"`
	Attack      int      `json:"attack"`
	Defense     int      `json:"defense"`
	Speed       int      `json:"speed"`
	Types       []string `json:"types"`
	Moves       []Move   `json:"moves"`
	Sprite      string   `json:"sprite"`
	IsLegendary bool     `json:"is_legendary"`
	IsMythical  bool     `json:"is_mythical"`
}

// PokemonDatabase holds all Pokemon data
type PokemonDatabase struct {
	Pokemon   []PokemonEntry `json:"pokemon"`
	Generated string         `json:"generated"`
	Version   string         `json:"version"`
	
	// Internal index for fast lookup
	pokemonByID map[int]*PokemonEntry
}

var (
	globalDatabase *PokemonDatabase
	loadOnce       sync.Once
	loadError      error
)

// LoadPokemonDatabase loads and parses the embedded Pokemon data
// Uses sync.Once to ensure it's only loaded once
func LoadPokemonDatabase() (*PokemonDatabase, error) {
	loadOnce.Do(func() {
		var db PokemonDatabase
		
		if err := json.Unmarshal(pokemonDataJSON, &db); err != nil {
			loadError = fmt.Errorf("failed to parse embedded Pokemon data: %w", err)
			return
		}
		
		// Build index for fast lookup
		db.pokemonByID = make(map[int]*PokemonEntry, len(db.Pokemon))
		for i := range db.Pokemon {
			db.pokemonByID[db.Pokemon[i].ID] = &db.Pokemon[i]
		}
		
		globalDatabase = &db
	})
	
	if loadError != nil {
		return nil, loadError
	}
	
	return globalDatabase, nil
}

// GetPokemonByID retrieves a Pokemon by its ID
func GetPokemonByID(id int) (*PokemonEntry, error) {
	db, err := LoadPokemonDatabase()
	if err != nil {
		return nil, err
	}
	
	pokemon, exists := db.pokemonByID[id]
	if !exists {
		return nil, fmt.Errorf("pokemon with ID %d not found", id)
	}
	
	return pokemon, nil
}

// GetRandomPokemon returns a random Pokemon from the database
// excludeLegendary: if true, excludes legendary Pokemon
// excludeMythical: if true, excludes mythical Pokemon
// Optimized to avoid copying entire Pokemon slice
func GetRandomPokemon(excludeLegendary, excludeMythical bool) (*PokemonEntry, error) {
	db, err := LoadPokemonDatabase()
	if err != nil {
		return nil, err
	}
	
	// Build index list instead of copying Pokemon entries (memory optimization)
	var candidateIndices []int
	for i := range db.Pokemon {
		pokemon := &db.Pokemon[i]
		if excludeLegendary && pokemon.IsLegendary {
			continue
		}
		if excludeMythical && pokemon.IsMythical {
			continue
		}
		candidateIndices = append(candidateIndices, i)
	}
	
	if len(candidateIndices) == 0 {
		return nil, fmt.Errorf("no Pokemon available with given filters")
	}
	
	// Select random Pokemon by index
	selectedIdx := candidateIndices[rand.Intn(len(candidateIndices))]
	return &db.Pokemon[selectedIdx], nil
}

// GetDatabaseStats returns statistics about the loaded database
func GetDatabaseStats() (total, legendary, mythical int, err error) {
	db, err := LoadPokemonDatabase()
	if err != nil {
		return 0, 0, 0, err
	}
	
	total = len(db.Pokemon)
	for _, pokemon := range db.Pokemon {
		if pokemon.IsLegendary {
			legendary++
		}
		if pokemon.IsMythical {
			mythical++
		}
	}
	
	return total, legendary, mythical, nil
}

// FetchRandomPokemonCardOffline returns a random Card using offline data
// Applies rarity logic: 0.01% mythical, 0.01% legendary, rest common/uncommon/rare
// No network calls - completely offline
func FetchRandomPokemonCardOffline() Card {
	mythicalOdds := 0.0001  // 0.01%
	legendaryOdds := 0.0001 // 0.01%
	
	roll := rand.Float64()
	var pokemon *PokemonEntry
	var err error
	
	if roll < mythicalOdds {
		// Try to get a mythical Pokemon
		pokemon, err = GetRandomPokemon(true, false) // exclude legendary, include mythical
		if err != nil {
			// Fallback to any Pokemon if no mythical available
			pokemon, _ = GetRandomPokemon(false, false)
		}
	} else if roll < mythicalOdds+legendaryOdds {
		// Try to get a legendary Pokemon
		pokemon, err = GetRandomPokemon(false, true) // include legendary, exclude mythical
		if err != nil {
			// Fallback to any Pokemon if no legendary available
			pokemon, _ = GetRandomPokemon(false, false)
		}
	} else {
		// Get a normal Pokemon (exclude legendary and mythical)
		pokemon, err = GetRandomPokemon(true, true)
		if err != nil {
			// Fallback to any Pokemon
			pokemon, _ = GetRandomPokemon(false, false)
		}
	}
	
	// Ultimate fallback to Pikachu if database is corrupted
	if pokemon == nil {
		return Card{
			Name:        "Pikachu",
			HP:          100,
			HPMax:       100,
			Stamina:     180,
			Attack:      55,
			Defense:     40,
			Speed:       90,
			Types:       []string{"electric"},
			Moves: []Move{
				{Name: "thunderbolt", Power: 90, StaminaCost: 30, Type: "electric"},
				{Name: "quick-attack", Power: 40, StaminaCost: 13, Type: "normal"},
				{Name: "iron-tail", Power: 100, StaminaCost: 33, Type: "steel"},
				{Name: "electro-ball", Power: 80, StaminaCost: 26, Type: "electric"},
			},
			Sprite:      "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png",
			Level:       1,
			XP:          0,
			IsLegendary: false,
			IsMythical:  false,
		}
	}
	
	// Build Card from PokemonEntry
	return buildCardFromEntry(pokemon)
}

// buildCardFromEntry converts a PokemonEntry to a Card
func buildCardFromEntry(entry *PokemonEntry) Card {
	stamina := entry.Speed * 2
	
	// Ensure moves are properly formatted
	moves := make([]Move, len(entry.Moves))
	copy(moves, entry.Moves)
	
	// Ensure at least one move
	if len(moves) == 0 {
		moves = []Move{
			{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
		}
	}
	
	return Card{
		Name:        entry.Name,
		HP:          entry.HP,
		HPMax:       entry.HP,
		Stamina:     stamina,
		Attack:      entry.Attack,
		Defense:     entry.Defense,
		Speed:       entry.Speed,
		Types:       entry.Types,
		Moves:       moves,
		Sprite:      entry.Sprite,
		Level:       1,
		XP:          0,
		IsLegendary: entry.IsLegendary,
		IsMythical:  entry.IsMythical,
	}
}

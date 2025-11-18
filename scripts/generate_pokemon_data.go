package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

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

// Move represents a Pokemon move
type Move struct {
	Name        string `json:"name"`
	Power       int    `json:"power"`
	StaminaCost int    `json:"stamina_cost"`
	Type        string `json:"attack_type"`
}

// PokemonDatabase holds all Pokemon data
type PokemonDatabase struct {
	Pokemon   []PokemonEntry `json:"pokemon"`
	Generated string         `json:"generated"`
	Version   string         `json:"version"`
}

var legendaryNames = []string{
	"articuno", "zapdos", "moltres", "mewtwo", "raikou", "entei", "suicune", "lugia", "ho-oh",
	"regirock", "regice", "registeel", "latias", "latios", "kyogre", "groudon", "rayquaza",
	"uxie", "mesprit", "azelf", "dialga", "palkia", "heatran", "regigigas", "giratina", "cresselia",
	"cobalion", "terrakion", "virizion", "tornadus", "thundurus", "reshiram", "zekrom", "landorus", "kyurem",
}

var mythicalNames = []string{
	"mew", "celebi", "jirachi", "deoxys", "phione", "manaphy", "darkrai", "shaymin", "arceus",
	"victini", "keldeo", "meloetta", "genesect",
}

func isLegendaryOrMythical(name string) (isLegendary bool, isMythical bool) {
	nameLower := strings.ToLower(name)
	if slices.Contains(mythicalNames, nameLower) {
		return false, true
	}
	if slices.Contains(legendaryNames, nameLower) {
		return true, false
	}
	return false, false
}

func fetchPokemonData(id int, client *http.Client) (*PokemonEntry, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)
	
	var resp *http.Response
	var err error
	
	// Retry logic with exponential backoff
	for attempt := 0; attempt < 3; attempt++ {
		resp, err = client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		
		// Wait before retry (exponential backoff)
		if attempt < 2 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon %d: %w", id, err)
	}
	
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("pokemon %d not found (status %d)", id, resp.StatusCode)
	}
	
	defer resp.Body.Close()
	
	var data struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Stats   []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
		Moves []struct {
			Move struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move"`
		} `json:"moves"`
		Sprites struct {
			FrontDefault string `json:"front_default"`
		} `json:"sprites"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode pokemon %d: %w", id, err)
	}
	
	// Extract stats
	var hp, attack, defense, speed int
	for _, stat := range data.Stats {
		switch stat.Stat.Name {
		case "hp":
			hp = stat.BaseStat
		case "attack":
			attack = stat.BaseStat
		case "defense":
			defense = stat.BaseStat
		case "speed":
			speed = stat.BaseStat
		}
	}
	
	// Adjust HP (same as web version)
	hp = hp + int(float64(hp)*0.5)
	
	// Extract types
	types := make([]string, len(data.Types))
	for i, t := range data.Types {
		types[i] = t.Type.Name
	}
	
	// Fetch moves (limit to 4)
	moves := fetchMoves(data.Moves, client)
	
	// Check legendary/mythical status
	isLegendary, isMythical := isLegendaryOrMythical(data.Name)
	
	return &PokemonEntry{
		ID:          data.ID,
		Name:        data.Name,
		HP:          hp,
		Attack:      attack,
		Defense:     defense,
		Speed:       speed,
		Types:       types,
		Moves:       moves,
		Sprite:      data.Sprites.FrontDefault,
		IsLegendary: isLegendary,
		IsMythical:  isMythical,
	}, nil
}

func fetchMoves(rawMoves []struct {
	Move struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"move"`
}, client *http.Client) []Move {
	const maxMoves = 4
	
	if len(rawMoves) == 0 {
		// Fallback move
		return []Move{
			{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
		}
	}
	
	// Shuffle moves to get random selection
	perm := rand.Perm(len(rawMoves))
	var moves []Move
	
	for _, i := range perm {
		if len(moves) >= maxMoves {
			break
		}
		
		moveURL := rawMoves[i].Move.URL
		resp, err := client.Get(moveURL)
		if err != nil {
			continue
		}
		
		var moveData struct {
			Name  string `json:"name"`
			Power int    `json:"power"`
			Type  struct {
				Name string `json:"name"`
			} `json:"type"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&moveData); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		
		// Only include moves with power
		if moveData.Power <= 0 {
			continue
		}
		
		moves = append(moves, Move{
			Name:        moveData.Name,
			Power:       moveData.Power,
			StaminaCost: moveData.Power / 3,
			Type:        moveData.Type.Name,
		})
	}
	
	// Ensure at least one move
	if len(moves) == 0 {
		moves = []Move{
			{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
		}
	}
	
	return moves
}

func main() {
	log.Println("Starting Pokemon data generation...")
	log.Println("Fetching Gen 1-5 Pokemon (IDs 1-649)")
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	database := PokemonDatabase{
		Pokemon:   make([]PokemonEntry, 0, 649),
		Generated: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}
	
	successCount := 0
	failCount := 0
	
	// Fetch Pokemon 1-649 (Gen 1-5)
	for id := 1; id <= 649; id++ {
		log.Printf("Fetching Pokemon %d/649...", id)
		
		pokemon, err := fetchPokemonData(id, client)
		if err != nil {
			log.Printf("Warning: %v", err)
			failCount++
			
			// Create fallback entry for failed fetches
			pokemon = &PokemonEntry{
				ID:      id,
				Name:    fmt.Sprintf("pokemon-%d", id),
				HP:      100,
				Attack:  50,
				Defense: 50,
				Speed:   50,
				Types:   []string{"normal"},
				Moves: []Move{
					{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
				},
				Sprite:      fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%d.png", id),
				IsLegendary: false,
				IsMythical:  false,
			}
		} else {
			successCount++
		}
		
		database.Pokemon = append(database.Pokemon, *pokemon)
		
		// Rate limiting: 1 request per 100ms
		time.Sleep(100 * time.Millisecond)
		
		// Progress update every 50 Pokemon
		if id%50 == 0 {
			log.Printf("Progress: %d/649 (Success: %d, Failed: %d)", id, successCount, failCount)
		}
	}
	
	log.Printf("Fetch complete! Success: %d, Failed: %d", successCount, failCount)
	
	// Create output directory
	outputDir := "internal/pokemon/data"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Write to file
	outputPath := filepath.Join(outputDir, "pokemon_data.json")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(database); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}
	
	// Get file size
	fileInfo, _ := file.Stat()
	sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	
	log.Printf("Successfully generated %s (%.2f MB)", outputPath, sizeMB)
	log.Printf("Total Pokemon: %d", len(database.Pokemon))
	log.Println("Generation complete!")
}

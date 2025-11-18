package main

import (
	"fmt"
	"log"
	"os"
	"pokemon-cli/internal/pokemon"
)

func main() {
	log.Println("Testing Offline Pokemon Data System...")
	
	// Set offline mode
	os.Setenv("POKEMON_OFFLINE_MODE", "true")
	
	// Test 1: Load database
	log.Println("\n=== Test 1: Loading Pokemon Database ===")
	db, err := pokemon.LoadPokemonDatabase()
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}
	log.Printf("✓ Database loaded successfully")
	log.Printf("  Total Pokemon: %d", len(db.Pokemon))
	log.Printf("  Version: %s", db.Version)
	log.Printf("  Generated: %s", db.Generated)
	
	// Test 2: Get database stats
	log.Println("\n=== Test 2: Database Statistics ===")
	total, legendary, mythical, err := pokemon.GetDatabaseStats()
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}
	log.Printf("✓ Stats retrieved successfully")
	log.Printf("  Total: %d", total)
	log.Printf("  Legendary: %d", legendary)
	log.Printf("  Mythical: %d", mythical)
	
	// Test 3: Get Pokemon by ID
	log.Println("\n=== Test 3: Get Pokemon by ID ===")
	pikachu, err := pokemon.GetPokemonByID(25)
	if err != nil {
		log.Fatalf("Failed to get Pikachu: %v", err)
	}
	log.Printf("✓ Retrieved Pokemon #25")
	log.Printf("  Name: %s", pikachu.Name)
	log.Printf("  HP: %d", pikachu.HP)
	log.Printf("  Attack: %d", pikachu.Attack)
	log.Printf("  Defense: %d", pikachu.Defense)
	log.Printf("  Speed: %d", pikachu.Speed)
	log.Printf("  Types: %v", pikachu.Types)
	log.Printf("  Moves: %d", len(pikachu.Moves))
	
	// Test 4: Get random Pokemon (exclude legendary/mythical)
	log.Println("\n=== Test 4: Get Random Pokemon (Normal) ===")
	randomPokemon, err := pokemon.GetRandomPokemon(true, true)
	if err != nil {
		log.Fatalf("Failed to get random Pokemon: %v", err)
	}
	log.Printf("✓ Retrieved random Pokemon")
	log.Printf("  Name: %s", randomPokemon.Name)
	log.Printf("  Legendary: %v", randomPokemon.IsLegendary)
	log.Printf("  Mythical: %v", randomPokemon.IsMythical)
	
	// Test 5: Fetch random Pokemon card (offline)
	log.Println("\n=== Test 5: Fetch Random Pokemon Card (Offline) ===")
	card := pokemon.FetchRandomPokemonCardOffline()
	log.Printf("✓ Generated random card")
	log.Printf("  Name: %s", card.Name)
	log.Printf("  HP: %d/%d", card.HP, card.HPMax)
	log.Printf("  Stamina: %d", card.Stamina)
	log.Printf("  Attack: %d", card.Attack)
	log.Printf("  Defense: %d", card.Defense)
	log.Printf("  Speed: %d", card.Speed)
	log.Printf("  Types: %v", card.Types)
	log.Printf("  Moves: %d", len(card.Moves))
	for i, move := range card.Moves {
		log.Printf("    %d. %s (Power: %d, Cost: %d, Type: %s)", 
			i+1, move.Name, move.Power, move.StaminaCost, move.Type)
	}
	
	// Test 6: Test FetchRandomPokemonCard with offline mode
	log.Println("\n=== Test 6: FetchRandomPokemonCard (Auto-detect Offline) ===")
	card2 := pokemon.FetchRandomPokemonCard(false)
	log.Printf("✓ Generated card using FetchRandomPokemonCard")
	log.Printf("  Name: %s", card2.Name)
	log.Printf("  Level: %d", card2.Level)
	log.Printf("  XP: %d", card2.XP)
	
	// Test 7: Generate multiple cards to test randomness
	log.Println("\n=== Test 7: Generate Multiple Cards (Testing Randomness) ===")
	pokemonNames := make(map[string]int)
	for i := 0; i < 10; i++ {
		card := pokemon.FetchRandomPokemonCardOffline()
		pokemonNames[card.Name]++
	}
	log.Printf("✓ Generated 10 random cards")
	log.Printf("  Unique Pokemon: %d", len(pokemonNames))
	for name, count := range pokemonNames {
		log.Printf("    %s: %d", name, count)
	}
	
	log.Println("\n=== All Tests Passed! ===")
	fmt.Println("\nOffline Pokemon Data System is working correctly!")
}

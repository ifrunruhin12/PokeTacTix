package setup

import (
	"fmt"
	"time"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
	"pokemon-cli/internal/pokemon"
)

// GenerateStarterDeck creates 5 random non-legendary, non-mythical Pokemon
// Returns a slice of PlayerCard with level 1 and 0 XP
func GenerateStarterDeck() ([]storage.PlayerCard, error) {
	starterCards := make([]storage.PlayerCard, 0, 5)
	usedIDs := make(map[int]bool) // Track used Pokemon IDs to avoid duplicates

	for i := 0; i < 5; i++ {
		// Keep trying until we get a unique Pokemon
		var pokemonEntry *pokemon.PokemonEntry
		var err error

		for {
			// Get random non-legendary, non-mythical Pokemon
			pokemonEntry, err = pokemon.GetRandomPokemon(true, true)
			if err != nil {
				return nil, fmt.Errorf("failed to generate starter Pokemon: %w", err)
			}

			// Check if we already have this Pokemon
			if !usedIDs[pokemonEntry.ID] {
				usedIDs[pokemonEntry.ID] = true
				break
			}
		}

		// Create PlayerCard from PokemonEntry
		card := storage.PlayerCard{
			ID:          i, // Temporary ID, will be reassigned when added to collection
			PokemonID:   pokemonEntry.ID,
			Name:        pokemonEntry.Name,
			Level:       1,
			XP:          0,
			BaseHP:      pokemonEntry.HP,
			BaseAttack:  pokemonEntry.Attack,
			BaseDefense: pokemonEntry.Defense,
			BaseSpeed:   pokemonEntry.Speed,
			Types:       pokemonEntry.Types,
			Moves:       pokemonEntry.Moves,
			Sprite:      pokemonEntry.Sprite,
			IsLegendary: pokemonEntry.IsLegendary,
			IsMythical:  pokemonEntry.IsMythical,
			AcquiredAt:  time.Now(),
		}

		starterCards = append(starterCards, card)
	}

	return starterCards, nil
}

// DisplayStarterPokemon shows the starter Pokemon with their details
func DisplayStarterPokemon(starterCards []storage.PlayerCard) {
	fmt.Println()
	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("=== YOUR STARTER POKEMON ===", ui.ColorBrightYellow))
	} else {
		fmt.Println("=== YOUR STARTER POKEMON ===")
	}
	fmt.Println()

	fmt.Println("You've been given 5 Pokemon to start your journey:")
	fmt.Println()

	for i, card := range starterCards {
		stats := card.GetCurrentStats()

		// Pokemon number and name
		nameDisplay := fmt.Sprintf("%d. %s", i+1, card.Name)
		if ui.GetColorSupport() {
			nameDisplay = ui.Colorize(nameDisplay, ui.ColorBrightCyan)
		}
		fmt.Println(nameDisplay)

		// Types
		typeDisplay := "   Types: "
		for j, t := range card.Types {
			if j > 0 {
				typeDisplay += ", "
			}
			if ui.GetColorSupport() {
				typeDisplay += ui.ColorizeType(t, t)
			} else {
				typeDisplay += t
			}
		}
		fmt.Println(typeDisplay)

		// Stats
		fmt.Printf("   Stats: HP: %d | Attack: %d | Defense: %d | Speed: %d | Stamina: %d\n",
			stats.HP, stats.Attack, stats.Defense, stats.Speed, stats.Stamina)

		// Moves
		fmt.Println("   Moves:")
		for _, move := range card.Moves {
			moveType := move.Type
			if ui.GetColorSupport() {
				moveType = ui.ColorizeType(move.Type, move.Type)
			}
			fmt.Printf("     • %s (%s) - Power: %d, Stamina: %d\n",
				move.Name, moveType, move.Power, move.StaminaCost)
		}

		fmt.Println()
	}

	fmt.Println(ui.RenderDivider(75, "─"))
	fmt.Println()

	if ui.GetColorSupport() {
		fmt.Println(ui.Colorize("These Pokemon will form your starting deck!", ui.ColorBrightGreen))
	} else {
		fmt.Println("These Pokemon will form your starting deck!")
	}
	fmt.Println("You can customize your deck later using the 'deck' command.")
	fmt.Println()
}

// ConfirmStarterDeck prompts the user to confirm they're ready to start
func ConfirmStarterDeck() error {
	fmt.Print("Press Enter to continue...")
	fmt.Scanln()
	return nil
}

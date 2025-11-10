package cards

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"pokemon-cli/internal/database"
	"pokemon-cli/internal/pokemon"
	"time"
)

// Service handles business logic for Pokemon cards
type Service struct {
	repository *Repository
}

// NewService creates a new card service
func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

// GenerateStarterDeck generates 5 random non-legendary Pokemon cards for a new user
// Requirements: 3.1, 3.2, 3.3, 3.4
func (s *Service) GenerateStarterDeck(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	const deckSize = 5
	const maxRetries = 50
	
	starterCards := make([]database.PlayerCard, 0, deckSize)
	usedNames := make(map[string]bool)
	
	// Generate 5 unique non-legendary, non-mythical Pokemon
	for len(starterCards) < deckSize && maxRetries > len(starterCards)*10 {
		// Generate random Pokemon ID (Gen 1-8, avoiding forms)
		pokemonID := rand.Intn(898) + 1
		pokemonName := fmt.Sprintf("%d", pokemonID)
		
		// Fetch Pokemon data
		poke, moves, err := pokemon.FetchPokemon(pokemonName)
		if err != nil {
			continue
		}
		
		// Check if already used
		if usedNames[poke.Name] {
			continue
		}
		
		// Check if legendary or mythical
		isLegendary, isMythical := pokemon.IsLegendaryOrMythical(poke.Name)
		if isLegendary || isMythical {
			continue
		}
		
		// Build card
		card := pokemon.BuildCardFromPokemon(poke, moves)
		
		// Convert to PlayerCard
		typesJSON, err := json.Marshal(card.Types)
		if err != nil {
			continue
		}
		
		movesJSON, err := json.Marshal(card.Moves)
		if err != nil {
			continue
		}
		
		deckPosition := len(starterCards) + 1
		playerCard := &database.PlayerCard{
			UserID:       userID,
			PokemonName:  card.Name,
			Level:        1,
			XP:           0,
			BaseHP:       card.HPMax,
			BaseAttack:   card.Attack,
			BaseDefense:  card.Defense,
			BaseSpeed:    card.Speed,
			Types:        typesJSON,
			Moves:        movesJSON,
			Sprite:       card.Sprite,
			IsLegendary:  false,
			IsMythical:   false,
			InDeck:       true,
			DeckPosition: &deckPosition,
		}
		
		// Create card in database
		createdCard, err := s.repository.Create(ctx, playerCard)
		if err != nil {
			return nil, fmt.Errorf("failed to create starter card: %w", err)
		}
		
		starterCards = append(starterCards, *createdCard)
		usedNames[poke.Name] = true
	}
	
	if len(starterCards) < deckSize {
		return nil, fmt.Errorf("failed to generate %d starter cards, only got %d", deckSize, len(starterCards))
	}
	
	return starterCards, nil
}

// AddXP adds XP to a card and handles level-ups
// Requirements: 18.1, 18.2, 18.3, 18.4, 18.5
func (s *Service) AddXP(ctx context.Context, cardID int, xp int) (*database.PlayerCard, error) {
	return s.repository.AddXP(ctx, cardID, xp)
}

// GetUserCards retrieves all cards for a user
func (s *Service) GetUserCards(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	return s.repository.GetUserCards(ctx, userID)
}

// GetUserDeck retrieves the user's current deck
func (s *Service) GetUserDeck(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	return s.repository.GetUserDeck(ctx, userID)
}

// UpdateDeck updates the user's deck configuration
func (s *Service) UpdateDeck(ctx context.Context, userID int, cardIDs []int) error {
	return s.repository.UpdateDeck(ctx, userID, cardIDs)
}

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

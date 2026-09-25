package cards

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	"pokemon-cli/internal/database"
	"pokemon-cli/internal/items"
	"pokemon-cli/internal/middleware"
	"pokemon-cli/internal/pokemon"
)

var (
	// ErrCardNotFound is returned when the card does not exist
	ErrCardNotFound = errors.New("card not found")

	// ErrNotCardOwner is returned when the card belongs to another player
	ErrNotCardOwner = errors.New("card does not belong to the requesting user")

	// ErrNotEligible is returned when the card has no evolution matching the
	// requested method/item
	ErrNotEligible = errors.New("pokemon is not eligible for this evolution")

	// ErrConcurrentEvolution is returned when the card changed between the
	// evolution pre-check and the transaction (e.g. a duplicate request raced)
	ErrConcurrentEvolution = errors.New("card was modified concurrently, evolution aborted")
)

// cardRepository is the persistence surface the card service depends on.
// *Repository implements it; tests provide their own.
type cardRepository interface {
	Create(ctx context.Context, card *database.PlayerCard) (*database.PlayerCard, error)
	GetByID(ctx context.Context, id int) (*database.PlayerCard, error)
	GetUserCards(ctx context.Context, userID int) ([]database.PlayerCard, error)
	GetUserDeck(ctx context.Context, userID int) ([]database.PlayerCard, error)
	UpdateDeck(ctx context.Context, userID int, cardIDs []int) error
	AddXP(ctx context.Context, cardID int, xp int) (*database.PlayerCard, error)
	GetPokemonID(ctx context.Context, cardID, userID int) (*int, error)
	SetPokemonID(ctx context.Context, cardID, userID, pokemonID int) error
	EvolveWithItem(ctx context.Context, userID, cardID, expectedPokemonID int, target *pokemon.Pokemon, itemID string) (int, error)
}

// EvolutionRules is the evolution-domain surface the card service needs. It is
// implemented by the concrete pokemon service.
type EvolutionRules interface {
	GetByName(ctx context.Context, name string) (*pokemon.Pokemon, error)
	// GetEvolutionForItem returns the Pokemon the given pokemon (by pokemon ID)
	// evolves into when the given item id (slug) is used, or nil if no
	// use-item evolution matches.
	GetEvolutionForItem(ctx context.Context, pokemonID int, itemID string) (*pokemon.Pokemon, error)
	// GetEvolutionOptions returns every evolution mechanism available to the
	// given pokemon (by pokemon ID).
	GetEvolutionOptions(ctx context.Context, pokemonID int) ([]pokemon.EvolutionOption, error)
}

// ItemCatalog is the item-domain surface the card service needs. It is
// implemented by the items service.
type ItemCatalog interface {
	GetItem(ctx context.Context, itemID string) (*items.Item, error)
	GetQuantity(ctx context.Context, userID int, itemID string) (int, error)
}

// Service handles business logic for Pokemon cards
type Service struct {
	repository cardRepository
	evolutions EvolutionRules
	itemShop   ItemCatalog
}

// NewService creates a new card service
func NewService(repository cardRepository, evolutions EvolutionRules, itemShop ItemCatalog) *Service {
	return &Service{
		repository: repository,
		evolutions: evolutions,
		itemShop:   itemShop,
	}
}

func (s *Service) GenerateStarterDeck(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	const deckSize = 5
	const maxRetries = 50

	starterCards := make([]database.PlayerCard, 0, deckSize)
	usedNames := make(map[string]bool)

	for len(starterCards) < deckSize && maxRetries > len(starterCards)*10 {
		pokemonID := rand.Intn(898) + 1
		pokemonName := fmt.Sprintf("%d", pokemonID)

		poke, moves, err := pokemon.FetchPokemon(pokemonName)
		if err != nil {
			continue
		}

		if usedNames[poke.Name] {
			continue
		}

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

// GetEvolutionInfo describes how a card can evolve: every available mechanism
// (level, item), the required level or item, and whether the player can act on
// it right now. Level-up evolutions apply automatically after battles; only
// item evolutions are player-initiated.
func (s *Service) GetEvolutionInfo(ctx context.Context, userID, cardID int) (*EvolutionInfo, error) {
	card, err := s.getOwnedCard(ctx, userID, cardID)
	if err != nil {
		return nil, err
	}

	pokemonID, err := s.resolvePokemonID(ctx, card)
	if err != nil {
		return nil, err
	}

	options, err := s.evolutions.GetEvolutionOptions(ctx, pokemonID)
	if err != nil {
		return nil, fmt.Errorf("failed to load evolution options for card %d: %w", cardID, err)
	}

	// Options must serialize as [] (never null) so clients can render an
	// empty list for final-form Pokemon without crashing.
	info := &EvolutionInfo{
		CardID:      card.ID,
		PokemonName: card.PokemonName,
		Level:       card.Level,
		Options:     []EvolutionOptionResponse{},
	}
	for _, opt := range options {
		resp := EvolutionOptionResponse{
			Method:       opt.Method,
			TargetName:   opt.TargetName,
			TargetSprite: opt.TargetSprite,
			MinLevel:     opt.MinLevel,
			RequiredItem: opt.Item,
		}

		switch opt.Method {
		case pokemon.MethodLevel:
			// Level evolutions are applied automatically when the level is
			// reached through battle rewards.
			resp.Eligible = card.Level >= opt.MinLevel
		case pokemon.MethodItem:
			if s.itemShop != nil {
				owned, err := s.itemShop.GetQuantity(ctx, userID, opt.Item)
				if err != nil {
					return nil, fmt.Errorf("failed to load inventory for item %s: %w", opt.Item, err)
				}
				resp.OwnedQuantity = owned
				resp.Eligible = owned > 0
				if item, err := s.itemShop.GetItem(ctx, opt.Item); err == nil {
					resp.RequiredItemName = item.Name
				}
			}
			if resp.Eligible {
				info.CanEvolveNow = true
			}
		}

		info.Options = append(info.Options, resp)
	}

	return info, nil
}

// GetEvolutionSummaries reports, for every card the user owns, whether its
// species has any evolution path. The deck UI uses this to hide the Evolve
// action on final-form Pokemon. Cards whose species cannot be resolved are
// skipped (treated as "unknown" by clients) rather than failing the listing.
func (s *Service) GetEvolutionSummaries(ctx context.Context, userID int) ([]EvolutionSummary, error) {
	cards, err := s.repository.GetUserCards(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cards for user %d: %w", userID, err)
	}

	summaries := make([]EvolutionSummary, 0, len(cards))
	for i := range cards {
		card := &cards[i]

		pokemonID, err := s.resolvePokemonID(ctx, card)
		if err != nil {
			slog.Warn("skipping evolution summary for unresolvable card", "card_id", card.ID, "error", err)
			continue
		}

		options, err := s.evolutions.GetEvolutionOptions(ctx, pokemonID)
		if err != nil {
			slog.Warn("skipping evolution summary on options error", "card_id", card.ID, "error", err)
			continue
		}

		summaries = append(summaries, EvolutionSummary{
			CardID:       card.ID,
			HasEvolution: len(options) > 0,
		})
	}
	return summaries, nil
}

// EvolveCardWithItem evolves a card using the given evolution item. The item
// is consumed from the player's inventory in the same transaction that applies
// the evolution; a failed evolution never consumes the item.
func (s *Service) EvolveCardWithItem(ctx context.Context, userID, cardID int, itemID string) (*EvolutionResult, error) {
	card, err := s.getOwnedCard(ctx, userID, cardID)
	if err != nil {
		return nil, err
	}

	if s.itemShop == nil {
		return nil, fmt.Errorf("item system is not configured")
	}

	// Reject unknown items before touching evolution rules or inventory.
	item, err := s.itemShop.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// The requested item must match the card's item-based evolution rule.
	// GetEvolutionForItem is data-driven from the evolution chain: it returns
	// nil for any item this Pokemon does not evolve with.
	pokemonID, err := s.resolvePokemonID(ctx, card)
	if err != nil {
		return nil, err
	}

	target, err := s.evolutions.GetEvolutionForItem(ctx, pokemonID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve item evolution for card %d: %w", cardID, err)
	}
	if target == nil {
		return nil, fmt.Errorf("%w: %s does not evolve with %s", ErrNotEligible, card.PokemonName, item.Name)
	}

	// Pre-check inventory for a clear error; the transaction below re-checks
	// with a guarded UPDATE so a race cannot over-consume.
	owned, err := s.itemShop.GetQuantity(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory for item %s: %w", itemID, err)
	}
	if owned < 1 {
		return nil, fmt.Errorf("%w: you own 0 of %s", items.ErrInsufficientItem, item.Name)
	}

	// Snapshot the species name before the transaction; the card row changes
	// inside it.
	fromName := card.PokemonName

	remaining, err := s.repository.EvolveWithItem(ctx, userID, cardID, pokemonID, target, itemID)
	if err != nil {
		return nil, err
	}

	updated, err := s.repository.GetByID(ctx, cardID)
	if err != nil {
		// Evolution committed; the re-read failing is not fatal to the result.
		slog.Warn("failed to reload card after evolution", "card_id", cardID, "error", err)
		updated = nil
	}

	middleware.EvolutionTotal.Inc()
	slog.Info("pokemon evolved via item",
		"user_id", userID, "card_id", cardID,
		"from", fromName, "into", target.Name, "item", itemID,
	)

	return &EvolutionResult{
		CardID:                cardID,
		EvolvedFrom:           fromName,
		EvolvedInto:           target.Name,
		ItemConsumed:          itemID,
		RemainingItemQuantity: remaining,
		Card:                  updated,
	}, nil
}

// getOwnedCard loads a card and enforces ownership.
func (s *Service) getOwnedCard(ctx context.Context, userID, cardID int) (*database.PlayerCard, error) {
	card, err := s.repository.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, ErrCardNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get card %d: %w", cardID, err)
	}
	if card.UserID != userID {
		return nil, ErrNotCardOwner
	}
	return card, nil
}

// resolvePokemonID resolves the canonical pokemon row behind a card. Legacy
// cards created before the pokemon_id column carry NULL; those are resolved
// by name and backfilled so future lookups are a single column read.
func (s *Service) resolvePokemonID(ctx context.Context, card *database.PlayerCard) (int, error) {
	pokemonID, err := s.repository.GetPokemonID(ctx, card.ID, card.UserID)
	if err != nil {
		return 0, fmt.Errorf("failed to read pokemon_id for card %d: %w", card.ID, err)
	}
	if pokemonID != nil {
		return *pokemonID, nil
	}

	p, err := s.evolutions.GetByName(ctx, card.PokemonName)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve pokemon %q for card %d: %w", card.PokemonName, card.ID, err)
	}

	// Best-effort backfill, mirroring the battle rewards path.
	if err := s.repository.SetPokemonID(ctx, card.ID, card.UserID, p.ID); err != nil {
		slog.Warn("failed to backfill pokemon_id", "card_id", card.ID, "error", err)
	}

	return p.ID, nil
}

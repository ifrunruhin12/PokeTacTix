package battle

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"pokemon-cli/internal/pokemon"
)

type SelectOptions struct {
	PlayerID      int
	Generation    *int
	MinLevel      int
	MaxLevel      int
	PokemonWindow int // K, default 5
	FamilyWindow  int // M, default 3
}

type EnemySelector interface {
	PickEnemy(ctx context.Context, opts SelectOptions) (*pokemon.Pokemon, error)
}

type enemySelector struct {
	pokemonService pokemon.PokemonService
	repo           pokemon.Repository
	history        History
}

func NewEnemySelector(ps pokemon.PokemonService, repo pokemon.Repository, history History) EnemySelector {
	return &enemySelector{
		pokemonService: ps,
		repo:           repo,
		history:        history,
	}
}

func (s *enemySelector) PickEnemy(ctx context.Context, opts SelectOptions) (*pokemon.Pokemon, error) {
	kWindow := opts.PokemonWindow
	if kWindow <= 0 {
		kWindow = 5
	}
	mWindow := opts.FamilyWindow
	if mWindow <= 0 {
		mWindow = 3
	}

	var recentPokemon []int
	var recentFamilies []int
	var err error

	if s.history != nil && opts.PlayerID > 0 {
		recentPokemon, _ = s.history.RecentPokemon(ctx, opts.PlayerID, kWindow)
		recentFamilies, _ = s.history.RecentFamilies(ctx, opts.PlayerID, mWindow)
	}

	var pool []*pokemon.Pokemon
	if s.repo != nil {
		pool, err = s.repo.CandidatePool(ctx, pokemon.PoolFilter{
			Generation: opts.Generation,
			MinLevel:   opts.MinLevel,
			MaxLevel:   opts.MaxLevel,
		})
	}

	// Fallback if repository pool is empty or not pre-seeded
	if err != nil || len(pool) == 0 {
		if s.pokemonService == nil {
			if err == nil {
				return nil, fmt.Errorf("no candidate pool available and pokemon service is nil")
			}
			return nil, fmt.Errorf("no candidate pool available and pokemon service is nil: %w", err)
		}

		// Try picking a random card from PokemonService
		card, err := s.pokemonService.GetRandomCard(ctx, false)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch fallback random pokemon: %w", err)
		}
		p, err := s.pokemonService.GetByID(ctx, card.CardID)
		if err != nil {
			return nil, fmt.Errorf("failed to load pokemon details for fallback card: %w", err)
		}
		pool = []*pokemon.Pokemon{p}
	}

	// Progressive relaxation filtering
	filtered := exclude(pool, recentPokemon, recentFamilies)
	if len(filtered) == 0 {
		// Drop family constraint first
		filtered = exclude(pool, recentPokemon, nil)
	}
	if len(filtered) == 0 {
		// Last resort: allow anything from candidate pool
		filtered = pool
	}

	chosen := weightedRandomPick(filtered)

	if s.history != nil && opts.PlayerID > 0 && chosen != nil {
		_ = s.history.RecordEncounter(ctx, opts.PlayerID, chosen)
	}

	return chosen, nil
}

// exclude filters out candidates whose ID is in recentPokemon or whose EvolutionChainID is in recentFamilies
func exclude(pool []*pokemon.Pokemon, recentPokemon []int, recentFamilies []int) []*pokemon.Pokemon {
	pMap := make(map[int]bool, len(recentPokemon))
	for _, id := range recentPokemon {
		pMap[id] = true
	}

	fMap := make(map[int]bool, len(recentFamilies))
	for _, famID := range recentFamilies {
		fMap[famID] = true
	}

	var result []*pokemon.Pokemon
	for _, p := range pool {
		if p == nil {
			continue
		}
		if pMap[p.ID] {
			continue
		}
		if fMap[p.EvolutionChainID] {
			continue
		}
		result = append(result, p)
	}
	return result
}

// weightedRandomPick selects a Pokémon from candidates using cryptographically safe random selection
func weightedRandomPick(pool []*pokemon.Pokemon) *pokemon.Pokemon {
	if len(pool) == 0 {
		return nil
	}
	if len(pool) == 1 {
		return pool[0]
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
	if err != nil {
		return pool[0]
	}
	return pool[nBig.Int64()]
}

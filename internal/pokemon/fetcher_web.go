package pokemon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	defaultService   PokemonService
	defaultServiceMu sync.RWMutex
)

// SetDefaultService sets global default PokemonService for legacy call-sites
func SetDefaultService(s PokemonService) {
	defaultServiceMu.Lock()
	defer defaultServiceMu.Unlock()
	defaultService = s
}

func getDefaultService() PokemonService {
	defaultServiceMu.RLock()
	defer defaultServiceMu.RUnlock()
	return defaultService
}

// GetMoves fetches move details from the API and returns up to 4 moves with
// positive power.
//
// Candidates are drawn in random order like before, but bounded to
// maxMoveCandidates and fetched concurrently. The previous sequential loop
// walked the Pokemon's entire move list one HTTP round trip at a time (many
// are zero-power status moves that get discarded), which regularly exceeded
// frontend timeouts on the shop purchase path even though the purchase
// itself succeeded server-side.
func GetMoves(rawMoves []RawMove) []Move {
	const (
		maxMoves          = 4
		maxMoveCandidates = 24
		maxConcurrency    = 8
	)

	if len(rawMoves) == 0 {
		return nil
	}

	perm := rand.Perm(len(rawMoves))
	if len(perm) > maxMoveCandidates {
		perm = perm[:maxMoveCandidates]
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	type moveResult struct {
		move Move
		ok   bool
	}

	results := make([]moveResult, len(perm))
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	for i, idx := range perm {
		wg.Add(1)
		go func(slot, rawIdx int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			moveURL := rawMoves[rawIdx].Move.URL
			resp, err := client.Get(moveURL)
			if err != nil {
				return
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
				return
			}

			if data.Power <= 0 {
				return
			}

			results[slot] = moveResult{
				ok: true,
				move: Move{
					Name:        data.Name,
					Power:       data.Power,
					StaminaCost: data.Power / 3,
					Type:        data.Type.Name,
				},
			}
		}(i, idx)
	}
	wg.Wait()

	// Preserve the shuffled candidate order when picking the final moves.
	var gameMoves []Move
	for _, res := range results {
		if !res.ok {
			continue
		}
		gameMoves = append(gameMoves, res.move)
		if len(gameMoves) == maxMoves {
			break
		}
	}

	return gameMoves
}

// FetchPokemon fetches Pokemon data from PokemonService (or legacy cold PokéAPI fallback)
func FetchPokemon(name string) (RawPokeAPIPokemon, []Move, error) {
	if svc := getDefaultService(); svc != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		p, err := svc.GetByName(ctx, name)
		if err == nil && p != nil {
			var raw RawPokeAPIPokemon
			if len(p.RawJSON) > 0 && json.Unmarshal(p.RawJSON, &raw) == nil && raw.Name != "" && len(raw.Stats) > 0 {
				return raw, GetMoves(raw.Moves), nil
			}
		}
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(name)

	// Create HTTP client with timeout to prevent hanging
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return RawPokeAPIPokemon{}, nil, fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return RawPokeAPIPokemon{}, nil, fmt.Errorf("Pokemon \"%s\" not found. Please check the name and try again", name)
	}

	defer resp.Body.Close()

	var poke RawPokeAPIPokemon

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&poke); err != nil {
		return RawPokeAPIPokemon{}, nil, fmt.Errorf("failed to decode pokemon data: %w", err)
	}

	pokeMoves := GetMoves(poke.Moves)
	return poke, pokeMoves, nil
}

// FetchRandomPokemonCard returns a random Card from PokemonService (or legacy cold path fallback)
func FetchRandomPokemonCard(allowSpecial bool) Card {
	if svc := getDefaultService(); svc != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		card, err := svc.GetRandomCard(ctx, allowSpecial)
		if err == nil {
			if len(card.Moves) == 0 {
				card.Moves = []Move{
					{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
				}
			}
			return card
		}
	}
	mythicalOdds := 0.0001  // 0.01%
	legendaryOdds := 0.0001 // 0.01%
	maxRetries := 5

	for range maxRetries {
		roll := rand.Float64()
		var name string
		if roll < mythicalOdds {
			name = mythicalNames[rand.Intn(len(mythicalNames))]
		} else if roll < mythicalOdds+legendaryOdds {
			name = legendaryNames[rand.Intn(len(legendaryNames))]
		} else {
			id := rand.Intn(649) + 1
			name = fmt.Sprintf("%d", id)
		}
		poke, moves, err := FetchPokemon(name)
		if err != nil {
			continue
		}
		card := BuildCardFromPokemon(poke, moves)
		if len(card.Moves) == 0 {
			card.Moves = []Move{
				{Name: "tackle", Power: 40, StaminaCost: 13, Type: "normal"},
			}
		}
		return card
	}

	// Fallback to Pikachu
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

	// Ultimate fallback
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

package pokemon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingChainClient struct {
	calls int
}

func (c *failingChainClient) FetchPokemonRaw(context.Context, string) ([]byte, error) {
	return nil, errors.New("unexpected pokemon fetch")
}

func (c *failingChainClient) FetchSpeciesRaw(context.Context, string) ([]byte, error) {
	return nil, errors.New("unexpected species fetch")
}

func (c *failingChainClient) FetchEvolutionChainRaw(context.Context, int) ([]byte, error) {
	c.calls++
	return nil, errors.New("pokeapi unavailable")
}

// TestGetEvolutionForLevelSchedulesWarmupOnMissingChain verifies that a lookup
// which finds no chain data at all returns the no-evolution answer without
// erroring, and hands the retry to the background warm-up.
func TestGetEvolutionForLevelSchedulesWarmupOnMissingChain(t *testing.T) {
	repo := newMockRepo()
	repo.pokemon[4] = &Pokemon{
		ID:               4,
		Name:             "charmander",
		SpeciesID:        4,
		EvolutionChainID: 2,
		Generation:       1,
		Types:            []string{"fire"},
		BaseStats:        BaseStats{HP: 39, Attack: 52, Defense: 43, Speed: 65},
		FetchedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	client := &failingChainClient{}
	svc := NewService(nil, repo, client).(*service)

	target, err := svc.GetEvolutionForLevel(context.Background(), 4, 20)
	require.NoError(t, err)
	assert.Nil(t, target, "missing chain data must degrade to no evolution, not an error")
}

// TestWarmEvolutionChainHonorsBackoff verifies the warm-up respects the
// failure backoff window instead of retrying a chain that just failed.
func TestWarmEvolutionChainHonorsBackoff(t *testing.T) {
	client := &failingChainClient{}
	svc := NewService(nil, nil, client).(*service)

	// First refresh fails and opens the backoff window.
	svc.warmEvolutionChain(2)
	assert.Equal(t, 1, client.calls, "first warm-up should attempt one fetch")

	// A warm-up inside the backoff window must not hit PokéAPI again.
	svc.warmEvolutionChain(2)
	assert.Equal(t, 1, client.calls, "warm-up inside backoff window must not refetch")
}

// TestWarmEvolutionChainPersistsResult verifies a successful warm-up stores
// the chain in the durable cache so later lookups skip the cold path.
func TestWarmEvolutionChainPersistsResult(t *testing.T) {
	repo := newMockRepo()
	client := &mockPokeAPIClient{
		rawEvolutionChain: []byte(`{
			"chain": {
				"species": {"url": "https://pokeapi.co/api/v2/pokemon-species/4/"},
				"evolves_to": [{
					"species": {"url": "https://pokeapi.co/api/v2/pokemon-species/5/"},
					"evolution_details": [{"min_level": 16, "trigger": {"name": "level-up"}}],
					"evolves_to": []
				}]
			}
		}`),
	}
	svc := NewService(nil, repo, client).(*service)

	svc.warmEvolutionChain(2)

	chain, err := repo.GetEvolutionChain(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, chain.Links, 1)
	assert.Equal(t, 5, chain.Links[0].ToSpeciesID)
}

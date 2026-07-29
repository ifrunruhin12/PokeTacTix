package battle

import (
	"context"
	"testing"

	"pokemon-cli/internal/pokemon"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHistory struct {
	recentPokemon  map[int][]int
	recentFamilies map[int][]int
	encounters     []int
}

func newMockHistory() *mockHistory {
	return &mockHistory{
		recentPokemon:  make(map[int][]int),
		recentFamilies: make(map[int][]int),
	}
}

func (m *mockHistory) RecentPokemon(ctx context.Context, playerID int, n int) ([]int, error) {
	list := m.recentPokemon[playerID]
	if len(list) > n {
		return list[:n], nil
	}
	return list, nil
}

func (m *mockHistory) RecentFamilies(ctx context.Context, playerID int, n int) ([]int, error) {
	list := m.recentFamilies[playerID]
	if len(list) > n {
		return list[:n], nil
	}
	return list, nil
}

func (m *mockHistory) RecordEncounter(ctx context.Context, playerID int, p *pokemon.Pokemon) error {
	m.recentPokemon[playerID] = append([]int{p.ID}, m.recentPokemon[playerID]...)
	m.recentFamilies[playerID] = append([]int{p.EvolutionChainID}, m.recentFamilies[playerID]...)
	m.encounters = append(m.encounters, p.ID)
	return nil
}

type mockRepo struct {
	pool []*pokemon.Pokemon
}

func (m *mockRepo) GetPokemon(ctx context.Context, id int) (*pokemon.Pokemon, error) {
	for _, p := range m.pool {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, pokemon.ErrPokemonNotFound
}

func (m *mockRepo) GetPokemonByName(ctx context.Context, name string) (*pokemon.Pokemon, error) {
	for _, p := range m.pool {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, pokemon.ErrPokemonNotFound
}

func (m *mockRepo) UpsertPokemon(ctx context.Context, p *pokemon.Pokemon) error {
	m.pool = append(m.pool, p)
	return nil
}

func (m *mockRepo) GetEvolutionChain(ctx context.Context, id int) (*pokemon.EvolutionChain, error) {
	return nil, pokemon.ErrEvolutionChainNotFound
}

func (m *mockRepo) UpsertEvolutionChain(ctx context.Context, ec *pokemon.EvolutionChain) error {
	return nil
}

func (m *mockRepo) CandidatePool(ctx context.Context, filter pokemon.PoolFilter) ([]*pokemon.Pokemon, error) {
	var result []*pokemon.Pokemon
	for _, p := range m.pool {
		if filter.Generation != nil && p.Generation != *filter.Generation {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

func TestExcludeLogic(t *testing.T) {
	pool := []*pokemon.Pokemon{
		{ID: 1, Name: "bulbasaur", EvolutionChainID: 1},  // Gen 1 Grass family
		{ID: 2, Name: "ivysaur", EvolutionChainID: 1},    // Gen 1 Grass family (same evolution chain)
		{ID: 4, Name: "charmander", EvolutionChainID: 2}, // Gen 1 Fire family
		{ID: 7, Name: "squirtle", EvolutionChainID: 3},   // Gen 1 Water family
	}

	// Exclude recent pokemon ID 1 and recent families 1 and 2
	filtered := exclude(pool, []int{1}, []int{1, 2})
	require.Len(t, filtered, 1)
	assert.Equal(t, 7, filtered[0].ID) // Only Squirtle remains!
}

func TestProgressiveRelaxation(t *testing.T) {
	ctx := context.Background()
	history := newMockHistory()

	// Pool has only Bulbasaur and Ivysaur (both evolution_chain_id = 1)
	pool := []*pokemon.Pokemon{
		{ID: 1, Name: "bulbasaur", EvolutionChainID: 1, BaseStats: pokemon.BaseStats{HP: 45}},
		{ID: 2, Name: "ivysaur", EvolutionChainID: 1, BaseStats: pokemon.BaseStats{HP: 60}},
	}

	repo := &mockRepo{pool: pool}
	selector := NewEnemySelector(nil, repo, history)

	// Simulate player fought Bulbasaur (#1, family #1)
	history.recentPokemon[100] = []int{1}
	history.recentFamilies[100] = []int{1}

	// Pick enemy: Since both #1 and family #1 are restricted, exclude(pool, [1], [1]) is empty.
	// Progressive relaxation drops family constraint: exclude(pool, [1], nil) -> leaves Ivysaur (#2)!
	chosen, err := selector.PickEnemy(ctx, SelectOptions{
		PlayerID:      100,
		PokemonWindow: 5,
		FamilyWindow:  3,
	})

	require.NoError(t, err)
	assert.Equal(t, 2, chosen.ID)
}

func TestAntiRepeatSequence(t *testing.T) {
	ctx := context.Background()
	history := newMockHistory()

	// Pool of 6 distinct evolution chains
	pool := []*pokemon.Pokemon{
		{ID: 1, Name: "bulbasaur", EvolutionChainID: 1, Generation: 1},
		{ID: 4, Name: "charmander", EvolutionChainID: 2, Generation: 1},
		{ID: 7, Name: "squirtle", EvolutionChainID: 3, Generation: 1},
		{ID: 10, Name: "caterpie", EvolutionChainID: 4, Generation: 1},
		{ID: 13, Name: "weedle", EvolutionChainID: 5, Generation: 1},
		{ID: 16, Name: "pidgey", EvolutionChainID: 6, Generation: 1},
	}

	repo := &mockRepo{pool: pool}
	selector := NewEnemySelector(nil, repo, history)

	playerID := 555
	chosenIDs := make([]int, 0, 5)

	for i := 0; i < 5; i++ {
		p, err := selector.PickEnemy(ctx, SelectOptions{
			PlayerID:      playerID,
			PokemonWindow: 5,
			FamilyWindow:  3,
		})
		require.NoError(t, err)
		chosenIDs = append(chosenIDs, p.ID)
	}

	// Verify all 5 picked Pokémon IDs are unique!
	uniqueMap := make(map[int]bool)
	for _, id := range chosenIDs {
		assert.False(t, uniqueMap[id], "Duplicate Pokémon ID %d found in sequence of 5 fights", id)
		uniqueMap[id] = true
	}
}

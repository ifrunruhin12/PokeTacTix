package pokemon

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCache struct {
	pokemon map[int]*Pokemon
	names   map[string]int
	chains  map[int]*EvolutionChain
}

func newMockCache() *mockCache {
	return &mockCache{
		pokemon: make(map[int]*Pokemon),
		names:   make(map[string]int),
		chains:  make(map[int]*EvolutionChain),
	}
}

func (m *mockCache) GetPokemon(ctx context.Context, id int) (*Pokemon, error) {
	if p, ok := m.pokemon[id]; ok {
		return p, nil
	}
	return nil, assert.AnError
}

func (m *mockCache) SetPokemon(ctx context.Context, p *Pokemon) error {
	m.pokemon[p.ID] = p
	m.names[p.Name] = p.ID
	return nil
}

func (m *mockCache) GetIDByName(ctx context.Context, name string) (int, error) {
	if id, ok := m.names[name]; ok {
		return id, nil
	}
	return 0, assert.AnError
}

func (m *mockCache) SetIDByName(ctx context.Context, name string, id int) error {
	m.names[name] = id
	return nil
}

func (m *mockCache) GetEvolutionChain(ctx context.Context, id int) (*EvolutionChain, error) {
	if ec, ok := m.chains[id]; ok {
		return ec, nil
	}
	return nil, assert.AnError
}

func (m *mockCache) SetEvolutionChain(ctx context.Context, ec *EvolutionChain) error {
	m.chains[ec.ID] = ec
	return nil
}

type mockRepo struct {
	pokemon map[int]*Pokemon
	chains  map[int]*EvolutionChain
}

type mockPokeAPIClient struct {
	rawPokemon          []byte
	rawSpecies          []byte
	rawEvolutionChain   []byte
	pokemonCalls        int
	speciesCalls        int
	evolutionChainCalls int
}

func (m *mockPokeAPIClient) FetchPokemonRaw(context.Context, string) ([]byte, error) {
	m.pokemonCalls++
	return m.rawPokemon, nil
}

func (m *mockPokeAPIClient) FetchSpeciesRaw(context.Context, string) ([]byte, error) {
	m.speciesCalls++
	return m.rawSpecies, nil
}

func (m *mockPokeAPIClient) FetchEvolutionChainRaw(context.Context, int) ([]byte, error) {
	m.evolutionChainCalls++
	return m.rawEvolutionChain, nil
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		pokemon: make(map[int]*Pokemon),
		chains:  make(map[int]*EvolutionChain),
	}
}

func (m *mockRepo) GetPokemon(ctx context.Context, id int) (*Pokemon, error) {
	if p, ok := m.pokemon[id]; ok {
		return p, nil
	}
	return nil, ErrPokemonNotFound
}

func (m *mockRepo) GetPokemonByName(ctx context.Context, name string) (*Pokemon, error) {
	for _, p := range m.pokemon {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrPokemonNotFound
}

func (m *mockRepo) UpsertPokemon(ctx context.Context, p *Pokemon) error {
	m.pokemon[p.ID] = p
	return nil
}

func (m *mockRepo) GetEvolutionChain(ctx context.Context, id int) (*EvolutionChain, error) {
	if ec, ok := m.chains[id]; ok {
		return ec, nil
	}
	return nil, ErrEvolutionChainNotFound
}

func (m *mockRepo) UpsertEvolutionChain(ctx context.Context, ec *EvolutionChain) error {
	m.chains[ec.ID] = ec
	return nil
}

func (m *mockRepo) CandidatePool(ctx context.Context, filter PoolFilter) ([]*Pokemon, error) {
	var list []*Pokemon
	for _, p := range m.pokemon {
		if filter.Generation != nil && p.Generation != *filter.Generation {
			continue
		}
		list = append(list, p)
	}
	return list, nil
}

func TestNormalizePokemon(t *testing.T) {
	rawPokemon := []byte(`{
		"id": 25,
		"name": "pikachu",
		"stats": [
			{"base_stat": 35, "stat": {"name": "hp"}},
			{"base_stat": 55, "stat": {"name": "attack"}},
			{"base_stat": 40, "stat": {"name": "defense"}},
			{"base_stat": 50, "stat": {"name": "special-attack"}},
			{"base_stat": 50, "stat": {"name": "special-defense"}},
			{"base_stat": 90, "stat": {"name": "speed"}}
		],
		"types": [
			{"type": {"name": "electric"}}
		],
		"sprites": {
			"front_default": "http://example.com/pikachu.png"
		}
	}`)

	rawSpecies := []byte(`{
		"id": 25,
		"generation": {"name": "generation-i"},
		"evolution_chain": {"url": "https://pokeapi.co/api/v2/evolution-chain/10/"}
	}`)

	p, evolutionChainURL, err := normalizePokemon(rawPokemon, rawSpecies)
	require.NoError(t, err)
	assert.Equal(t, "https://pokeapi.co/api/v2/evolution-chain/10/", evolutionChainURL)
	assert.Equal(t, 25, p.ID)
	assert.Equal(t, "pikachu", p.Name)
	assert.Equal(t, 1, p.Generation)
	assert.Equal(t, 35, p.BaseStats.HP)
	assert.Equal(t, 55, p.BaseStats.Attack)
	assert.Equal(t, []string{"electric"}, p.Types)
	assert.Equal(t, "http://example.com/pikachu.png", p.SpriteURL)

	card := p.ToCard()
	assert.Equal(t, 25, card.CardID)
	assert.Equal(t, "pikachu", card.Name)
	assert.Equal(t, 52, card.HP)
	assert.Equal(t, 180, card.Stamina) // speed 90 * 2
}

func TestTieredFetchOrder(t *testing.T) {
	ctx := context.Background()
	cache := newMockCache()
	repo := newMockRepo()

	// Seed DB with Charmander
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

	svc := NewService(cache, repo, nil)

	// Fetch ID 4: Misses Redis, hits DB, backfills Redis
	p, err := svc.GetByID(ctx, 4)
	require.NoError(t, err)
	assert.Equal(t, "charmander", p.Name)

	// Verify Redis now has Charmander
	cachedP, err := cache.GetPokemon(ctx, 4)
	require.NoError(t, err)
	assert.Equal(t, "charmander", cachedP.Name)
}

func TestGetByIDReusesSpeciesEvolutionChainURL(t *testing.T) {
	client := &mockPokeAPIClient{
		rawPokemon: []byte(`{"id":25,"name":"pikachu"}`),
		rawSpecies: []byte(`{
			"id": 25,
			"evolution_chain": {"url": "https://pokeapi.co/api/v2/evolution-chain/10/"}
		}`),
		rawEvolutionChain: []byte(`{
			"chain": {
				"species": {"url": "https://pokeapi.co/api/v2/pokemon-species/25/"},
				"evolves_to": []
			}
		}`),
	}

	p, err := NewService(nil, nil, client).GetByID(context.Background(), 25)
	require.NoError(t, err)
	assert.Equal(t, 10, p.EvolutionChainID)
	assert.Equal(t, 1, client.pokemonCalls)
	assert.Equal(t, 1, client.speciesCalls)
	assert.Equal(t, 1, client.evolutionChainCalls)
}

func TestExtractMemberSpeciesIDs(t *testing.T) {
	rawChain := []byte(`{
		"chain": {
			"species": {"url": "https://pokeapi.co/api/v2/pokemon-species/1/"},
			"evolves_to": [
				{
					"species": {"url": "https://pokeapi.co/api/v2/pokemon-species/2/"},
					"evolves_to": [
						{
							"species": {"url": "https://pokeapi.co/api/v2/pokemon-species/3/"},
							"evolves_to": []
						}
					]
				}
			]
		}
	}`)

	ids, err := ExtractMemberSpeciesIDs(rawChain)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, ids)
}

package cards

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"pokemon-cli/internal/database"
	"pokemon-cli/internal/items"
	"pokemon-cli/internal/pokemon"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mocks -----------------------------------------------------------------

type evolveCall struct {
	userID, cardID, expectedPokemonID int
	target                            *pokemon.Pokemon
	itemID                            string
}

type mockCardRepository struct {
	cards      map[int]*database.PlayerCard
	pokemonIDs map[int]int // cardID -> pokemon_id (absent key = NULL)
	setIDs     map[int]int // cardID -> backfilled pokemon_id

	evolveCalls []evolveCall
	evolveErr   error // when set, EvolveWithItem fails without consuming
}

func (m *mockCardRepository) Create(ctx context.Context, card *database.PlayerCard) (*database.PlayerCard, error) {
	card.ID = len(m.cards) + 1
	m.cards[card.ID] = card
	return card, nil
}

func (m *mockCardRepository) GetByID(ctx context.Context, id int) (*database.PlayerCard, error) {
	if card, ok := m.cards[id]; ok {
		return card, nil
	}
	return nil, ErrCardNotFound
}

func (m *mockCardRepository) GetUserCards(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	ids := make([]int, 0, len(m.cards))
	for id, card := range m.cards {
		if card.UserID == userID {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	owned := make([]database.PlayerCard, 0, len(ids))
	for _, id := range ids {
		owned = append(owned, *m.cards[id])
	}
	return owned, nil
}

func (m *mockCardRepository) GetUserDeck(ctx context.Context, userID int) ([]database.PlayerCard, error) {
	return nil, nil
}

func (m *mockCardRepository) UpdateDeck(ctx context.Context, userID int, cardIDs []int) error {
	return nil
}

func (m *mockCardRepository) AddXP(ctx context.Context, cardID int, xp int) (*database.PlayerCard, error) {
	return m.GetByID(ctx, cardID)
}

func (m *mockCardRepository) GetPokemonID(ctx context.Context, cardID, userID int) (*int, error) {
	if id, ok := m.pokemonIDs[cardID]; ok {
		return &id, nil
	}
	return nil, nil
}

func (m *mockCardRepository) SetPokemonID(ctx context.Context, cardID, userID, pokemonID int) error {
	if m.setIDs == nil {
		m.setIDs = make(map[int]int)
	}
	m.setIDs[cardID] = pokemonID
	return nil
}

func (m *mockCardRepository) EvolveWithItem(ctx context.Context, userID, cardID, expectedPokemonID int, target *pokemon.Pokemon, itemID string) (int, error) {
	m.evolveCalls = append(m.evolveCalls, evolveCall{
		userID: userID, cardID: cardID, expectedPokemonID: expectedPokemonID,
		target: target, itemID: itemID,
	})
	if m.evolveErr != nil {
		return 0, m.evolveErr
	}
	// Simulate the transaction: card becomes the target species.
	if card, ok := m.cards[cardID]; ok {
		card.PokemonName = target.Name
	}
	return 0, nil // remaining quantity; irrelevant to the service layer
}

type itemKey struct {
	pokemonID int
	itemID    string
}

type mockEvolutionRules struct {
	byName      map[string]*pokemon.Pokemon
	itemTargets map[itemKey]*pokemon.Pokemon
	options     map[int][]pokemon.EvolutionOption
}

func (m *mockEvolutionRules) GetByName(ctx context.Context, name string) (*pokemon.Pokemon, error) {
	if p, ok := m.byName[name]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("pokemon %q not found", name)
}

func (m *mockEvolutionRules) GetEvolutionForItem(ctx context.Context, pokemonID int, itemID string) (*pokemon.Pokemon, error) {
	return m.itemTargets[itemKey{pokemonID: pokemonID, itemID: itemID}], nil
}

func (m *mockEvolutionRules) GetEvolutionOptions(ctx context.Context, pokemonID int) ([]pokemon.EvolutionOption, error) {
	return m.options[pokemonID], nil
}

type mockItemCatalog struct {
	items      map[string]*items.Item
	quantities map[int]map[string]int // userID -> itemID -> quantity
}

func (m *mockItemCatalog) GetItem(ctx context.Context, itemID string) (*items.Item, error) {
	if item, ok := m.items[itemID]; ok {
		return item, nil
	}
	return nil, items.ErrItemNotFound
}

func (m *mockItemCatalog) GetQuantity(ctx context.Context, userID int, itemID string) (int, error) {
	return m.quantities[userID][itemID], nil
}

// --- fixtures --------------------------------------------------------------

const (
	userA = 1
	userB = 2
)

var pikachuPokemon = &pokemon.Pokemon{ID: 25, Name: "pikachu", SpeciesID: 25}
var raichuPokemon = &pokemon.Pokemon{ID: 26, Name: "raichu", SpeciesID: 26, Types: []string{"electric"}}

func thunderStone() *items.Item {
	return &items.Item{ID: "thunder-stone", Name: "Thunder Stone", ItemType: items.TypeEvolution, Price: 800}
}

func newEvolutionTestService(t *testing.T) (*Service, *mockCardRepository, *mockEvolutionRules, *mockItemCatalog) {
	t.Helper()

	repo := &mockCardRepository{
		cards: map[int]*database.PlayerCard{
			101: {ID: 101, UserID: userA, PokemonName: "pikachu", Level: 12},
			102: {ID: 102, UserID: userA, PokemonName: "raichu", Level: 20},
			103: {ID: 103, UserID: userB, PokemonName: "pikachu", Level: 12},
		},
		pokemonIDs: map[int]int{101: 25, 102: 26, 103: 25},
	}

	rules := &mockEvolutionRules{
		byName: map[string]*pokemon.Pokemon{
			"pikachu": pikachuPokemon,
			"raichu":  raichuPokemon,
		},
		itemTargets: map[itemKey]*pokemon.Pokemon{
			{pokemonID: 25, itemID: "thunder-stone"}: raichuPokemon,
		},
		options: map[int][]pokemon.EvolutionOption{
			25: {
				{Method: pokemon.MethodItem, ToSpeciesID: 26, Item: "thunder-stone"},
			},
			26: nil, // final form
		},
	}

	catalog := &mockItemCatalog{
		items: map[string]*items.Item{
			"thunder-stone": thunderStone(),
			"fire-stone":    {ID: "fire-stone", Name: "Fire Stone", ItemType: items.TypeEvolution, Price: 800},
		},
		quantities: map[int]map[string]int{},
	}

	svc := NewService(repo, rules, catalog)
	return svc, repo, rules, catalog
}

func give(t *testing.T, catalog *mockItemCatalog, userID int, itemID string, qty int) {
	t.Helper()
	if catalog.quantities[userID] == nil {
		catalog.quantities[userID] = make(map[string]int)
	}
	catalog.quantities[userID][itemID] = qty
}

// --- evolution -------------------------------------------------------------

func TestEvolveCardWithItem(t *testing.T) {
	t.Run("pikachu with thunder stone evolves into raichu", func(t *testing.T) {
		svc, repo, _, catalog := newEvolutionTestService(t)
		give(t, catalog, userA, "thunder-stone", 1)

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 101, "thunder-stone")
		require.NoError(t, err)

		assert.Equal(t, "pikachu", result.EvolvedFrom)
		assert.Equal(t, "raichu", result.EvolvedInto)
		assert.Equal(t, "thunder-stone", result.ItemConsumed)

		// Exactly one consumption attempt, against the right card and item.
		require.Len(t, repo.evolveCalls, 1)
		call := repo.evolveCalls[0]
		assert.Equal(t, userA, call.userID)
		assert.Equal(t, 101, call.cardID)
		assert.Equal(t, 25, call.expectedPokemonID)
		assert.Equal(t, "thunder-stone", call.itemID)
		assert.Same(t, raichuPokemon, call.target)

		// The card itself became a raichu.
		assert.Equal(t, "raichu", repo.cards[101].PokemonName)
	})

	t.Run("pikachu without thunder stone is rejected", func(t *testing.T) {
		svc, repo, _, _ := newEvolutionTestService(t) // no inventory

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 101, "thunder-stone")
		require.ErrorIs(t, err, items.ErrInsufficientItem)
		assert.Nil(t, result)

		// The item was never touched: no evolve transaction ran.
		assert.Empty(t, repo.evolveCalls)
		assert.Equal(t, "pikachu", repo.cards[101].PokemonName)
	})

	t.Run("pikachu with the wrong item is rejected", func(t *testing.T) {
		svc, repo, _, catalog := newEvolutionTestService(t)
		give(t, catalog, userA, "fire-stone", 3)

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 101, "fire-stone")
		require.ErrorIs(t, err, ErrNotEligible)
		assert.Nil(t, result)

		assert.Empty(t, repo.evolveCalls)
		assert.Equal(t, "pikachu", repo.cards[101].PokemonName)
	})

	t.Run("non-eligible pokemon with thunder stone is rejected", func(t *testing.T) {
		svc, repo, _, catalog := newEvolutionTestService(t)
		give(t, catalog, userA, "thunder-stone", 2) // card 102 is already raichu

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 102, "thunder-stone")
		require.ErrorIs(t, err, ErrNotEligible)
		assert.Nil(t, result)

		assert.Empty(t, repo.evolveCalls)
		assert.Equal(t, "raichu", repo.cards[102].PokemonName)
	})

	t.Run("unknown item is rejected", func(t *testing.T) {
		svc, repo, _, _ := newEvolutionTestService(t)

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 101, "moon-stone")
		require.ErrorIs(t, err, items.ErrItemNotFound)
		assert.Nil(t, result)
		assert.Empty(t, repo.evolveCalls)
	})

	t.Run("failed evolution does not consume the item", func(t *testing.T) {
		svc, repo, _, catalog := newEvolutionTestService(t)
		give(t, catalog, userA, "thunder-stone", 1)
		repo.evolveErr = errors.New("consume guard failed") // tx fails after pre-checks

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 101, "thunder-stone")
		require.Error(t, err)
		assert.Nil(t, result)

		// The card is untouched: the whole transaction rolled back.
		assert.Equal(t, "pikachu", repo.cards[101].PokemonName)
	})

	t.Run("player cannot evolve another player's pokemon", func(t *testing.T) {
		svc, repo, _, catalog := newEvolutionTestService(t)
		give(t, catalog, userA, "thunder-stone", 1) // owner has the stone, card belongs to userB

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 103, "thunder-stone")
		require.ErrorIs(t, err, ErrNotCardOwner)
		assert.Nil(t, result)
		assert.Empty(t, repo.evolveCalls)
		assert.Equal(t, "pikachu", repo.cards[103].PokemonName)
	})

	t.Run("missing card is rejected", func(t *testing.T) {
		svc, _, _, _ := newEvolutionTestService(t)

		result, err := svc.EvolveCardWithItem(context.Background(), userA, 999, "thunder-stone")
		require.ErrorIs(t, err, ErrCardNotFound)
		assert.Nil(t, result)
	})
}

func TestGetEvolutionInfo(t *testing.T) {
	t.Run("item evolution reports required item and ownership", func(t *testing.T) {
		svc, _, _, catalog := newEvolutionTestService(t)
		give(t, catalog, userA, "thunder-stone", 2)

		info, err := svc.GetEvolutionInfo(context.Background(), userA, 101)
		require.NoError(t, err)

		assert.Equal(t, "pikachu", info.PokemonName)
		assert.True(t, info.CanEvolveNow)
		require.Len(t, info.Options, 1)

		opt := info.Options[0]
		assert.Equal(t, pokemon.MethodItem, opt.Method)
		assert.Equal(t, "thunder-stone", opt.RequiredItem)
		assert.Equal(t, "Thunder Stone", opt.RequiredItemName)
		assert.Equal(t, 2, opt.OwnedQuantity)
		assert.True(t, opt.Eligible)
		assert.Empty(t, opt.MinLevel) // item evolution must not invent a level requirement
	})

	t.Run("item evolution with empty inventory is not eligible", func(t *testing.T) {
		svc, _, _, _ := newEvolutionTestService(t)

		info, err := svc.GetEvolutionInfo(context.Background(), userA, 101)
		require.NoError(t, err)

		assert.False(t, info.CanEvolveNow)
		require.Len(t, info.Options, 1)
		assert.False(t, info.Options[0].Eligible)
		assert.Zero(t, info.Options[0].OwnedQuantity)
	})

	t.Run("final form has no options", func(t *testing.T) {
		svc, _, _, _ := newEvolutionTestService(t)

		info, err := svc.GetEvolutionInfo(context.Background(), userA, 102)
		require.NoError(t, err)
		assert.Empty(t, info.Options)
		assert.False(t, info.CanEvolveNow)
	})

	t.Run("ownership is enforced", func(t *testing.T) {
		svc, _, _, _ := newEvolutionTestService(t)

		info, err := svc.GetEvolutionInfo(context.Background(), userA, 103)
		require.ErrorIs(t, err, ErrNotCardOwner)
		assert.Nil(t, info)
	})

	t.Run("legacy card without pokemon_id is resolved by name and backfilled", func(t *testing.T) {
		svc, repo, _, _ := newEvolutionTestService(t)
		delete(repo.pokemonIDs, 101) // legacy NULL pokemon_id

		info, err := svc.GetEvolutionInfo(context.Background(), userA, 101)
		require.NoError(t, err)
		assert.Equal(t, "pikachu", info.PokemonName)

		require.Len(t, repo.evolveCalls, 0) // info must never evolve or consume
		assert.Equal(t, 25, repo.setIDs[101], "pokemon_id should be backfilled")
	})
}

// Level-based evolution shares the EvolutionOption infrastructure but must not
// require an item, and item evolution must not become eligible from level alone.
func TestGetEvolutionInfoLevelRules(t *testing.T) {
	charmander := &pokemon.Pokemon{ID: 4, Name: "charmander", SpeciesID: 4}

	repo := &mockCardRepository{
		cards: map[int]*database.PlayerCard{
			201: {ID: 201, UserID: userA, PokemonName: "charmander", Level: 16},
		},
		pokemonIDs: map[int]int{201: 4},
	}
	rules := &mockEvolutionRules{
		byName: map[string]*pokemon.Pokemon{"charmander": charmander},
		options: map[int][]pokemon.EvolutionOption{
			4: {
				{Method: pokemon.MethodLevel, ToSpeciesID: 5, MinLevel: 16, TargetName: "charmeleon"},
			},
		},
	}
	catalog := &mockItemCatalog{items: map[string]*items.Item{}, quantities: map[int]map[string]int{}}

	svc := NewService(repo, rules, catalog)

	info, err := svc.GetEvolutionInfo(context.Background(), userA, 201)
	require.NoError(t, err)

	require.Len(t, info.Options, 1)
	opt := info.Options[0]
	assert.Equal(t, pokemon.MethodLevel, opt.Method)
	assert.Equal(t, 16, opt.MinLevel)
	assert.True(t, opt.Eligible, "level requirement met")
	assert.Empty(t, opt.RequiredItem, "level evolution must not require an item")
	assert.False(t, info.CanEvolveNow, "level evolution is applied automatically, not from the deck")
}

func TestGetEvolutionInfoFinalForm(t *testing.T) {
	svc, _, _, _ := newEvolutionTestService(t)

	// Raichu has no evolution links; Options must serialize as [] (not null)
	// so clients can render the empty list without crashing.
	info, err := svc.GetEvolutionInfo(context.Background(), userA, 102)
	require.NoError(t, err)

	assert.NotNil(t, info.Options, "options must never be a nil slice")
	assert.Empty(t, info.Options)
	assert.False(t, info.CanEvolveNow)
}

func TestGetEvolutionSummaries(t *testing.T) {
	svc, _, _, _ := newEvolutionTestService(t)

	summaries, err := svc.GetEvolutionSummaries(context.Background(), userA)
	require.NoError(t, err)

	byCard := make(map[int]bool, len(summaries))
	for _, s := range summaries {
		byCard[s.CardID] = s.HasEvolution
	}

	// Only the user's own cards are summarized.
	assert.Len(t, summaries, 2)
	assert.True(t, byCard[101], "pikachu has an evolution path")
	assert.False(t, byCard[102], "raichu is a final form")
	assert.NotContains(t, byCard, 103, "another player's card is not summarized")
}

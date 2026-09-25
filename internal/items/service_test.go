package items

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type purchaseRecord struct {
	userID   int
	itemID   string
	quantity int
}

type mockCatalogRepository struct {
	items       map[string]*Item
	quantities  map[int]map[string]int
	purchases   []purchaseRecord
	purchaseErr error

	boosts          []ActiveBoost
	activated       []string // item ids activated via ActivateBooster
	activateErr     error
	activeBoostsErr error
	tickedBattles   []int // user ids passed to ConsumeBattleBoosts
}

func (m *mockCatalogRepository) ListItems(ctx context.Context) ([]Item, error) {
	var result []Item
	for _, item := range m.items {
		result = append(result, *item)
	}
	return result, nil
}

func (m *mockCatalogRepository) GetItem(ctx context.Context, itemID string) (*Item, error) {
	if item, ok := m.items[itemID]; ok {
		return item, nil
	}
	return nil, ErrItemNotFound
}

func (m *mockCatalogRepository) GetInventory(ctx context.Context, userID int) ([]InventoryEntry, error) {
	var result []InventoryEntry
	for itemID, qty := range m.quantities[userID] {
		if qty <= 0 || m.items[itemID] == nil {
			continue
		}
		result = append(result, InventoryEntry{Item: *m.items[itemID], Quantity: qty})
	}
	return result, nil
}

func (m *mockCatalogRepository) GetQuantity(ctx context.Context, userID int, itemID string) (int, error) {
	return m.quantities[userID][itemID], nil
}

// Purchase simulates the SQL upsert: quantity accumulates per (user, item).
func (m *mockCatalogRepository) Purchase(ctx context.Context, userID int, item *Item, quantity int) (int, error) {
	m.purchases = append(m.purchases, purchaseRecord{userID: userID, itemID: item.ID, quantity: quantity})
	if m.purchaseErr != nil {
		return 0, m.purchaseErr
	}
	if m.quantities[userID] == nil {
		m.quantities[userID] = make(map[string]int)
	}
	m.quantities[userID][item.ID] += quantity
	return m.quantities[userID][item.ID], nil
}

// ActivateBooster simulates the SQL transaction: consume one from the
// inventory and append the active boost.
func (m *mockCatalogRepository) ActivateBooster(ctx context.Context, userID int, item *Item, effect BoosterEffect) (*ActiveBoost, error) {
	if m.activateErr != nil {
		return nil, m.activateErr
	}
	if m.quantities[userID][item.ID] <= 0 {
		return nil, ErrInsufficientItem
	}
	m.quantities[userID][item.ID]--
	m.activated = append(m.activated, item.ID)
	m.boosts = append(m.boosts, ActiveBoost{
		ID: len(m.boosts) + 1, ItemID: item.ID, ItemName: item.Name,
		Stat: effect.Stat, Bonus: effect.Bonus, BattlesRemaining: effect.Battles,
	})
	return &m.boosts[len(m.boosts)-1], nil
}

func (m *mockCatalogRepository) ActiveBoosts(ctx context.Context, userID int) ([]ActiveBoost, error) {
	if m.activeBoostsErr != nil {
		return nil, m.activeBoostsErr
	}
	return m.boosts, nil
}

func (m *mockCatalogRepository) ConsumeBattleBoosts(ctx context.Context, userID int) error {
	m.tickedBattles = append(m.tickedBattles, userID)
	kept := m.boosts[:0]
	for _, b := range m.boosts {
		b.BattlesRemaining--
		if b.BattlesRemaining > 0 {
			kept = append(kept, b)
		}
	}
	m.boosts = kept
	return nil
}

func newTestService() (*Service, *mockCatalogRepository) {
	repo := &mockCatalogRepository{
		items: map[string]*Item{
			"thunder-stone": {ID: "thunder-stone", Name: "Thunder Stone", ItemType: TypeEvolution, Price: 800},
			"fire-stone":    {ID: "fire-stone", Name: "Fire Stone", ItemType: TypeEvolution, Price: 800},
			"attack-booster": {ID: "attack-booster", Name: "Attack Booster", ItemType: TypeBooster, Price: 300,
				Effect: []byte(`{"stat":"attack","bonus":5,"battles":3}`)},
		},
		quantities: map[int]map[string]int{},
	}
	return NewService(repo), repo
}

func TestPurchaseItem(t *testing.T) {
	t.Run("purchase increases inventory quantity", func(t *testing.T) {
		svc, _ := newTestService()

		item, newQty, err := svc.Purchase(context.Background(), 1, "thunder-stone", 1)
		require.NoError(t, err)
		assert.Equal(t, "thunder-stone", item.ID)
		assert.Equal(t, 1, newQty)

		qty, err := svc.GetQuantity(context.Background(), 1, "thunder-stone")
		require.NoError(t, err)
		assert.Equal(t, 1, qty)
	})

	t.Run("multiple purchases accumulate quantity", func(t *testing.T) {
		svc, _ := newTestService()

		_, newQty, err := svc.Purchase(context.Background(), 1, "thunder-stone", 2)
		require.NoError(t, err)
		assert.Equal(t, 2, newQty)

		_, newQty, err = svc.Purchase(context.Background(), 1, "thunder-stone", 1)
		require.NoError(t, err)
		assert.Equal(t, 3, newQty)
	})

	t.Run("omitted quantity defaults to one", func(t *testing.T) {
		svc, repo := newTestService()

		_, newQty, err := svc.Purchase(context.Background(), 1, "fire-stone", 0)
		require.NoError(t, err)
		assert.Equal(t, 1, newQty)
		require.Len(t, repo.purchases, 1)
		assert.Equal(t, 1, repo.purchases[0].quantity)
	})

	t.Run("invalid quantities are rejected", func(t *testing.T) {
		svc, repo := newTestService()

		for _, qty := range []int{-1, 100, 1000} {
			_, _, err := svc.Purchase(context.Background(), 1, "thunder-stone", qty)
			require.Error(t, err, "quantity %d must be rejected", qty)
		}
		assert.Empty(t, repo.purchases, "no purchase may reach the repository")
	})

	t.Run("unknown item is rejected", func(t *testing.T) {
		svc, repo := newTestService()

		_, _, err := svc.Purchase(context.Background(), 1, "moon-stone", 1)
		require.ErrorIs(t, err, ErrItemNotFound)
		assert.Empty(t, repo.purchases)
	})

	t.Run("insufficient coins surfaces the repository error", func(t *testing.T) {
		svc, repo := newTestService()
		repo.purchaseErr = ErrInsufficientCoins

		_, _, err := svc.Purchase(context.Background(), 1, "thunder-stone", 1)
		require.ErrorIs(t, err, ErrInsufficientCoins)
	})
}

func TestGetInventory(t *testing.T) {
	svc, repo := newTestService()
	repo.quantities[1] = map[string]int{
		"thunder-stone": 3,
		"fire-stone":    0, // exhausted items are not listed
	}

	inventory, err := svc.GetInventory(context.Background(), 1)
	require.NoError(t, err)

	require.Len(t, inventory, 1)
	assert.Equal(t, "thunder-stone", inventory[0].ID)
	assert.Equal(t, 3, inventory[0].Quantity)

	empty, err := svc.GetInventory(context.Background(), 2)
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestErrorSentinels(t *testing.T) {
	// Errors are wrapped with context at call sites; sentinel matching must
	// still work for handlers.
	wrapped := errors.Join(ErrInsufficientItem, errors.New("detail"))
	assert.ErrorIs(t, wrapped, ErrInsufficientItem)
}

func TestUseItem(t *testing.T) {
	t.Run("activation returns committed boost without reloading active boosts", func(t *testing.T) {
		svc, repo := newTestService()
		repo.quantities[1] = map[string]int{"attack-booster": 1}
		repo.activeBoostsErr = errors.New("active boosts unavailable")
		boost, err := svc.UseItem(context.Background(), 1, "attack-booster")
		require.NoError(t, err)
		assert.Equal(t, "attack-booster", boost.ItemID)
		assert.Zero(t, repo.quantities[1]["attack-booster"])
	})
	t.Run("using a booster activates it and consumes one from inventory", func(t *testing.T) {
		svc, repo := newTestService()
		repo.quantities[1] = map[string]int{"attack-booster": 2}

		boost, err := svc.UseItem(context.Background(), 1, "attack-booster")
		require.NoError(t, err)
		assert.Equal(t, "attack-booster", boost.ItemID)
		assert.Equal(t, StatAttack, boost.Stat)
		assert.Equal(t, 5, boost.Bonus)
		assert.Equal(t, 3, boost.BattlesRemaining)

		assert.Equal(t, 1, repo.quantities[1]["attack-booster"], "one booster must be consumed")
		require.Len(t, repo.boosts, 1)
	})

	t.Run("evolution items are rejected", func(t *testing.T) {
		svc, repo := newTestService()
		repo.quantities[1] = map[string]int{"thunder-stone": 1}

		_, err := svc.UseItem(context.Background(), 1, "thunder-stone")
		require.ErrorIs(t, err, ErrNotBooster)
		assert.Empty(t, repo.activated, "nothing may be activated")
		assert.Equal(t, 1, repo.quantities[1]["thunder-stone"], "item must not be consumed")
	})

	t.Run("zero inventory is rejected", func(t *testing.T) {
		svc, repo := newTestService()
		repo.quantities[1] = map[string]int{"attack-booster": 0}

		_, err := svc.UseItem(context.Background(), 1, "attack-booster")
		require.ErrorIs(t, err, ErrInsufficientItem)
		assert.Empty(t, repo.activated)
	})

	t.Run("unknown item is rejected", func(t *testing.T) {
		svc, repo := newTestService()

		_, err := svc.UseItem(context.Background(), 1, "moon-booster")
		require.ErrorIs(t, err, ErrItemNotFound)
		assert.Empty(t, repo.activated)
	})

	t.Run("invalid effect payloads are rejected", func(t *testing.T) {
		tests := []struct {
			name   string
			effect string
		}{
			{"missing effect", ``},
			{"unsupported stat", `{"stat":"speed","bonus":5,"battles":3}`},
			{"non-positive bonus", `{"stat":"hp","bonus":0,"battles":3}`},
			{"non-positive battles", `{"stat":"hp","bonus":10,"battles":0}`},
			{"malformed json", `{"stat":`},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				item := &Item{ID: "bad-booster", Name: "Bad", ItemType: TypeBooster, Effect: []byte(tc.effect)}
				_, err := ParseBoosterEffect(item)
				require.Error(t, err)
			})
		}
	})
}

func TestConsumeBattleBoosts(t *testing.T) {
	svc, repo := newTestService()
	repo.boosts = []ActiveBoost{
		{ID: 1, ItemID: "attack-booster", Stat: StatAttack, Bonus: 5, BattlesRemaining: 3},
		{ID: 2, ItemID: "hp-booster", Stat: StatHP, Bonus: 10, BattlesRemaining: 1},
	}

	require.NoError(t, svc.ConsumeBattleBoosts(context.Background(), 1))

	require.Len(t, repo.boosts, 1, "exhausted boost is removed")
	assert.Equal(t, 2, repo.boosts[0].BattlesRemaining)
	assert.Equal(t, []int{1}, repo.tickedBattles)
}

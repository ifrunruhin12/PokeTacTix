package shop

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"pokemon-cli/internal/items"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockItemsRepo backs a real items.Service so the shop handler test exercises
// genuine purchase validation instead of duplicating it in a mock.
type mockItemsRepo struct {
	items map[string]*items.Item
}

func (m *mockItemsRepo) ListItems(ctx context.Context) ([]items.Item, error) {
	// Deterministic order mirroring the SQL query (price ASC, id ASC) so
	// order-sensitive assertions are stable.
	ids := make([]string, 0, len(m.items))
	for id := range m.items {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	result := make([]items.Item, 0, len(ids))
	for _, id := range ids {
		result = append(result, *m.items[id])
	}
	return result, nil
}

func (m *mockItemsRepo) GetItem(ctx context.Context, itemID string) (*items.Item, error) {
	if item, ok := m.items[itemID]; ok {
		return item, nil
	}
	return nil, items.ErrItemNotFound
}

func (m *mockItemsRepo) GetInventory(ctx context.Context, userID int) ([]items.InventoryEntry, error) {
	return nil, nil
}

func (m *mockItemsRepo) GetQuantity(ctx context.Context, userID int, itemID string) (int, error) {
	return 0, nil
}

func (m *mockItemsRepo) Purchase(ctx context.Context, userID int, item *items.Item, quantity int) (int, error) {
	return 0, items.ErrInsufficientCoins
}

func (m *mockItemsRepo) ActivateBooster(ctx context.Context, userID int, item *items.Item, effect items.BoosterEffect) error {
	return nil
}

func (m *mockItemsRepo) ActiveBoosts(ctx context.Context, userID int) ([]items.ActiveBoost, error) {
	return nil, nil
}

func (m *mockItemsRepo) ConsumeBattleBoosts(ctx context.Context, userID int) error {
	return nil
}

func newTestItems() map[string]*items.Item {
	return map[string]*items.Item{
		"thunder-stone": {ID: "thunder-stone", Name: "Thunder Stone", ItemType: items.TypeEvolution, Price: 800},
		"fire-stone":    {ID: "fire-stone", Name: "Fire Stone", ItemType: items.TypeEvolution, Price: 800},
		"water-stone":   {ID: "water-stone", Name: "Water Stone", ItemType: items.TypeEvolution, Price: 800},
	}
}

func newTestApp(t *testing.T, catalog ItemCatalog, withUser bool) *fiber.App {
	t.Helper()

	app := fiber.New()
	handler := NewHandler(NewService(), nil, nil, catalog)
	RegisterRoutes(app, handler, func(c *fiber.Ctx) error {
		if withUser {
			c.Locals("user_id", 1)
		}
		return c.Next()
	})
	return app
}

func TestGetInventoryIncludesItemCategory(t *testing.T) {
	// Unauthenticated inventory read: middleware sets no user_id, matching the
	// public shop browsing path.
	app := newTestApp(t, items.NewService(&mockItemsRepo{items: newTestItems()}), false)

	resp, err := app.Test(httptest.NewRequest("GET", "/api/shop/inventory", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var payload struct {
		Items      []ShopItem `json:"items"`
		GameTokens any        `json:"game_tokens"`
		GameItems  []struct {
			ID       string `json:"id"`
			ItemType string `json:"item_type"`
			Price    int    `json:"price"`
		} `json:"game_items"`
	}
	require.NoError(t, json.Unmarshal(body, &payload))

	// The existing Pokemon category is intact.
	assert.NotNil(t, payload.Items, "pokemon category must keep working")

	// The new item category lists the catalog (sorted by id).
	require.Len(t, payload.GameItems, 3)
	assert.Equal(t, "fire-stone", payload.GameItems[0].ID)
	assert.Equal(t, items.TypeEvolution, payload.GameItems[0].ItemType)
	assert.Equal(t, 800, payload.GameItems[0].Price)
	assert.Contains(t, []string{"thunder-stone", "water-stone"}, payload.GameItems[1].ID)
	assert.Contains(t, []string{"thunder-stone", "water-stone"}, payload.GameItems[2].ID)
}

func TestGetInventoryItemCatalogFailure(t *testing.T) {
	repo := &mockItemsRepo{items: newTestItems()}
	failing := &failingItemsRepo{mockItemsRepo: *repo}
	app := newTestApp(t, items.NewService(failing), false)

	resp, err := app.Test(httptest.NewRequest("GET", "/api/shop/inventory", nil))
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

type failingItemsRepo struct {
	mockItemsRepo
}

func (m *failingItemsRepo) ListItems(ctx context.Context) ([]items.Item, error) {
	return nil, context.Canceled
}

func TestPurchaseItemValidation(t *testing.T) {
	app := newTestApp(t, items.NewService(&mockItemsRepo{items: newTestItems()}), true)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{name: "malformed body", body: `not-json`, wantStatus: fiber.StatusBadRequest},
		{name: "missing item id", body: `{}`, wantStatus: fiber.StatusBadRequest},
		{name: "blank item id", body: `{"item_id":"  "}`, wantStatus: fiber.StatusBadRequest},
		{name: "unknown item", body: `{"item_id":"moon-stone"}`, wantStatus: fiber.StatusNotFound},
		{name: "invalid quantity", body: `{"item_id":"thunder-stone","quantity":1000}`, wantStatus: fiber.StatusBadRequest},
		{name: "insufficient coins", body: `{"item_id":"thunder-stone"}`, wantStatus: fiber.StatusPaymentRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/shop/items/purchase", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

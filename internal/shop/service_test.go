package shop

import (
	"testing"
	"time"
)

// newTestService builds a Service with a deterministic inventory, avoiding
// NewService's random generation so price assertions are stable.
func newTestService() *Service {
	return &Service{
		inventory: &ShopInventory{
			Items: []ShopItem{
				{PokemonName: "Pikachu", Price: 100, Rarity: "common", InStock: true, Types: []string{"electric"}},
				{PokemonName: "Moltres", Price: 2500, Rarity: "legendary", IsLegendary: true, InStock: true, Types: []string{"fire", "flying"}},
			},
		},
		lastRefresh: time.Now(),
	}
}

// GetInventory must hand out copies: the handler annotates item prices with
// per-request discount pricing, and writing those into the shared inventory
// compounded the discount on every request.
func TestGetInventoryDiscountDoesNotCompound(t *testing.T) {
	s := newTestService()
	if err := s.ApplyDiscount(30, time.Hour); err != nil {
		t.Fatalf("ApplyDiscount failed: %v", err)
	}

	// Simulate what the inventory handler does on every request
	for request := 0; request < 3; request++ {
		inventory := s.GetInventory()
		for i := range inventory.Items {
			inventory.Items[i].Price = s.GetItemPrice(inventory.Items[i])
		}
		if got := inventory.Items[1].Price; got != 1500 {
			t.Errorf("request %d: discounted legendary price = %d, want 1500", request, got)
		}
	}

	// The shared inventory must still hold the undiscounted base price
	fresh := s.GetInventory()
	if got := fresh.Items[1].Price; got != 2500 {
		t.Errorf("base price after repeated requests = %d, want 2500 (discount compounded into shared state)", got)
	}
}

func TestGetInventoryDiscountExpires(t *testing.T) {
	s := newTestService()
	if err := s.ApplyDiscount(30, time.Hour); err != nil {
		t.Fatalf("ApplyDiscount failed: %v", err)
	}
	s.discountEndTime = time.Now().Add(-time.Minute)

	inventory := s.GetInventory()
	if inventory.DiscountActive {
		t.Error("discount still active after end time passed")
	}
	if got := s.GetItemPrice(inventory.Items[1]); got != 2500 {
		t.Errorf("legendary price after discount expiry = %d, want 2500", got)
	}
}

func TestGetInventorySnapshotIsolation(t *testing.T) {
	s := newTestService()

	first := s.GetInventory()
	first.Items[0].Price = 1
	first.Items[0].Types[0] = "corrupted"
	first.GameTokens = &GameTokenItem{Price: 1}

	second := s.GetInventory()
	if got := second.Items[0].Price; got != 100 {
		t.Errorf("price leaked between snapshots: got %d, want 100", got)
	}
	if got := second.Items[0].Types[0]; got != "electric" {
		t.Errorf("types leaked between snapshots: got %q, want %q", got, "electric")
	}
	if second.GameTokens != nil {
		t.Error("game tokens leaked between snapshots")
	}
}

func TestApplyDiscountValidatesPercent(t *testing.T) {
	s := newTestService()

	if err := s.ApplyDiscount(-1, time.Hour); err == nil {
		t.Error("negative percent accepted")
	}
	if err := s.ApplyDiscount(101, time.Hour); err == nil {
		t.Error("percent above 100 accepted")
	}
	if err := s.ApplyDiscount(50, time.Hour); err != nil {
		t.Errorf("valid percent rejected: %v", err)
	}
}

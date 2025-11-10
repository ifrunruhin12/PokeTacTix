package shop

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"pokemon-cli/internal/pokemon"
)

// Service handles shop business logic
type Service struct {
	inventory       *ShopInventory
	mu              sync.RWMutex
	lastRefresh     time.Time
	discountEndTime time.Time
}

// NewService creates a new shop service
func NewService() *Service {
	s := &Service{
		inventory: &ShopInventory{
			Items:           []ShopItem{},
			DiscountActive:  false,
			DiscountPercent: 0,
			RefreshTime:     time.Now().Add(24 * time.Hour),
		},
		lastRefresh: time.Now(),
	}
	
	// Generate initial inventory
	s.generateInventory()
	
	return s
}

// GetInventory returns the current shop inventory
func (s *Service) GetInventory() *ShopInventory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Check if inventory needs refresh (every 24 hours)
	if time.Since(s.lastRefresh) >= 24*time.Hour {
		s.mu.RUnlock()
		s.mu.Lock()
		s.generateInventory()
		s.lastRefresh = time.Now()
		s.inventory.RefreshTime = time.Now().Add(24 * time.Hour)
		s.mu.Unlock()
		s.mu.RLock()
	}
	
	// Check if discount has expired
	if s.inventory.DiscountActive && time.Now().After(s.discountEndTime) {
		s.mu.RUnlock()
		s.mu.Lock()
		s.inventory.DiscountActive = false
		s.inventory.DiscountPercent = 0
		s.mu.Unlock()
		s.mu.RLock()
	}
	
	return s.inventory
}

// generateInventory creates a new shop inventory
func (s *Service) generateInventory() {
	items := []ShopItem{}
	
	// Add 10-15 common/uncommon Pokemon
	commonUncommonCount := 10 + rand.Intn(6) // 10-15
	for i := 0; i < commonUncommonCount; i++ {
		// Random Pokemon ID from Gen 1-8 (excluding legendaries)
		pokemonID := rand.Intn(898) + 1
		pokemonName := fmt.Sprintf("%d", pokemonID)
		
		poke, moves, err := pokemon.FetchPokemon(pokemonName)
		if err != nil {
			continue
		}
		
		// Skip if legendary or mythical
		isLegendary, isMythical := pokemon.IsLegendaryOrMythical(poke.Name)
		if isLegendary || isMythical {
			continue
		}
		
		card := pokemon.BuildCardFromPokemon(poke, moves)
		
		// Determine rarity and price based on base stats
		totalStats := card.Attack + card.Defense + card.Speed + card.HPMax
		var rarity string
		var price int
		
		if totalStats < 300 {
			rarity = "common"
			price = 100
		} else if totalStats < 400 {
			rarity = "uncommon"
			price = 250
		} else {
			rarity = "rare"
			price = 500
		}
		
		items = append(items, ShopItem{
			PokemonName: card.Name,
			Price:       price,
			Rarity:      rarity,
			IsLegendary: false,
			IsMythical:  false,
			Sprite:      card.Sprite,
			Types:       card.Types,
			InStock:     true,
			BaseHP:      card.HPMax,
			BaseAttack:  card.Attack,
			BaseDefense: card.Defense,
			BaseSpeed:   card.Speed,
			Moves:       card.Moves,
		})
	}
	
	// Add 5-8 rare Pokemon
	rareCount := 5 + rand.Intn(4) // 5-8
	for i := 0; i < rareCount; i++ {
		pokemonID := rand.Intn(898) + 1
		pokemonName := fmt.Sprintf("%d", pokemonID)
		
		poke, moves, err := pokemon.FetchPokemon(pokemonName)
		if err != nil {
			continue
		}
		
		// Skip if legendary or mythical
		isLegendary, isMythical := pokemon.IsLegendaryOrMythical(poke.Name)
		if isLegendary || isMythical {
			continue
		}
		
		card := pokemon.BuildCardFromPokemon(poke, moves)
		
		items = append(items, ShopItem{
			PokemonName: card.Name,
			Price:       500,
			Rarity:      "rare",
			IsLegendary: false,
			IsMythical:  false,
			Sprite:      card.Sprite,
			Types:       card.Types,
			InStock:     true,
			BaseHP:      card.HPMax,
			BaseAttack:  card.Attack,
			BaseDefense: card.Defense,
			BaseSpeed:   card.Speed,
			Moves:       card.Moves,
		})
	}
	
	// 15% chance to include 1-2 legendary/mythical Pokemon
	if rand.Float64() < 0.15 {
		specialCount := 1 + rand.Intn(2) // 1-2
		
		for i := 0; i < specialCount; i++ {
			// 50/50 chance between legendary and mythical
			var pokemonName string
			var isLegendary, isMythical bool
			var price int
			var rarity string
			
			if rand.Float64() < 0.5 {
				// Legendary
				legendaryList := []string{
					"articuno", "zapdos", "moltres", "mewtwo", "raikou", "entei", "suicune", "lugia", "ho-oh",
					"regirock", "regice", "registeel", "latias", "latios", "kyogre", "groudon", "rayquaza",
				}
				pokemonName = legendaryList[rand.Intn(len(legendaryList))]
				isLegendary = true
				isMythical = false
				price = 2500
				rarity = "legendary"
			} else {
				// Mythical
				mythicalList := []string{
					"mew", "celebi", "jirachi", "deoxys", "manaphy", "darkrai", "shaymin",
					"victini", "keldeo", "meloetta", "genesect", "diancie", "hoopa", "volcanion",
				}
				pokemonName = mythicalList[rand.Intn(len(mythicalList))]
				isLegendary = false
				isMythical = true
				price = 5000
				rarity = "mythical"
			}
			
			poke, moves, err := pokemon.FetchPokemon(pokemonName)
			if err != nil {
				continue
			}
			
			card := pokemon.BuildCardFromPokemon(poke, moves)
			
			items = append(items, ShopItem{
				PokemonName: card.Name,
				Price:       price,
				Rarity:      rarity,
				IsLegendary: isLegendary,
				IsMythical:  isMythical,
				Sprite:      card.Sprite,
				Types:       card.Types,
				InStock:     true,
				BaseHP:      card.HPMax,
				BaseAttack:  card.Attack,
				BaseDefense: card.Defense,
				BaseSpeed:   card.Speed,
				Moves:       card.Moves,
			})
		}
	}
	
	s.inventory.Items = items
}

// ApplyDiscount activates a discount event
func (s *Service) ApplyDiscount(percent int, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if percent < 0 || percent > 100 {
		return fmt.Errorf("discount percent must be between 0 and 100")
	}
	
	s.inventory.DiscountActive = true
	s.inventory.DiscountPercent = percent
	s.discountEndTime = time.Now().Add(duration)
	
	return nil
}

// GetItemPrice returns the price of an item, applying discounts if active
func (s *Service) GetItemPrice(item ShopItem) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	price := item.Price
	
	if s.inventory.DiscountActive {
		// Apply 40% discount to legendary, 30% to mythical
		if item.IsLegendary {
			price = int(float64(price) * 0.6) // 40% off
		} else if item.IsMythical {
			price = int(float64(price) * 0.7) // 30% off
		}
	}
	
	return price
}

// FindItem finds a shop item by Pokemon name
func (s *Service) FindItem(pokemonName string) (*ShopItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, item := range s.inventory.Items {
		if item.PokemonName == pokemonName {
			return &item, nil
		}
	}
	
	return nil, fmt.Errorf("pokemon not found in shop inventory")
}

// RefreshInventory manually refreshes the shop inventory
func (s *Service) RefreshInventory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.generateInventory()
	s.lastRefresh = time.Now()
	s.inventory.RefreshTime = time.Now().Add(24 * time.Hour)
}

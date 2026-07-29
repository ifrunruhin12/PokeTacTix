package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"pokemon-cli/internal/database"
	"pokemon-cli/internal/pokemon"
	"pokemon-cli/pkg/config"
)

// Generation ranges for Pokémon national dex IDs
var genRanges = map[int][2]int{
	1: {1, 151},
	2: {152, 251},
	3: {252, 386},
	4: {387, 493},
	5: {494, 649},
	6: {650, 721},
	7: {722, 809},
	8: {810, 905},
	9: {906, 1025},
}

func main() {
	log.Println("🌱 PokeTacTix Seed CLI starting...")

	cfg := config.Load()

	// Initialize Database Pool
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("❌ Failed to initialize Database: %v", err)
	}
	defer database.CloseDB()
	dbPool := database.GetDB()

	// Initialize Redis Client
	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Printf("⚠️ Redis connection failed (%v). Seed will populate PostgreSQL only.", err)
	} else {
		defer database.CloseRedis()
		log.Println("✅ Connected to Redis cache")
	}
	redisClient := database.GetRedis()

	// Create dependencies
	repo := pokemon.NewPostgresRepository(dbPool)
	var cache pokemon.Cache
	if redisClient != nil {
		cache = pokemon.NewRedisCache(redisClient, cfg.Redis.TTL)
	}
	pokeClient := pokemon.NewPokeAPIClient(os.Getenv("POKEAPI_BASE_URL"), 10*time.Second)

	pokeService := pokemon.NewService(cache, repo, pokeClient)

	// Determine targeted IDs from SEED_GENERATIONS or default (Gen 1-3)
	targetIDs := determineTargetIDs()
	log.Printf("📦 Preparing to seed %d Pokémon into durable cache...", len(targetIDs))

	// Worker pool settings
	workerCount := 6
	if wStr := os.Getenv("SEED_WORKERS"); wStr != "" {
		if w, err := strconv.Atoi(wStr); err == nil && w > 0 {
			workerCount = w
		}
	}

	jobs := make(chan int, len(targetIDs))
	for _, id := range targetIDs {
		jobs <- id
	}
	close(jobs)

	var wg sync.WaitGroup
	var successCount int64
	var errorCount int64

	startTime := time.Now()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for id := range jobs {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				p, err := pokeService.GetByID(ctx, id)
				cancel()

				if err != nil {
					atomic.AddInt64(&errorCount, 1)
					log.Printf("[Worker %d] ❌ Failed to fetch/seed ID %d: %v", workerID, id, err)
				} else {
					atomic.AddInt64(&successCount, 1)
					if successCount%25 == 0 {
						log.Printf("PROGRESS: %d/%d seeded (%s)", successCount, len(targetIDs), p.Name)
					}
				}
				time.Sleep(50 * time.Millisecond) // Throttle to be polite to PokéAPI
			}
		}(i + 1)
	}

	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("🎉 Seeding completed in %s! Successfully seeded: %d, Errors: %d", duration, successCount, errorCount)
}

func determineTargetIDs() []int {
	genEnv := os.Getenv("SEED_GENERATIONS")
	if genEnv == "" {
		genEnv = "1,2,3" // Default to Gen 1-3
	}

	idMap := make(map[int]bool)
	parts := strings.Split(genEnv, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if genNum, err := strconv.Atoi(part); err == nil {
			if r, ok := genRanges[genNum]; ok {
				for id := r[0]; id <= r[1]; id++ {
					idMap[id] = true
				}
			}
		}
	}

	// Also support SEED_MAX_ID override
	if maxStr := os.Getenv("SEED_MAX_ID"); maxStr != "" {
		if maxID, err := strconv.Atoi(maxStr); err == nil && maxID > 0 {
			for id := 1; id <= maxID; id++ {
				idMap[id] = true
			}
		}
	}

	var targetIDs []int
	for id := 1; id <= 1025; id++ {
		if idMap[id] {
			targetIDs = append(targetIDs, id)
		}
	}

	if len(targetIDs) == 0 {
		// Fallback to Gen 1 (1-151)
		for id := 1; id <= 151; id++ {
			targetIDs = append(targetIDs, id)
		}
	}

	return targetIDs
}

package battle

import (
	"context"
	"fmt"

	"pokemon-cli/internal/pokemon"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type History interface {
	RecentPokemon(ctx context.Context, playerID int, n int) ([]int, error)
	RecentFamilies(ctx context.Context, playerID int, n int) ([]int, error)
	RecordEncounter(ctx context.Context, playerID int, p *pokemon.Pokemon) error
}

type historyManager struct {
	redisClient *redis.Client
	dbPool      *pgxpool.Pool
	maxWindow   int
}

func NewHistoryManager(redisClient *redis.Client, dbPool *pgxpool.Pool, maxWindow int) History {
	if maxWindow <= 0 {
		maxWindow = 10 // Safe default max window
	}
	return &historyManager{
		redisClient: redisClient,
		dbPool:      dbPool,
		maxWindow:   maxWindow,
	}
}

func (h *historyManager) RecentPokemon(ctx context.Context, playerID int, n int) ([]int, error) {
	if n <= 0 {
		return nil, nil
	}

	key := fmt.Sprintf("player:%d:recent_pokemon", playerID)

	// 1. Check Redis
	if h.redisClient != nil {
		vals, err := h.redisClient.LRange(ctx, key, 0, int64(n-1)).Result()
		if err == nil && len(vals) > 0 {
			var result []int
			for _, val := range vals {
				var id int
				if _, fmtErr := fmt.Sscanf(val, "%d", &id); fmtErr == nil {
					result = append(result, id)
				}
			}
			return result, nil
		}
	}

	// 2. Fallback to PostgreSQL battle_encounter_log
	if h.dbPool == nil {
		return nil, nil
	}

	query := `
		SELECT pokemon_id
		FROM battle_encounter_log
		WHERE player_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := h.dbPool.Query(ctx, query, playerID, n)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent pokemon history from db: %w", err)
	}
	defer rows.Close()

	var result []int
	for rows.Next() {
		var pID int
		if err := rows.Scan(&pID); err == nil {
			result = append(result, pID)
		}
	}

	// Lazy backfill into Redis if we found results
	if h.redisClient != nil && len(result) > 0 {
		pipe := h.redisClient.Pipeline()
		for i := len(result) - 1; i >= 0; i-- {
			pipe.LPush(ctx, key, result[i])
		}
		pipe.LTrim(ctx, key, 0, int64(h.maxWindow-1))
		_, _ = pipe.Exec(ctx)
	}

	return result, rows.Err()
}

func (h *historyManager) RecentFamilies(ctx context.Context, playerID int, n int) ([]int, error) {
	if n <= 0 {
		return nil, nil
	}

	key := fmt.Sprintf("player:%d:recent_families", playerID)

	// 1. Check Redis
	if h.redisClient != nil {
		vals, err := h.redisClient.LRange(ctx, key, 0, int64(n-1)).Result()
		if err == nil && len(vals) > 0 {
			var result []int
			for _, val := range vals {
				var id int
				if _, fmtErr := fmt.Sscanf(val, "%d", &id); fmtErr == nil {
					result = append(result, id)
				}
			}
			return result, nil
		}
	}

	// 2. Fallback to PostgreSQL
	if h.dbPool == nil {
		return nil, nil
	}

	query := `
		SELECT evolution_chain_id
		FROM battle_encounter_log
		WHERE player_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := h.dbPool.Query(ctx, query, playerID, n)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent families history from db: %w", err)
	}
	defer rows.Close()

	var result []int
	for rows.Next() {
		var fID int
		if err := rows.Scan(&fID); err == nil {
			result = append(result, fID)
		}
	}

	// Lazy backfill into Redis
	if h.redisClient != nil && len(result) > 0 {
		pipe := h.redisClient.Pipeline()
		for i := len(result) - 1; i >= 0; i-- {
			pipe.LPush(ctx, key, result[i])
		}
		pipe.LTrim(ctx, key, 0, int64(h.maxWindow-1))
		_, _ = pipe.Exec(ctx)
	}

	return result, rows.Err()
}

func (h *historyManager) RecordEncounter(ctx context.Context, playerID int, p *pokemon.Pokemon) error {
	if p == nil {
		return nil
	}

	// 1. Record into Redis lists
	if h.redisClient != nil {
		pokeKey := fmt.Sprintf("player:%d:recent_pokemon", playerID)
		famKey := fmt.Sprintf("player:%d:recent_families", playerID)

		pipe := h.redisClient.Pipeline()
		pipe.LPush(ctx, pokeKey, p.ID)
		pipe.LTrim(ctx, pokeKey, 0, int64(h.maxWindow-1))

		pipe.LPush(ctx, famKey, p.EvolutionChainID)
		pipe.LTrim(ctx, famKey, 0, int64(h.maxWindow-1))

		_, _ = pipe.Exec(ctx)
	}

	// 2. Record into PostgreSQL audit log
	if h.dbPool != nil {
		query := `
			INSERT INTO battle_encounter_log (player_id, pokemon_id, evolution_chain_id)
			VALUES ($1, $2, $3)
		`
		_, err := h.dbPool.Exec(ctx, query, playerID, p.ID, p.EvolutionChainID)
		if err != nil {
			return fmt.Errorf("failed to insert battle_encounter_log: %w", err)
		}
	}

	return nil
}

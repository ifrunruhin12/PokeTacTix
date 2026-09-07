package pokemon

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetPokemon(ctx context.Context, id int) (*Pokemon, error)
	SetPokemon(ctx context.Context, p *Pokemon) error
	GetIDByName(ctx context.Context, name string) (int, error)
	SetIDByName(ctx context.Context, name string, id int) error
	GetEvolutionChain(ctx context.Context, id int) (*EvolutionChain, error)
	SetEvolutionChain(ctx context.Context, ec *EvolutionChain) error
}

type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) Cache {
	if ttl <= 0 {
		ttl = 720 * time.Hour // 30 days default
	}
	return &redisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *redisCache) GetPokemon(ctx context.Context, id int) (*Pokemon, error) {
	if c.client == nil {
		return nil, redis.Nil
	}
	key := fmt.Sprintf("pokemon:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var p Pokemon
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached pokemon: %w", err)
	}
	return &p, nil
}

func (c *redisCache) SetPokemon(ctx context.Context, p *Pokemon) error {
	if c.client == nil || p == nil {
		return nil
	}
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal pokemon for cache: %w", err)
	}

	key := fmt.Sprintf("pokemon:%d", p.ID)
	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set redis pokemon cache: %w", err)
	}

	// Also index name -> id
	_ = c.SetIDByName(ctx, p.Name, p.ID)
	return nil
}

func (c *redisCache) GetIDByName(ctx context.Context, name string) (int, error) {
	if c.client == nil {
		return 0, redis.Nil
	}
	key := fmt.Sprintf("pokemon:name:%s", strings.ToLower(name))
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (c *redisCache) SetIDByName(ctx context.Context, name string, id int) error {
	if c.client == nil {
		return nil
	}
	key := fmt.Sprintf("pokemon:name:%s", strings.ToLower(name))
	return c.client.Set(ctx, key, id, c.ttl).Err()
}

func (c *redisCache) GetEvolutionChain(ctx context.Context, id int) (*EvolutionChain, error) {
	if c.client == nil {
		return nil, redis.Nil
	}
	key := fmt.Sprintf("evochain:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var ec EvolutionChain
	if err := json.Unmarshal([]byte(val), &ec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached evolution chain: %w", err)
	}
	return &ec, nil
}

func (c *redisCache) SetEvolutionChain(ctx context.Context, ec *EvolutionChain) error {
	if c.client == nil || ec == nil {
		return nil
	}
	data, err := json.Marshal(ec)
	if err != nil {
		return fmt.Errorf("failed to marshal evolution chain for cache: %w", err)
	}

	key := fmt.Sprintf("evochain:%d", ec.ID)
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

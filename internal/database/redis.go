package database

import (
	"context"
	"fmt"
	"time"

	"pokemon-cli/pkg/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis initializes the Redis client
func InitRedis(cfg *config.RedisConfig) error {
	if cfg == nil || cfg.URL == "" {
		return fmt.Errorf("redis URL is not configured")
	}

	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return fmt.Errorf("unable to parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return fmt.Errorf("unable to ping redis: %w", err)
	}

	RedisClient = client
	return nil
}

// CloseRedis closes the Redis client connection
func CloseRedis() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}

// GetRedis returns the global Redis client
func GetRedis() *redis.Client {
	return RedisClient
}

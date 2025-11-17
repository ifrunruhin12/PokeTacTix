package database

import (
	"context"
	"fmt"
	"pokemon-cli/pkg/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// InitDB initializes the database connection pool
func InitDB(cfg *config.DatabaseConfig) error {
	if cfg.URL == "" {
		return fmt.Errorf("database URL is not set")
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = cfg.IdleTimeout
	poolConfig.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	DB = pool
	return nil
}

// CloseDB closes the database connection pool
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// GetDB returns the database connection pool
func GetDB() *pgxpool.Pool {
	return DB
}

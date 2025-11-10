package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	maxConns := getEnvInt("DB_MAX_CONNECTIONS", 20)
	config.MaxConns = int32(maxConns)
	
	minConns := getEnvInt("DB_MIN_CONNECTIONS", 2)
	config.MinConns = int32(minConns)
	
	idleTimeout := getEnvInt("DB_IDLE_TIMEOUT", 300)
	config.MaxConnIdleTime = time.Duration(idleTimeout) * time.Second
	
	maxLifetime := getEnvInt("DB_MAX_LIFETIME", 1800)
	config.MaxConnLifetime = time.Duration(maxLifetime) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
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

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetDB() *pgxpool.Pool {
	return DB
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

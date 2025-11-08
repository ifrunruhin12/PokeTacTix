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

// InitDB initializes the database connection pool
func InitDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Parse pool configuration
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Set connection pool settings
	maxConns := getEnvInt("DB_MAX_CONNECTIONS", 20)
	config.MaxConns = int32(maxConns)
	
	minConns := getEnvInt("DB_MIN_CONNECTIONS", 2)
	config.MinConns = int32(minConns)
	
	idleTimeout := getEnvInt("DB_IDLE_TIMEOUT", 300)
	config.MaxConnIdleTime = time.Duration(idleTimeout) * time.Second
	
	maxLifetime := getEnvInt("DB_MAX_LIFETIME", 1800)
	config.MaxConnLifetime = time.Duration(maxLifetime) * time.Second

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
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

// getEnvInt retrieves an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

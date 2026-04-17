package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"prompt-management/internal/config"
)

// NewPool creates a new PostgreSQL connection pool based on the provided configuration.
func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	// Apply Max Connections
	if cfg.DBMaxConns != "" {
		max, err := strconv.Atoi(cfg.DBMaxConns)
		if err == nil {
			poolCfg.MaxConns = int32(max)
		}
	}

	// Apply Min Connections
	if cfg.DBMinConns != "" {
		min, err := strconv.Atoi(cfg.DBMinConns)
		if err == nil {
			poolCfg.MinConns = int32(min)
		}
	}

	// Apply Max Idle Time
	if cfg.DBMaxConnIdleTime != "" {
		dur, err := time.ParseDuration(cfg.DBMaxConnIdleTime)
		if err == nil {
			poolCfg.MaxConnIdleTime = dur
		}
	}

	// Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return pool, nil
}

// Ping checks the connectivity to the database.
func Ping(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.Ping(ctx)
}

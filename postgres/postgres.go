package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns        = 10
	defaultMinConns        = 0
	defaultMaxConnLifetime = time.Hour
	defaultMaxConnIdleTime = 30 * time.Minute
	defaultHealthCheck     = 15 * time.Second
	pgConnRetryTimeout     = 30 * time.Second
)

// New creates a pgxpool with retry (exponential backoff) until the database is reachable.
func New(cfg *Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}
	maxConns := defaultMaxConns
	if cfg.MaxConns > 0 {
		maxConns = cfg.MaxConns
	}
	minConns := defaultMinConns
	if cfg.MinConns > 0 {
		minConns = cfg.MinConns
	}
	poolCfg.MaxConns = int32(maxConns)
	poolCfg.MinConns = int32(minConns)
	poolCfg.MaxConnLifetime = defaultMaxConnLifetime
	poolCfg.MaxConnIdleTime = defaultMaxConnIdleTime
	poolCfg.HealthCheckPeriod = defaultHealthCheck

	var pool *pgxpool.Pool
	operation := func() error {
		var createErr error
		pool, createErr = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if createErr != nil {
			return createErr
		}
		if pingErr := pool.Ping(context.Background()); pingErr != nil {
			pool.Close()
			pool = nil
			return pingErr
		}
		return nil
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = pgConnRetryTimeout
	if err := backoff.Retry(operation, bo); err != nil {
		return nil, fmt.Errorf("postgres - New: %w", err)
	}

	return pool, nil
}

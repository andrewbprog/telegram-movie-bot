package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	MaxCon      = 10
	MinCon      = 1
	HealthCheck = 15 * time.Second
	Timeout     = 5 * time.Second
)

func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	cfg.MaxConns = MaxCon
	cfg.MinConns = MinCon
	cfg.HealthCheckPeriod = HealthCheck

	ctx2, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx2, cfg)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	return pool, nil
}

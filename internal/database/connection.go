package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

func NewConnection(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// config.MaxConnIdleTime = time.Second
	// config.MinConns = 0
	// config.MaxConns = 10

	Conn, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return Conn, nil
}

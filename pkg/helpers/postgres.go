package helpers

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitPostgres(ctx context.Context, dsn string, migrationPath string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// migrations

	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		return nil, err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	if err := m.Up(); err != nil {
		return nil, err
	}

	return pool, nil
}

package postgres

import (
	"context"
	"fmt"
	"time"

	"payroll/internal/infrastructure/config"
	"payroll/internal/infrastructure/database/postgres/repository"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg config.DatabaseConfig) (repository.Querier, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	conn, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	conn.Config().MaxConns = 10
	conn.Config().MinConns = 2
	conn.Config().MaxConnLifetime = 5 * time.Minute
	conn.Config().MaxConnIdleTime = 3 * time.Minute
	conn.Config().HealthCheckPeriod = 1 * time.Minute

	conn.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return repository.New(conn), nil
}

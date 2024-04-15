package psql

import (
	"context"
	"fmt"
	internal "github.com/MamushevArup/typeracer/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func New(ctx context.Context, cfg *internal.Config) (*pgxpool.Pool, error) {
	pg := cfg.Postgres

	dbPasswd := os.Getenv("POSTGRES_PASSWORD")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pg.User, dbPasswd, pg.Host, pg.Port, pg.Database)

	conn, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", cfg.Postgres.Database, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	log.Println("connect to the database successfully")

	return conn, nil
}

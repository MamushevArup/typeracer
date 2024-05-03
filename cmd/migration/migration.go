package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"time"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
	dir             = "schema"
	driverName      = "postgres"
)

var (
	NoMigrationFiles = errors.New("goose: no migration files")
	NoDbConnection   = errors.New("goose: environment variable not declared: PG_URL")
)

func Run(cfg *config.Config) error {
	databaseURL := cfg.Postgres.PgUrl
	if len(databaseURL) == 0 {
		return NoDbConnection
	}
	databaseURL += "?sslmode=disable"
	var (
		attempts = defaultAttempts
		err      error
		db       *sql.DB
	)

	err = goose.SetDialect(driverName)
	if err != nil {
		return fmt.Errorf("goose: set dialect error: %w", err)
	}

	for attempts > 0 {
		db, err = sql.Open(driverName, databaseURL)
		if err == nil {
			break
		}

		log.Printf("goose: postgres db is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return fmt.Errorf("goose: postgres db connect error: %w", err)
	}

	defer func() { _ = db.Close() }()

	err = goose.Up(db, dir)
	if err != nil {
		if errors.Is(err, goose.ErrAlreadyApplied) {
			log.Println("goose: up migration applied")
			return nil
		}
		if errors.Is(err, goose.ErrNoMigrationFiles) {
			return NoMigrationFiles
		}
		return err
	}
	return nil
}

package psql

import (
	"context"
	"fmt"
	internal "github.com/MamushevArup/typeracer/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func NewDBConnector(cfg *internal.Config) *pgxpool.Pool {
	pg := cfg.Postgres
	dbPasswd := os.Getenv("POSTGRES_PASSWORD")
	fmt.Println(dbPasswd, "DATABSE PASSWORD HERRE")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pg.User, dbPasswd, pg.Host, pg.Port, pg.Database)
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("unable to connect to database, config issue %v\n", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("unable to ping database %v\n", err)
	}
	log.Println("connect to the database successfully")
	return conn
}

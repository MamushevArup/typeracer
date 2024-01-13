package psql

import (
	"context"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConnector interface {
	DBConnector() *pgxpool.Pool
}

type psql struct {
	lg  *logger.Logger
	url string
}

func NewDBConnector(lg *logger.Logger, url string) DBConnector {
	return &psql{
		lg:  lg,
		url: url,
	}
}

func (p *psql) postgresUp() *pgxpool.Pool {
	// postgres://username:password@localhost:5432/database_name
	conn, err := pgxpool.New(context.Background(), p.url)
	if err != nil {
		p.lg.Fatalf("unable to connect to database, config issue %v\n", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		p.lg.Fatalf("unable to ping database %v\n", err)
	}
	defer conn.Close()
	p.lg.Infof("connect to the database successfully")
	return conn
}

func (p *psql) DBConnector() *pgxpool.Pool {
	return p.postgresUp()
}

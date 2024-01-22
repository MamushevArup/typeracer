package repository

import (
	"github.com/MamushevArup/typeracer/internal/repository/postgres"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Starter postgres.Starter
}

func NewRepo(lg *logger.Logger, db *pgxpool.Pool) *Repo {
	return &Repo{
		Starter: postgres.NewSingle(lg, db),
	}
}

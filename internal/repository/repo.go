package repository

import (
	"github.com/MamushevArup/typeracer/internal/repository/postgres/contributor"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/single"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Starter     single.Starter
	Contributor contributor.Contributor
}

func NewRepo(lg *logger.Logger, db *pgxpool.Pool) *Repo {
	return &Repo{
		Starter:     single.NewSingle(lg, db),
		Contributor: contributor.NewContributor(lg, db),
	}
}

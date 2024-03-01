package repository

import (
	"github.com/MamushevArup/typeracer/internal/repository/postgres/authr"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/contributor"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/multiple"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/single"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Starter     single.Starter
	Contributor contributor.Contributor
	Auth        authr.Auth
	Multiple    multiple.Multiple
}

func NewRepo(lg *logger.Logger, db *pgxpool.Pool) *Repo {
	return &Repo{
		Starter:     single.NewSingle(lg, db),
		Contributor: contributor.NewContributor(lg, db),
		Auth:        authr.NewUser(db, lg),
		Multiple:    multiple.NewMultiple(lg, db),
	}
}

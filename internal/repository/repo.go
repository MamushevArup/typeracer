package repository

import (
	"github.com/MamushevArup/typeracer/internal/repository/postgres/admin"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/authr"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/contributor"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/multiple/link"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/multiple/race"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/racer"
	"github.com/MamushevArup/typeracer/internal/repository/postgres/single"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Starter     single.Starter
	Auth        authr.Auth
	Contributor contributor.Contributor
	Link        link.Manager
	Multiple    race.Multiple
	Admin       admin.Moderation
	Racer       racer.Profile
}

func NewRepo(lg *logger.Logger, db *pgxpool.Pool) *Repo {
	return &Repo{
		Starter:     single.New(lg, db),
		Contributor: contributor.NewContributor(lg, db),
		Auth:        authr.NewUser(db, lg),
		Link:        link.NewManager(db, lg),
		Multiple:    race.NewRace(lg, db),
		Admin:       admin.NewModeration(db, lg),
		Racer:       racer.New(db, lg),
	}
}

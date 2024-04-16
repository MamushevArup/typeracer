package services

import (
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/MamushevArup/typeracer/internal/services/auth"
	"github.com/MamushevArup/typeracer/internal/services/contribute"
	"github.com/MamushevArup/typeracer/internal/services/multiple/link"
	"github.com/MamushevArup/typeracer/internal/services/multiple/race"
	"github.com/MamushevArup/typeracer/internal/services/single"
)

type Service struct {
	Single     single.PracticeY
	Contribute contribute.Contributor
	Auth       auth.Auth
	Multiple   race.Racer
	Link       link.Checker
}

func NewService(repo *repository.Repo, cfg *config.Config) *Service {
	return &Service{
		Single:     single.NewPracticeY(repo),
		Contribute: contribute.NewContribute(repo),
		Auth:       auth.NewAuth(repo),
		Multiple:   race.NewMultiple(repo, cfg),
		Link:       link.NewLink(repo),
	}
}

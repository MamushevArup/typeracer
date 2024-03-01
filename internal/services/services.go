package services

import (
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/MamushevArup/typeracer/internal/services/auth"
	"github.com/MamushevArup/typeracer/internal/services/contribute"
	"github.com/MamushevArup/typeracer/internal/services/multiple"
	"github.com/MamushevArup/typeracer/internal/services/single"
)

type Service struct {
	PracticeY  single.PracticeY
	Contribute contribute.Contributor
	Auth       auth.Auth
	Multiple   multiple.Racer
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		PracticeY:  single.NewPracticeY(repo),
		Contribute: contribute.NewContribute(repo),
		Auth:       auth.NewAuth(repo),
		Multiple:   multiple.NewMultiple(repo),
	}
}

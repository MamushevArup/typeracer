package services

import (
	"github.com/MamushevArup/typeracer/adapters/avatar/aws"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/MamushevArup/typeracer/internal/services/admin"
	"github.com/MamushevArup/typeracer/internal/services/auth"
	"github.com/MamushevArup/typeracer/internal/services/contribute"
	"github.com/MamushevArup/typeracer/internal/services/multiple/link"
	"github.com/MamushevArup/typeracer/internal/services/multiple/race"
	"github.com/MamushevArup/typeracer/internal/services/racer"
	"github.com/MamushevArup/typeracer/internal/services/single"
)

type Service struct {
	Single     single.Practice
	Contribute contribute.Contributor
	Auth       auth.Auth
	Multiple   race.Racer
	Link       link.Checker
	Admin      admin.Admin
	Racer      racer.Profile
}

func NewService(repo *repository.Repo, cfg *config.Config, s3 aws.CloudService) *Service {
	return &Service{
		Single:     single.NewPracticeY(repo),
		Contribute: contribute.NewContribute(repo),
		Auth:       auth.NewAuth(repo, cfg),
		Multiple:   race.NewMultiple(repo, cfg),
		Link:       link.NewLink(repo),
		Admin:      admin.New(repo, s3),
		Racer:      racer.New(repo),
	}
}

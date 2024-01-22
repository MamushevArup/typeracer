package services

import (
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/MamushevArup/typeracer/internal/services/single"
)

type Service struct {
	PracticeY single.PracticeY
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		PracticeY: single.NewPracticeY(repo),
	}
}

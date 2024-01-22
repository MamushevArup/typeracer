package single

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
)

// this package stands for implement the practice yourself section only

type PracticeY interface {
	StartRace(ctx context.Context, userId uuid.UUID) (*models.Single, error)
}

type service struct {
	repo *repository.Repo
}

func NewPracticeY(repo *repository.Repo) PracticeY {
	return &service{
		repo: repo,
	}
}

func (s *service) StartRace(ctx context.Context, userId uuid.UUID) (*models.Single, error) {
	//TODO implement me
	single, err := s.repo.Starter.StartSingle(ctx, userId)
	if err != nil {
		return nil, errors.New("you can't start a race")
	}
	return single, nil
}

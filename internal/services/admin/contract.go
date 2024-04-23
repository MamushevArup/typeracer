package admin

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
)

type Admin interface {
	ShowContentToModerate(ctx context.Context, limit, offset, sort string) ([]models.ModerationServiceResponse, error)
	TextDetails(ctx context.Context, modId string) (models.ModerationTextDetails, error)
	ApproveContent(ctx context.Context, modId string) error
}

type service struct {
	repo *repository.Repo
}

func New(repo *repository.Repo) Admin {
	return &service{repo: repo}
}

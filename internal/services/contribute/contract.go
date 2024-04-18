package contribute

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
)

type Contributor interface {
	ContributeText(ctx context.Context, contributeDTO models.ContributeHandlerRequest) error
}

type service struct {
	repo *repository.Repo
}

func NewContribute(repo *repository.Repo) Contributor {
	return &service{
		repo: repo,
	}
}

package link

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"time"
)

type Checker interface {
	Create(ctx context.Context, id string) (uuid.UUID, error)
	Check(ctx context.Context, link string) error
	Kill(ticker *time.Ticker)
}

type link struct {
	repo *repository.Repo
}

func NewLink(repo *repository.Repo) Checker {
	return &link{
		repo: repo,
	}
}

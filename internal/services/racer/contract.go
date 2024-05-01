package racer

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
)

type Profile interface {
	Details(ctx context.Context, racerId string) (models.RacerHandler, error)
	Avatars(ctx context.Context) ([]models.Avatar, error)
	UpdateAvatar(ctx context.Context, avatar models.AvatarUpdate) error
}

type service struct {
	repo *repository.Repo
}

func New(repo *repository.Repo) Profile {
	return &service{repo: repo}
}

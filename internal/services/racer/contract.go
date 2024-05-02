package racer

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
)

type Profile interface {
	Details(ctx context.Context, racerId string) (models.RacerHandler, error)
	Avatars(ctx context.Context) ([]models.Avatar, error)
	UpdateAvatar(ctx context.Context, avatar models.AvatarUpdate) error
	UpdateRacerInfo(ctx context.Context, racer models.RacerUpdate) error
	SingleHistory(ctx context.Context, id, limit, offset string) ([]models.SingleHistoryHandler, error)
	HistorySingleText(ctx context.Context, racerId string, singleID uuid.UUID) (models.SingleHistoryText, error)
}

type service struct {
	repo *repository.Repo
}

func New(repo *repository.Repo) Profile {
	return &service{repo: repo}
}

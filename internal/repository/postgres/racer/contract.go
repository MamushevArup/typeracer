package racer

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile interface {
	Info(ctx context.Context, racerId uuid.UUID) (models.RacerRepository, error)
	SelectAvatar(ctx context.Context) ([]models.Avatar, error)
	UpdateAvatar(ctx context.Context, avatar models.AvatarUpdateRepo) error
	UpdateRacer(ctx context.Context, racer models.RacerUpdateRepo) error
	SingleHistoryRows(ctx context.Context, racerId uuid.UUID, limit, offset int) ([]models.SingleHistory, error)
	SelectHistoryText(ctx context.Context, racerUUID, singleID uuid.UUID) (models.SingleHistoryText, error)
}

type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

func New(db *pgxpool.Pool, lg *logger.Logger) Profile {
	return &repo{db: db, lg: lg}
}

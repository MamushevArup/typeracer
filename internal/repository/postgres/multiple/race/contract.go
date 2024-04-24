package race

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Multiple interface {
	Texts(ctx context.Context) ([]uuid.UUID, error)
	Text(ctx context.Context, id uuid.UUID) (string, error)
	AddRacers(ctx context.Context, mlt models.MultipleRace) error
	User(ctx context.Context, id uuid.UUID) (models.RacerM, error)
	RacerID(ctx context.Context, email string) (uuid.UUID, error)
	InsertSession(ctx context.Context, r *models.RacerRepoM) error
	UpdateRacerHistory(ctx context.Context, currSpeed int, racerID, textID uuid.UUID) error
	GuestAvatar(ctx context.Context) (string, error)
}

type multiple struct {
	lg *logger.Logger
	db *pgxpool.Pool
}

func NewRace(lg *logger.Logger, db *pgxpool.Pool) Multiple {
	return &multiple{
		lg: lg,
		db: db,
	}
}

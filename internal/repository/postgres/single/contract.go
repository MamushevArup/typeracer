package single

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Starter interface {
	RacerInfo(ctx context.Context, id uuid.UUID) (models.RacerInfo, error)
	TextInfo(ctx context.Context) (models.TextInfo, error)
	GuestAvatar(ctx context.Context) (string, error)
	EndSingleRace(ctx context.Context, req models.RespEndSingle) error
	GetTextLen(ctx context.Context, textID uuid.UUID) (int, error)
	RacerExist(ctx context.Context, id uuid.UUID) (bool, error)
}

// This is responsible for the practice yourself section
type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

// New return instance of the implemented interface
func New(lg *logger.Logger, db *pgxpool.Pool) Starter {
	return &repo{
		db: db,
		lg: lg,
	}
}

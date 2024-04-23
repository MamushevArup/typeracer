package admin

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Moderation interface {
	SelectModeration(ctx context.Context, limit, offset int, sort string) ([]models.ModerationRepoResponse, error)
	ContentDetails(ctx context.Context, modId uuid.UUID) (models.ModerationTextDetails, error)
	ModerationById(ctx context.Context, modId uuid.UUID) (models.ModerationApprove, error)
	ApproveContent(ctx context.Context, transaction models.TextAcceptTransaction) error
	DeleteModerationText(ctx context.Context, modId uuid.UUID) error
	RejectContent(ctx context.Context, reject models.ModerationRejectToRepo) error
}

type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

func NewModeration(db *pgxpool.Pool, lg *logger.Logger) Moderation {
	return &repo{
		db: db,
		lg: lg,
	}
}

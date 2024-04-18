package contributor

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Contributor interface {
	InsertToModeration(ctx context.Context, c models.ContributeServiceRequest) error
	DuplicateContent(ctx context.Context, racerID uuid.UUID, contentHash string) (bool, error)
}

type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

// NewContributor return instance of the implemented interface
func NewContributor(lg *logger.Logger, db *pgxpool.Pool) Contributor {
	return &repo{
		db: db,
		lg: lg,
	}
}

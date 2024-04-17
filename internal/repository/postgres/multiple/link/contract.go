package link

import (
	"context"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Manager interface {
	Add(ctx context.Context, link uuid.UUID, creator string, time time.Time) error
	Check(ctx context.Context, link uuid.UUID) (bool, error)
	Remove(ctx context.Context, currentTime time.Time) error
}

type link struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

func NewManager(db *pgxpool.Pool, lg *logger.Logger) Manager {
	return &link{
		db: db,
		lg: lg,
	}
}

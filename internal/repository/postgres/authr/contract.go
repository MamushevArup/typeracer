package authr

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth interface {
	GetUserPasswordByEmail(ctx context.Context, email string) (uuid.UUID, string, string, error)
	UserByEmail(ctx context.Context, email string) (bool, error)
	InsertUser(ctx context.Context, racerAuth models.RacerAuth) error
	DeleteSession(ctx context.Context, fng, refresh string) error
	InsertSession(ctx context.Context, r models.RacerAuth) error
	UserSession(ctx context.Context, token string, id uuid.UUID) (bool, error)
	Fingerprint(ctx context.Context, fng, refresh string) (models.RacerAuth, error)
	DeleteRefreshSession(ctx context.Context, refresh string) error
	UpdateAdmin(ctx context.Context, username, refresh string) error
	AdminByRefresh(ctx context.Context, refresh string) (string, error)
	UpdateRefresh(ctx context.Context, id, refreshNew string) error
}

type auth struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

func NewUser(db *pgxpool.Pool, lg *logger.Logger) Auth {
	return &auth{
		db: db,
		lg: lg,
	}
}

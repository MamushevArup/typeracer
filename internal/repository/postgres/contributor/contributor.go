package contributor

import (
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Contributor interface {
	//Contribute(ctx context.Context)
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

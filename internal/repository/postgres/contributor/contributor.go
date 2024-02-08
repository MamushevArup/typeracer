package contributor

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Contributor interface {
	Contribute(ctx context.Context, ctr models.ContributeText) error
	RacerExist(ctx context.Context, id uuid.UUID) (bool, error)
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

func (r *repo) Contribute(ctx context.Context, ctr models.ContributeText) error {
	query := "INSERT INTO moderation values($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, query, ctr.RacerID, ctr.Content, ctr.Author, ctr.Length, ctr.Source, ctr.SourceTitle, ctr.SentAt)
	if err != nil {
		r.lg.Errorf("unable to insert to moderation table due to %v", err)
		return err
	}
	return nil
}

func (r *repo) RacerExist(ctx context.Context, id uuid.UUID) (bool, error) {
	var ex bool
	query := "SELECT EXISTS(SELECT 1 FROM racer WHERE id = $1)"
	err := r.db.QueryRow(ctx, query, id).Scan(&ex)
	if err != nil {
		r.lg.Errorf("can't execute use exist check %v", err)
		return false, err
	}
	return ex, nil
}

// Moderation return all entries from moderation table
func (r *repo) Moderation(ctx context.Context) {

}

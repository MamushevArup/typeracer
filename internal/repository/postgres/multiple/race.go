package multiple

import (
	"context"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Multiple interface {
	InsertLink(ctx context.Context, link uuid.UUID, creator string, time time.Time) error
	Link(ctx context.Context, link uuid.UUID, time time.Time) (bool, error)
	//DeleteLink(ctx context.Context, currentTime time.Time) error
}

//func (r *repo) DeleteLink(ctx context.Context, currentTime time.Time) error {
//	query := `DELETE FROM link_management`
//}

func (r *repo) Link(ctx context.Context, link uuid.UUID, t time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM link_management WHERE link=$1 AND created_at>=$2 AND created_at<=NOW()`

	var count int

	err := r.db.QueryRow(ctx, query, link, t).Scan(&count)
	if err != nil {
		return false, err
	}
	// if true link exist we pass user further if no close the door
	return count >= 1, nil
}

func (r *repo) InsertLink(ctx context.Context, link uuid.UUID, creator string, time time.Time) error {
	query := `INSERT INTO link_management VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, link, creator, time)
	if err != nil {
		r.lg.Errorf("unable to insert link dut to %v", err)
		return err
	}
	return nil
}

type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

// NewMultiple return instance of the implemented interface
func NewMultiple(lg *logger.Logger, db *pgxpool.Pool) Multiple {
	return &repo{
		db: db,
		lg: lg,
	}
}

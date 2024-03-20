package link

import (
	"context"
	"errors"
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

func NewManager(db *pgxpool.Pool, lg *logger.Logger) Manager {
	return &link{
		db: db,
		lg: lg,
	}
}

type link struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

func (l *link) Remove(ctx context.Context, currentTime time.Time) error {
	query := `UPDATE link_management SET is_expired=true WHERE EXTRACT(EPOCH FROM ($1::timestamp - created_at)) > 3600`

	commandTag, err := l.db.Exec(ctx, query, currentTime)
	if err != nil {
		l.lg.Errorf("unable to update expiry due to %v", err)
		return err
	}
	if !commandTag.Update() {
		return errors.New("update not affected false")
	}
	return nil
}
func (l *link) Check(ctx context.Context, link uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM link_management WHERE link=$1 AND is_expired=$2`

	var count int

	err := l.db.QueryRow(ctx, query, link, false).Scan(&count)
	if err != nil {
		return false, err
	}
	// if true link exist we pass user further if no close the door
	return count >= 1, nil
}

func (l *link) Add(ctx context.Context, link uuid.UUID, creator string, time time.Time) error {
	query := `INSERT INTO link_management VALUES ($1, $2, $3)`
	_, err := l.db.Exec(ctx, query, link, creator, time)
	if err != nil {
		l.lg.Errorf("unable to insert link dut to %v", err)
		return err
	}
	return nil
}

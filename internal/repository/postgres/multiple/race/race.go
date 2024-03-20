package race

import (
	"context"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Multiple interface {
	Texts(ctx context.Context) ([]uuid.UUID, error)
	Text(ctx context.Context, id uuid.UUID) (string, error)
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

func (m *multiple) Texts(ctx context.Context) ([]uuid.UUID, error) {
	query := `SELECT id  FROM text`
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		m.lg.Errorf("unable to get text due to %v", err)
		return nil, err
	}
	defer rows.Close()

	var uuids []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		err = rows.Scan(&id)
		if err != nil {
			m.lg.Errorf("unable to scan row due to %v", err)
			return nil, err
		}

		uuids = append(uuids, id)
	}

	if err = rows.Err(); err != nil {
		m.lg.Errorf("error while reading rows due to %v", err)
		return nil, err
	}

	return uuids, nil
}

func (m *multiple) Text(ctx context.Context, id uuid.UUID) (string, error) {
	query := `SELECT content from text where id = $1`
	var content string
	err := m.db.QueryRow(ctx, query, id).Scan(&content)
	if err != nil {
		m.lg.Errorf("unable to get text due to %v", err)
		return "", err
	}
	return content, nil
}

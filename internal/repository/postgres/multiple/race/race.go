package race

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Multiple interface {
	Texts(ctx context.Context) ([]uuid.UUID, error)
	Text(ctx context.Context, id uuid.UUID) (string, error)
	AddRacers(ctx context.Context, mult models.MultipleRace) error
	User(ctx context.Context, id uuid.UUID) (models.RacerM, error)
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

func (m *multiple) User(ctx context.Context, id uuid.UUID) (models.RacerM, error) {
	query := `SELECT email, username, avatar, role FROM racer WHERE id = $1`

	var racer models.RacerM

	err := m.db.QueryRow(ctx, query, id).Scan(&racer.Email, &racer.Username, &racer.Avatar, &racer.Role)
	m.lg.Infof("racer %v", racer)
	return racer, err
}

func (m *multiple) Texts(ctx context.Context) ([]uuid.UUID, error) {
	query := `SELECT text_id  FROM random_text`
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

func (m *multiple) AddRacers(ctx context.Context, mult models.MultipleRace) error {
	query := `INSERT INTO multiple (generated_link, creator_id, track_name, created_at, racers, text_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := m.db.Exec(ctx, query, mult.GeneratedLink, mult.CreatorId, mult.TrackName, mult.CreatedAt, mult.Racers, mult.Text)
	if err != nil {
		m.lg.Errorf("unable to insert multiple")
		return err
	}
	return nil
}

package race

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Multiple interface {
	Texts(ctx context.Context) ([]uuid.UUID, error)
	Text(ctx context.Context, id uuid.UUID) (string, error)
	AddRacers(ctx context.Context, mult models.MultipleRace) error
	User(ctx context.Context, id uuid.UUID) (models.RacerM, error)
	RacerID(ctx context.Context, email string) (uuid.UUID, error)
	InsertSession(ctx context.Context, r *models.RacerRepoM) error
	UpdateRacerHistory(ctx context.Context, currSpeed int, racerID, textID uuid.UUID) error
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

func (m *multiple) RacerID(ctx context.Context, email string) (uuid.UUID, error) {
	query := `SELECT id FROM racer WHERE email = $1`
	var id uuid.UUID
	err := m.db.QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			m.lg.Errorf("no racer with such email %v", email)
			return uuid.Nil, nil
		}
		m.lg.Errorf("unable to get racer id due to %v", err)
		return uuid.Nil, err
	}
	return id, nil
}

func (m *multiple) InsertSession(ctx context.Context, r *models.RacerRepoM) error {
	query := `INSERT INTO multiple_session VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	queryHistory := `INSERT INTO multiple_history values ($1, $2)`

	fmt.Println(r, "MODEL RACERRESULT")

	b, err := m.db.Begin(ctx)
	if err != nil {
		m.lg.Errorf("unable to start transaction due to %v", err)
		return err
	}

	ex, err := b.Exec(ctx, query, r.GeneratedLink, r.RacerId, r.Duration, r.Wpm, r.Accuracy, r.StartTime, r.Winner, r.Place, r.TrackSize)
	if err != nil {
		m.lg.Errorf("unable to insert session due to %v", err)
		return err
	}
	if !ex.Insert() {
		return errors.New("unable insert user try again later")
	}

	_, err = b.Exec(ctx, queryHistory, r.GeneratedLink, r.RacerId)
	if err != nil {
		m.lg.Errorf("unable to insert history due to %v", err)
		return err
	}

	err = b.Commit(ctx)
	if err != nil {
		err2 := b.Rollback(ctx)
		if err2 != nil {
			m.lg.Errorf("unable to rollback transaction due to %v", err2)
			return err2
		}
		m.lg.Errorf("unable to commit transaction due to %v", err)
		return err
	}

	m.lg.Infof("successfully inserted session")
	return nil
}

func (m *multiple) UpdateRacerHistory(ctx context.Context, currSpeed int, racerID, textID uuid.UUID) error {
	updateRacer := `UPDATE racer
			  SET 
    				total_speed = total_speed + $1, 
    				last_race_speed = $1,
   					races = races + 1, 
   					best_speed = GREATEST(best_speed, $1),
    				avg_speed = (total_speed + $1) / (races + 1)
			  WHERE id = $2;`

	updateText := `UPDATE text 
			SET
				total_speed = total_speed + $1,
				occurrence = occurrence + 1,
				avg_speed = (total_speed + $1) / (occurrence + 1)
				where id = $2;`
	b, err := m.db.Begin(ctx)
	if err != nil {
		m.lg.Errorf("unable to start transaction due to %v", err)
		return err
	}

	_, err = b.Exec(ctx, updateRacer, currSpeed, racerID)
	if err != nil {
		m.lg.Errorf("unable to update racer due to %v", err)
		return err
	}

	_, err = b.Exec(ctx, updateText, currSpeed, textID)
	if err != nil {
		m.lg.Errorf("unable to update text due to %v", err)
		return err
	}

	err = b.Commit(ctx)
	if err != nil {
		if err := b.Rollback(ctx); err != nil {
			m.lg.Errorf("unable to rollback transaction due to %v", err)
			return err
		}
		m.lg.Errorf("unable to commit transaction due to %v", err)
		return err
	}

	m.lg.Infof("successfully updated racer and text")
	return nil
}

package authr

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (a *auth) DeleteRefreshSession(ctx context.Context, refresh string) error {
	query := `DELETE FROM session WHERE refresh_token = $1;`

	_, err := a.db.Exec(ctx, query, refresh)
	if err != nil {
		a.lg.Errorf("can't delete from session %v", err)
		return err
	}
	return nil
}

func (a *auth) InsertUser(ctx context.Context, racerAuth models.RacerAuth) error {
	begin, err := a.db.Begin(ctx)
	if err != nil {
		a.lg.Errorf("can't start a transaction %v", err)
		return err
	}
	query := `INSERT INTO racer(id, email, password, username, created_at, last_login, refresh_token, role)
      					values ($1, $2, $3, $4, $5, $6, $7, $8)`
	r := racerAuth
	exec, err := begin.Exec(ctx, query, r.ID, r.Email, r.Password, r.Username, r.CreatedAt, r.LastLogin, r.RefreshToken, r.Role)
	if err != nil {
		a.lg.Errorf("error with insert %v", err)
		return err
	}
	if !exec.Insert() {
		a.lg.Error("can't insert to the racer database")
		return errors.New("no insert")
	}

	querySession := `INSERT INTO session(user_id, last_login, role, refresh_token, fingerprint) values($1, $2, $3, $4, $5)`
	_, err = begin.Exec(ctx, querySession, r.ID, r.LastLogin, r.Role, r.RefreshToken, r.Fingerprint)
	if err != nil {
		a.lg.Errorf("error with insert %v", err)
		return err
	}

	if err != nil {
		a.lg.Errorf("can't in the sessionInsert %v", err)
		return err
	}

	err = begin.Commit(ctx)
	if err != nil {
		a.lg.Errorf("error with commit transaction %v", err)
		err = begin.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (a *auth) UserByEmail(ctx context.Context, email string) (bool, error) {
	a.lg.Info("in user by email method repo layer")
	query := `SELECT COUNT(*) FROM racer where email=$1`
	var racers int
	err := a.db.QueryRow(ctx, query, email).Scan(&racers)
	if err != nil {
		a.lg.Errorf("can't get racer data %v", err)
		return false, err
	}
	return racers == 1, nil
}

func (a *auth) GetUserPasswordByEmail(ctx context.Context, email string) (uuid.UUID, string, string, error) {
	var passwd string
	var token string
	var id uuid.UUID
	a.lg.Info("Dive into GetUser method")
	query := `SELECT id, refresh_token, password FROM racer where email=$1`
	err := a.db.QueryRow(ctx, query, email).Scan(&id, &token, &passwd)
	if err != nil {
		a.lg.Errorf("row not found %v", err)
		return [16]byte{}, "", "", err
	}
	return id, token, passwd, nil
}

func (a *auth) InsertSession(ctx context.Context, r models.RacerAuth) error {
	querySession := `INSERT INTO session(user_id, last_login, role, refresh_token, fingerprint) values($1, $2, $3, $4, $5)`
	_, err := a.db.Exec(ctx, querySession, r.ID, r.LastLogin, r.Role, r.RefreshToken, r.Fingerprint)
	if err != nil {
		a.lg.Errorf("error with insert %v", err)
		return err
	}
	return nil
}

// Fingerprint find fingerprint of the browser by refresh token and fingerprint
func (a *auth) Fingerprint(ctx context.Context, fng, refresh string) (models.RacerAuth, error) {
	var r models.RacerAuth
	query := `SELECT user_id, role, fingerprint FROM session where fingerprint=$1 and refresh_token=$2`
	// inject fingerprint here
	err := pgxscan.Get(ctx, a.db, &r, query, fng, refresh)
	a.lg.Infof("error is %v", err)
	return r, err
}

func (a *auth) DeleteSession(ctx context.Context, fng, refresh string) error {
	query := `DELETE FROM session WHERE refresh_token = $1 AND fingerprint = $2;`

	_, err := a.db.Exec(ctx, query, refresh, fng)
	if err != nil {
		a.lg.Errorf("can't delete from session %v", err)
		return err
	}
	return nil
}
func (a *auth) UserSession(ctx context.Context, token string, id uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM session where refresh_token=$1 and user_id=$2`
	var counter int
	err := a.db.QueryRow(ctx, query, token, id).Scan(&counter)
	if err != nil {
		a.lg.Errorf("can't fetch from session %v", err)
		return false, err
	}
	return counter == 1, nil
}

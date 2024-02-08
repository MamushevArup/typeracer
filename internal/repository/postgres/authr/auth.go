package authr

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth interface {
	GetUserPasswordByEmail(ctx context.Context, email string) (string, error)
	UserByEmail(ctx context.Context, email string) (bool, error)
	InsertUser(ctx context.Context, racerAuth models.RacerAuth) error
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
	err = begin.Commit(ctx)
	if err != nil {
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

func (a *auth) GetUserPasswordByEmail(ctx context.Context, email string) (string, error) {
	var passwd string
	a.lg.Info("Dive into GetUser method")
	query := `SELECT password FROM racer where email=$1`
	err := a.db.QueryRow(ctx, query, email).Scan(&passwd)
	if err != nil {
		a.lg.Errorf("row not found %v", err)
		return "", err
	}
	return passwd, nil
}

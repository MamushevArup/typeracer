package authr

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

func (a *auth) InsertSession(ctx context.Context, r models.RacerAuth) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Insert("session").
		Columns("user_id", "last_login", "role", "refresh_token", "fingerprint").
		Values(r.ID, r.LastLogin, r.Role, r.RefreshToken, r.Fingerprint).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query user=%v, err=%v", r.ID, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = a.db.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("fail exec query user=%v, err=%v", r.ID, err)
		return fmt.Errorf("fail insert err=%w", err)
	}

	return nil
}

// Fingerprint find fingerprint of the browser by refresh token and fingerprint
func (a *auth) Fingerprint(ctx context.Context, fng, refresh string) (models.RacerAuth, error) {

	var r models.RacerAuth

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("user_id, role, fingerprint").
		From("session").
		Where(squirrel.Eq{"fingerprint": fng, "refresh_token": refresh}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query err=%v", err)
		return r, fmt.Errorf("fail build query err=%w", err)
	}

	err = pgxscan.Get(ctx, a.db, &r, sql, args...)
	if err != nil {
		a.lg.Errorf("can't scan storage err=%v", err)
		return r, fmt.Errorf("fail scan storage err=%w", err)
	}

	a.lg.Info("successfully proceed query")

	return r, nil
}

func (a *auth) UserSession(ctx context.Context, token string, id uuid.UUID) (bool, error) {

	var counter int

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("COUNT(*)").
		From("session").
		Where(squirrel.Eq{"refresh_token": token, "user_id": id}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query id=%v, err=%v", id, err)
		return false, fmt.Errorf("fail build query err=%w", err)
	}

	err = a.db.QueryRow(ctx, sql, args...).Scan(&counter)
	if err != nil {
		a.lg.Errorf("error scanning storage user=%v, err=%v", id, err)
		return false, fmt.Errorf("fail fetch from session err=%w", err)
	}

	return counter == 1, nil
}

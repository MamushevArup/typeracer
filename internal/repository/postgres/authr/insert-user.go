package authr

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (a *auth) InsertUser(ctx context.Context, racerAuth models.RacerAuth) error {

	begin, err := a.db.Begin(ctx)
	if err != nil {
		a.lg.Errorf("fail start trasaction user=%v, err=%v", racerAuth.ID, err)
		return fmt.Errorf("fail start transaction err=%w", err)
	}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	r := racerAuth

	sql, args, err := sq.
		Insert("racer").
		Columns("id", "email", "password", "username", "created_at", "last_login", "refresh_token", "role").
		Values(r.ID, r.Email, r.Password, r.Username, r.CreatedAt, r.LastLogin, r.RefreshToken, r.Role).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query user=%v, err=%v", r.ID, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	exec, err := begin.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("error with insert %v", err)
		return err
	}

	if !exec.Insert() {
		a.lg.Errorf("can't insert to the racer database user=%v", r.ID)
		return fmt.Errorf("fail insert to the racer database")
	}

	sql, args, err = sq.
		Insert("session").
		Columns("user_id", "last_login", "role", "refresh_token", "fingerprint").
		Values(r.ID, r.LastLogin, r.Role, r.RefreshToken, r.Fingerprint).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query user=%v, err=%v", r.ID, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = begin.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("error with insert user=%v, err=%v", r.ID, err)
		return fmt.Errorf("fail insert err=%w", err)
	}

	err = begin.Commit(ctx)
	if err != nil {

		a.lg.Errorf("error with commit transaction user=%v, err=%v", r.ID, err)

		if err := begin.Rollback(ctx); err != nil {
			a.lg.Errorf("fail rollback user=%v, err=%v", r.ID, err)
			return fmt.Errorf("fail rollback err=%w", err)
		}

		return fmt.Errorf("fail commit err=%w", err)
	}

	return nil
}

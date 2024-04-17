package authr

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (a *auth) DeleteRefreshSession(ctx context.Context, refresh string) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Delete("session").
		Where(squirrel.Eq{"refresh_token": refresh}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query err=%v", err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = a.db.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("can't delete from session %v", err)
		return fmt.Errorf("fail delete from session err=%w", err)
	}

	return nil
}

func (a *auth) DeleteSession(ctx context.Context, fng, refresh string) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Delete("session").
		Where(squirrel.Eq{"refresh_token": refresh, "fingerprint": fng}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query fng=%v,  err=%v", fng, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = a.db.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("can't delete from session %v", err)
		return fmt.Errorf("fail delete from session err=%w", err)
	}

	return nil
}

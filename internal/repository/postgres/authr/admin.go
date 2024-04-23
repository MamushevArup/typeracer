package authr

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"golang.org/x/net/context"
)

func (a *auth) UpdateAdmin(ctx context.Context, username, refresh string) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Update("admin").
		Set("refresh_token", refresh).
		Set("username", username).
		Where(squirrel.Eq{"id": 1}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query for admin, err=%v", err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = a.db.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("error with insert to the admin database err=%v", err)
		return fmt.Errorf("fail insert to the admin database err=%w", err)
	}

	return nil
}

func (a *auth) AdminByRefresh(ctx context.Context, refresh string) (string, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("id").
		From("admin").
		Where(squirrel.Eq{"refresh_token": refresh}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query for admin, err=%v", err)
		return "", fmt.Errorf("fail build query err=%w", err)
	}

	var id string

	err = a.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		a.lg.Errorf("can't scan admin, err=%v", err)
		return "", fmt.Errorf("fail to scan admin %w", err)
	}

	return id, nil
}

func (a *auth) UpdateRefresh(ctx context.Context, id, refreshNew string) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Update("admin").
		Set("refresh_token", refreshNew).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query for admin, err=%v", err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	updated, err := a.db.Exec(ctx, sql, args...)
	if err != nil {
		a.lg.Errorf("can't update admin, err=%v", err)
		return fmt.Errorf("fail to update admin %w", err)
	}

	if !updated.Update() {
		a.lg.Errorf("can't update admin, id=%v, err=%v", id, err)
		return fmt.Errorf("fail to update admin %w", err)
	}

	return nil
}

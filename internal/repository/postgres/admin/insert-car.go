package admin

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (r *repo) InsertAvatarURL(ctx context.Context, url string) (bool, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Insert("avatar").
		Columns("url").
		Values(url).
		ToSql()

	if err != nil {
		r.lg.Errorf("failed to build sql query: %v", err)
		return false, fmt.Errorf("failed to build sql query: %w", err)
	}

	cond, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("failed to execute query: %v", err)
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	if cond.RowsAffected() == 0 {
		r.lg.Errorf("failed to insert avatar url")
		return false, fmt.Errorf("failed to insert avatar url")
	}

	return true, nil

}

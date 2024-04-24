package race

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (m *multiple) GuestAvatar(ctx context.Context) (string, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("url").
		From("avatar").
		OrderBy("RANDOM()").
		Limit(1).
		ToSql()

	if err != nil {
		m.lg.Error("unable to generate random avatar err=%v", err)
		return "", fmt.Errorf("unable to generate random avatar: %w", err)
	}

	var avatar string
	err = m.db.QueryRow(ctx, sql, args...).Scan(&avatar)
	if err != nil {
		m.lg.Error("unable to get random avatar err=%v", err)
		return "", fmt.Errorf("unable to get random avatar: %w", err)
	}

	return avatar, nil
}

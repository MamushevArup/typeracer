package authr

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

func (a *auth) UserAvatar(ctx context.Context, email string) (string, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("url").
		From("avatar").
		Join("racer on racer.avatar_id=avatar.id").
		Where(squirrel.Eq{"racer.email": email}).ToSql()

	if err != nil {
		a.lg.Errorf("fail constructing query user=%v, err=%v", email, err)
		return "", fmt.Errorf("fail to construct query err=%w", err)
	}

	var url string

	err = a.db.QueryRow(ctx, sql, args...).Scan(&url)
	if err != nil {
		a.lg.Errorf("fail fetching from storage user=%v, err=%v", email, err)
		return "", err
	}

	return url, nil
}

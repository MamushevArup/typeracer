package link

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (l *link) Check(ctx context.Context, link uuid.UUID) (bool, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("COUNT(*)").
		From("link_management").
		Where(squirrel.Eq{"link": link, "is_expired": false}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("fail to construct query link=%v, err=%w", link, err)
	}

	var count int

	err = l.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("fail to execute query link=%v, err=%w", link, err)
	}
	// if true link exist we pass user further if no close the door
	return count >= 1, nil
}

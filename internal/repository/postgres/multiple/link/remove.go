package link

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"
)

func (l *link) Remove(ctx context.Context, currentTime time.Time) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Update("link_management").
		Set("is_expired", true).
		Where(squirrel.Expr("EXTRACT(EPOCH FROM (?::timestamp - created_at)) > 3600", currentTime)).
		ToSql()

	if err != nil {
		return fmt.Errorf("fail to construct query time=%v, err=%w", currentTime, err)
	}

	commandTag, err := l.db.Exec(ctx, sql, args...)
	if err != nil {
		l.lg.Errorf("unable to update expiry due to %v", err)
		return fmt.Errorf("fail to update expiry %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("fail to update expiry")
	}

	return nil
}

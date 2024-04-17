package link

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

func (l *link) Add(ctx context.Context, link uuid.UUID, creator string, time time.Time) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Insert("link_management").
		Columns("link", "creator_id", "created_at").
		Values(link, creator, time).
		ToSql()

	if err != nil {
		return fmt.Errorf("fail to construct query link=%v, err=%w", link, err)
	}

	_, err = l.db.Exec(ctx, sql, args...)
	if err != nil {
		l.lg.Errorf("unable to insert link dut to %v", err)
		return fmt.Errorf("fail to insert link %w", err)
	}

	return nil
}

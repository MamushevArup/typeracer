package admin

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) DeleteModerationText(ctx context.Context, modId uuid.UUID) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Delete("moderation").
		Where(squirrel.Eq{"moderation_id": modId}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for delete moderation text, mod_id=%v err=%v", modId, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't delete moderation text, mod_id=%v err=%v", modId, err)
		return fmt.Errorf("fail delete moderation text err=%w", err)
	}

	return nil
}

package contributor

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) DuplicateContent(ctx context.Context, racerID uuid.UUID, contentHash string) (bool, error) {

	// if there is duplicate content, return count > 0
	var count int

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.Select("COUNT(*)").
		From("moderation").
		Where(squirrel.Eq{"content_hash": contentHash, "status": 1}).
		ToSql()

	if err != nil {
		r.lg.Errorf("fail build query user=%v, err=%v", racerID, err)
		return false, fmt.Errorf("sql query error: %w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		r.lg.Errorf("fail select user=%v, err=%v", racerID, err)
		return false, fmt.Errorf("fail select from moderation: %w", err)
	}

	// return true if there is duplicate content
	return count >= 1, nil
}

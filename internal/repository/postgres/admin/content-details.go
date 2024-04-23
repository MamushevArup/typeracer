package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) ContentDetails(ctx context.Context, modId uuid.UUID) (models.ModerationTextDetails, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	details := models.ModerationTextDetails{}

	sql, args, err := sq.
		Select("content", "author", "source", "source_title").
		From("moderation").
		Where(squirrel.Eq{"moderation_id": modId}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for select content details, mod_id=%v err=%v", modId, err)
		return details, fmt.Errorf("fail build query err=%w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&details.Content, &details.Author, &details.Source, &details.SourceTitle)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.lg.Errorf("content not found, mod_id=%v", modId)
			return details, fmt.Errorf("content not found")
		}
		r.lg.Errorf("can't select content details, mod_id=%v err=%v", modId, err)
		return details, fmt.Errorf("fail select content details err=%w", err)
	}

	return details, nil
}

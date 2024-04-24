package contributor

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (r *repo) InsertToModeration(ctx context.Context, c models.ContributeServiceRequest) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Insert("moderation").
		Columns("moderation_id", "racer_id", "content", "author",
			"length", "source", "source_title", "sent_at", "status", "content_hash").
		Values(c.ModerationID, c.RacerID, c.Content, c.Author,
			c.Length, c.Source, c.SourceTitle, c.SentAt, c.Status, c.ContentHash).
		ToSql()

	if err != nil {
		r.lg.Errorf("fail build query user=%v, err=%v", c.RacerID, err)
		return fmt.Errorf("sql query error: %w", err)
	}

	i, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("fail insert user=%v, err=%v", c.RacerID, err)
		return fmt.Errorf("fail insert to moderation: %w", err)
	}

	if i.RowsAffected() == 0 {
		r.lg.Errorf("no rows added user=%v, err=%v", c.RacerID, err)
		return fmt.Errorf("no rows added: %w", err)
	}

	return nil
}

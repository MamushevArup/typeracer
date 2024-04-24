package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (r *repo) SelectModeration(ctx context.Context, limit, offset int, sort string) ([]models.ModerationRepoResponse, error) {

	var resp []models.ModerationRepoResponse

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("moderation_id", "sent_at", "username").
		From("moderation").
		Join("racer on moderation.racer_id=racer.id").
		Where(squirrel.Eq{"moderation.status": 0}).
		OrderBy("sent_at " + sort).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for select content, err=%v", err)
		return resp, fmt.Errorf("fail build query err=%w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't select content, err=%v", err)
		return resp, fmt.Errorf("fail select content err=%w", err)
	}

	for rows.Next() {
		var currResponse models.ModerationRepoResponse
		err = rows.Scan(&currResponse.ModerationID, &currResponse.SentAt, &currResponse.Username)
		if err != nil {
			r.lg.Errorf("can't scan content, err=%v", err)
			return resp, fmt.Errorf("fail scan content err=%w", err)
		}
		resp = append(resp, currResponse)
	}

	defer rows.Close()

	if len(resp) == 0 {
		return resp, fmt.Errorf("no content to moderate")
	}

	return resp, nil
}

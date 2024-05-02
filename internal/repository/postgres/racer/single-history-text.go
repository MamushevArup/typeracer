package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

func (r *repo) SelectHistoryText(ctx context.Context, racerUUID, singleID uuid.UUID) (models.SingleHistoryText, error) {
	var hst models.SingleHistoryText

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("content", "author", "source", "source_title", "username as contributor").
		From("text").
		Join("single on single.text_id = text.id").
		Join("racer on racer.id=text.contributor_id").
		Where(squirrel.Eq{"racer_id": racerUUID}).
		Where(squirrel.Eq{"single.id": singleID}).
		ToSql()

	if err != nil {
		r.lg.Errorf("fail to construct query raceId=%v, racerId=%v, err=%v", singleID, racerUUID, err)
		return hst, ErrQueryBuild
	}

	err = pgxscan.Get(ctx, r.db, &hst, sql, args...)
	if err != nil {
		r.lg.Errorf("fail to get history text, raceId=%v, err=%v", singleID, err)
		return hst, fmt.Errorf("fail to get history text %w", err)
	}

	r.lg.Info("successfully get history text")
	return hst, nil
}

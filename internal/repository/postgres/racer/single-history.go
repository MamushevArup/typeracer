package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) SingleHistoryRows(ctx context.Context, racerId uuid.UUID, limit, offset int) ([]models.SingleHistory, error) {
	var hst []models.SingleHistory
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("id", "speed", "duration", "accuracy", "started_time").
		From("single").
		Where(squirrel.Eq{"racer_id": racerId}).
		OrderBy("started_time DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for select single history userId=%v, err=%v", racerId, err)
		return hst, ErrQueryBuild
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't select single history userId=%v, err=%v", racerId, err)
		return hst, fmt.Errorf("fail select single history err=%w", err)
	}

	for rows.Next() {
		var curr models.SingleHistory
		err = rows.Scan(&curr.SingleID, &curr.Speed, &curr.Duration, &curr.Accuracy, &curr.StartedAt)
		if err != nil {
			r.lg.Errorf("can't scan single history userId=%v, err=%v", racerId, err)
			return hst, fmt.Errorf("fail scan single history err=%w", err)
		}
		hst = append(hst, curr)
	}
	defer rows.Close()

	r.lg.Info("single history rows selected")
	return hst, nil
}

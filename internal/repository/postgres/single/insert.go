package single

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (r *repo) EndSingleRace(ctx context.Context, resp models.RespEndSingle) error {

	begin, err := r.db.Begin(ctx)
	if err != nil {
		r.lg.Errorf("unable to start transaction %v", err)
		return fmt.Errorf("unable to start transaction %w user_id=%v", err, resp.RacerId)
	}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Insert("single").
		Columns("id", "speed", "duration", "accuracy", "start_time", "racer_id", "text_id").
		Values(resp.RaceId, resp.Wpm, resp.Duration, resp.Accuracy, resp.StartedTime, resp.RacerId, resp.TextId).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query %w user_id=%v", err, resp.RacerId)
	}

	_, err = begin.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("error with inserting into single transaction failed %v, user_id=%v, race_id=%v", err, resp.RacerId, resp.RaceId)
		return fmt.Errorf("error with inserting into single transaction failed %w user_id=%v", err, resp.RacerId)
	}

	sql, args, err = sq.
		Insert("race_history").
		Columns("single_id", "racer_id", "text_id").
		Values(resp.RaceId, resp.RacerId, resp.TextId).ToSql()

	_, err = begin.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("fail insert race_history user_id=%v, err: %v", resp.RacerId, err)
		return fmt.Errorf("fail insert race_history %w user_id=%v", err, resp.RacerId)
	}

	if err = begin.Commit(ctx); err != nil {

		if err = begin.Rollback(ctx); err != nil {
			r.lg.Errorf("unable to rollback %v, user_id=%v", err, resp.RacerId)
			return fmt.Errorf("unable to rollback %w user_id=%v", err, resp.RacerId)
		}

		r.lg.Errorf("unable to commit %v, user_id=%v", err, resp.RacerId)
		return fmt.Errorf("unable to commit %w user_id=%v", err, resp.RacerId)
	}

	return nil
}

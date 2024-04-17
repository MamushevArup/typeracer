package race

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (m *multiple) AddRacers(ctx context.Context, mlt models.MultipleRace) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Insert("multiple").
		Columns("generated_link", "creator_id", "track_name", "created_at", "racers", "text_id").
		Values(mlt.GeneratedLink, mlt.CreatorId, mlt.TrackName, mlt.CreatedAt, mlt.Racers, mlt.Text).
		ToSql()

	if err != nil {
		return fmt.Errorf("fail to construct query user=%v, link=%v, err=%w", mlt.CreatorId, mlt.GeneratedLink, err)
	}

	_, err = m.db.Exec(ctx, sql, args...)
	if err != nil {
		m.lg.Errorf("unable to insert multiple user=%v, link=%v due to %v", mlt.CreatorId, mlt.GeneratedLink, err)
		return fmt.Errorf("fail to insert multiple %w", err)
	}

	return nil
}

func (m *multiple) InsertSession(ctx context.Context, r *models.RacerRepoM) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := sq.
		Insert("multiple_session").
		Columns("generated_link", "racer_id", "duration", "wpm", "accuracy", "start_time", "winner", "place", "track_size").
		Values(r.GeneratedLink, r.RacerId, r.Duration, r.Wpm, r.Accuracy, r.StartTime, r.Winner, r.Place, r.TrackSize).
		ToSql()

	if err != nil {
		m.lg.Errorf("unable to construct query user=%v, link=%v, err=%v", r.RacerId, r.GeneratedLink, err)
		return fmt.Errorf("unable to construct query user=%v, link=%v, err=%w", r.RacerId, r.GeneratedLink, err)
	}

	b, err := m.db.Begin(ctx)
	if err != nil {
		m.lg.Errorf("unable to start transaction due to %v", err)
		return fmt.Errorf("fail to start transaction due to %w", err)
	}

	ex, err := b.Exec(ctx, sql, args...)
	if err != nil {
		m.lg.Errorf("unable to insert session user=%v, due to %v", r.RacerId, err)
		return fmt.Errorf("fail to insert session due to %w", err)
	}

	if !ex.Insert() {
		return fmt.Errorf("fail to insert user=%v, link=%v", r.RacerId, r.GeneratedLink)
	}

	sql, args, err = sq.
		Insert("multiple_history").
		Columns("generated_link", "racer_id").
		Values(r.GeneratedLink, r.RacerId).
		ToSql()

	if err != nil {
		return fmt.Errorf("unable to construct query user=%v, link=%v, err=%w", r.RacerId, r.GeneratedLink, err)
	}

	_, err = b.Exec(ctx, sql, args...)
	if err != nil {
		m.lg.Errorf("unable to insert history user=%v, link=%v  due to %v", r.RacerId, r.GeneratedLink, err)
		return fmt.Errorf("fail to insert history due to %w", err)
	}

	err = b.Commit(ctx)
	if err != nil {

		err2 := b.Rollback(ctx)
		if err2 != nil {
			m.lg.Errorf("unable to rollback transaction due to %v", err2)
			return fmt.Errorf("fail to rollback transaction due to %w", err2)
		}

		m.lg.Errorf("unable to commit transaction due to %v", err)
		return fmt.Errorf("fail to commit transaction due to %w", err)
	}

	m.lg.Infof("successfully inserted session")
	return nil
}

package race

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (m *multiple) UpdateRacerHistory(ctx context.Context, currSpeed int, racerID, textID uuid.UUID) error {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sql, args, err := sq.Update("racer").
		Set("total_speed", squirrel.Expr("total_speed + ?", currSpeed)).
		Set("last_race_speed", currSpeed).
		Set("races", squirrel.Expr("races + 1")).
		Set("best_speed", squirrel.Expr("GREATEST(best_speed, ?)", currSpeed)).
		Set("avg_speed", squirrel.Expr("(total_speed + ?) / (races + 1)", currSpeed)).
		Where("id = ?", racerID).
		ToSql()

	if err != nil {
		return fmt.Errorf("fail to construct query user=%v, err=%w", racerID, err)
	}

	b, err := m.db.Begin(ctx)
	if err != nil {
		m.lg.Errorf("unable to start transaction user=%v due to %v", racerID, err)
		return fmt.Errorf("fail to start transaction due to %w", err)
	}

	_, err = b.Exec(ctx, sql, args...)
	if err != nil {
		m.lg.Errorf("unable to update racer user=%v due to %v", racerID, err)
		return fmt.Errorf("fail to update racer due to %w", err)
	}

	sql, args, err = sq.Update("text").
		Set("total_speed", squirrel.Expr("total_speed + ?", currSpeed)).
		Set("occurrence", squirrel.Expr("occurrence + 1")).
		Set("avg_speed", squirrel.Expr("(total_speed + ?) / (occurrence + 1)", currSpeed)).
		Where("id = ?", textID).
		ToSql()

	if err != nil {
		return fmt.Errorf("fail to construct query user=%v, err=%w", racerID, err)
	}

	_, err = b.Exec(ctx, sql, args...)
	if err != nil {
		m.lg.Errorf("unable to update text text=%v, user=%v due to %v", textID, racerID, err)
		return fmt.Errorf("fail to update text due to %w", err)
	}

	err = b.Commit(ctx)
	if err != nil {

		if err2 := b.Rollback(ctx); err2 != nil {
			m.lg.Errorf("unable to rollback transaction user=%v, text=%v due to %v", racerID, textID, err2)
			return fmt.Errorf("fail to rollback transaction due to %w", err2)
		}

		m.lg.Errorf("unable to commit transaction user=%v, text=%v due to %v", racerID, textID, err)
		return fmt.Errorf("fail to commit transaction due to %w", err)
	}

	m.lg.Infof("successfully updated racer and text")
	return nil
}

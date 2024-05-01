package racer

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

var ErrQueryBuild = errors.New("can't build query")

func (r *repo) Info(ctx context.Context, racerId uuid.UUID) (models.RacerRepository, error) {
	var racer models.RacerRepository

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("username", "url as avatar", "created_at", "avg_speed", "last_race_speed", "best_speed", "races").
		From("racer").
		Join("avatar on racer.avatar_id = avatar.id").
		Where(squirrel.Eq{"racer.id": racerId}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query id=%v, err=%v", racerId, err)
		return racer, ErrQueryBuild
	}

	err = pgxscan.Get(ctx, r.db, &racer, sql, args...)
	if err != nil {
		r.lg.Errorf("fail to select racer id=%v, err=%v", racerId, err)
		return racer, fmt.Errorf("no rows found err=%w", err)
	}

	r.lg.Errorf("info method work fine id=%v", racerId)
	return racer, nil
}

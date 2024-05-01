package racer

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

var ErrUpdate = errors.New("unable update data")

func (r *repo) UpdateAvatar(ctx context.Context, avatar models.AvatarUpdateRepo) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Update("racer").
		Set("avatar_id", avatar.Id).
		Where(squirrel.Eq{"id": avatar.RacerID}).
		Where(squirrel.Expr("EXISTS (SELECT 1 FROM avatar WHERE id = ?)", avatar.Id)).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query %v", err)
		return ErrQueryBuild
	}

	flag, err := r.db.Exec(ctx, sql, args...)

	if flag.RowsAffected() == 0 {
		r.lg.Errorf("can't update avatar %v", err)
		return ErrUpdate
	}

	r.lg.Info("avatar updated successfully")
	return nil
}

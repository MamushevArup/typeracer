package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (r *repo) SelectAvatar(ctx context.Context) ([]models.Avatar, error) {
	avatars := make([]models.Avatar, 0)

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("id, url").
		From("avatar").
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query %v", err)
		return nil, ErrQueryBuild
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't get avatars %v", err)
		return nil, fmt.Errorf("can't get avatars %w", err)
	}

	for rows.Next() {
		var avatar models.Avatar
		err = rows.Scan(&avatar.Id, &avatar.Url)
		if err != nil {
			r.lg.Errorf("can't scan avatar %v", err)
			return nil, fmt.Errorf("can't scan avatar %w", err)
		}
		avatars = append(avatars, avatar)
	}
	defer rows.Close()
	r.lg.Info("avatars fetched successfully")
	return avatars, nil
}

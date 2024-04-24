package race

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (m *multiple) User(ctx context.Context, id uuid.UUID) (models.RacerM, error) {
	var racer models.RacerM

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, args, err := sq.
		Select("email, username, url as avatar, role").
		From("racer").
		Join("avatar a on racer.avatar_id = a.id").
		Where("racer.id = ?", id).
		ToSql()

	if err != nil {
		return racer, fmt.Errorf("fail to construct query user=%v, err=%w", id, err)
	}

	err = m.db.QueryRow(ctx, query, args...).Scan(&racer.Email, &racer.Username, &racer.Avatar, &racer.Role)
	if err != nil {
		return racer, fmt.Errorf("fail to get user info user=%v, err=%w", id, err)
	}

	m.lg.Infof("racer %v", racer)
	return racer, nil
}

func (m *multiple) RacerID(ctx context.Context, email string) (uuid.UUID, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("id").
		From("racer").
		Where("email = ?", email).
		ToSql()

	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to construct query email=%v, err=%w", email, err)
	}

	var id uuid.UUID

	err = m.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			m.lg.Errorf("no racer with such email %v", email)
			return uuid.Nil, nil
		}

		m.lg.Errorf("unable to get racer id due to %v", err)
		return uuid.Nil, fmt.Errorf("fail to get racer id %w", err)
	}

	return id, nil
}

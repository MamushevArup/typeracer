package authr

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (a *auth) UserByEmail(ctx context.Context, email string) (*models.SignInService, error) {
	a.lg.Info("in user by email method repo layer")

	var racerInfo models.SignInService

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("username, url").
		From("racer").
		Join("avatar on avatar.id=racer.avatar_id").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query user=%v, err=%v", email, err)
		return &racerInfo, fmt.Errorf("can't build query: %w", err)
	}

	err = a.db.QueryRow(ctx, sql, args...).Scan(&racerInfo.Username, &racerInfo.Avatar)
	if err != nil {
		a.lg.Errorf("can't get racer data email=%v, err=%v", email, err)
		return &racerInfo, fmt.Errorf("fail get racer data: %w", err)
	}

	return &racerInfo, nil
}

func (a *auth) GetUserPasswordByEmail(ctx context.Context, email string) (uuid.UUID, string, string, error) {

	var passwd string
	var token string
	var id uuid.UUID

	a.lg.Info("Dive into GetUser method")

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("id, refresh_token, password").
		From("racer").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		a.lg.Errorf("can't build query email=%v, err=%v", email, err)
		return uuid.Nil, "", "", fmt.Errorf("fail build query err=%w", err)
	}

	err = a.db.QueryRow(ctx, sql, args...).Scan(&id, &token, &passwd)
	if err != nil {
		a.lg.Errorf("row not found for email=%v  err=%v", email, err)
		return uuid.Nil, "", "", fmt.Errorf("row not found err=%w", err)
	}

	return id, token, passwd, nil
}

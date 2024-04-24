package single

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

func (r *repo) RacerInfo(ctx context.Context, id uuid.UUID) (models.RacerInfo, error) {
	var racer models.RacerInfo

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.
		Select("username, url as avatar").
		From("racer").
		Join("avatar a on racer.avatar_id = a.id").
		Where(squirrel.Eq{"racer.id": id}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query %v", err)
		return racer, fmt.Errorf("can't build query %w", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&racer.Username, &racer.Avatar)
	if err != nil {
		return models.RacerInfo{}, fmt.Errorf("no rows found %w", err)
	}

	return racer, nil
}

func (r *repo) TextInfo(ctx context.Context) (models.TextInfo, error) {

	var textInfo models.TextInfo

	textUUID := r.randomText(ctx)
	textInfo.TextID = textUUID

	// fetch data from text and contributor table and place it in single
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := sq.Select("content, source, source_title, author, r.username as contributor_name").
		From("racer r").
		Join("text t on r.id = t.contributor_id").
		Where("t.id = ?", textUUID)

	sql, args, err := query.ToSql()
	if err != nil {
		return textInfo, fmt.Errorf("fatal in query building query is : %v fix: %w", sql, err)
	}

	err = pgxscan.Get(ctx, r.db, &textInfo, sql, args...)
	if err != nil {
		r.lg.Errorf("can't scan text data %v", err)
		return textInfo, fmt.Errorf("query : %v, args : %v : %w", sql, args, err)
	}

	r.lg.Infof("text info %v", textInfo)
	return textInfo, nil
}

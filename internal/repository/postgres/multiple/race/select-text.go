package race

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (m *multiple) Texts(ctx context.Context) ([]uuid.UUID, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query, _, err := sq.
		Select("text_id").
		From("random_text").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("fail to construct query %w", err)
	}

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		m.lg.Errorf("fail to execute query due to %v", err)
		return nil, fmt.Errorf("fail to get texts %w", err)
	}

	defer rows.Close()

	var uuids []uuid.UUID

	for rows.Next() {
		var id uuid.UUID

		err = rows.Scan(&id)
		if err != nil {
			m.lg.Errorf("unable scan value:%v, err=%v", id, err)
			return nil, fmt.Errorf("fail to get textID:%v, err=%w", id, err)
		}

		uuids = append(uuids, id)
	}

	if err = rows.Err(); err != nil {
		m.lg.Errorf("error while reading rows due to %v", err)
		return nil, fmt.Errorf("fail occur in rows reading %w", err)
	}

	return uuids, nil
}

func (m *multiple) Text(ctx context.Context, textID uuid.UUID) (string, error) {

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.
		Select("content").
		From("text").
		Where("id = ?", textID).
		ToSql()

	if err != nil {
		m.lg.Errorf("fail to construct query text=%v, err=%v", textID, err)
		return "", fmt.Errorf("fail to construct query text=%v, err=%w", textID, err)
	}

	var content string

	err = m.db.QueryRow(ctx, query, args...).Scan(&content)
	if err != nil {
		m.lg.Errorf("unable to get text due to %v", err)
		return "", fmt.Errorf("fail to execute query %w", err)
	}

	return content, nil
}

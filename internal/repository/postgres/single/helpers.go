package single

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"math/big"
)

func (r *repo) randomText(ctx context.Context) uuid.UUID {
	var uuids []uuid.UUID

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, _, err := sq.
		Select("*").
		From("random_text").
		ToSql()

	if err != nil {
		r.lg.Errorf("error constructing query %v, error %v", sql, err)
		return uuid.Nil
	}

	txUUIDS, err := r.db.Query(ctx, sql)
	if err != nil {
		r.lg.Errorf("query failed due to %v", err)
		return uuid.Nil
	}

	for txUUIDS.Next() {
		var id uuid.UUID
		if err = txUUIDS.Scan(&id); err != nil {
			r.lg.Errorf("can't scan uuid %v value in uuid's set %v", id, err)
		}

		uuids = append(uuids, id)
	}

	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(uuids))))
	if err != nil {
		r.lg.Errorf("unable randomize array %v due to %v", uuids, err)
		return uuid.Nil

	}

	random := randomIndex.Int64()
	return uuids[random]
}

func (r *repo) GetTextLen(ctx context.Context, textID uuid.UUID) (int, error) {
	var length int

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.Select("length").From("text").Where("id = ?", textID).ToSql()
	if err != nil {
		r.lg.Errorf("textID: %v, err : %v", textID, err)
		return 0, fmt.Errorf("failed to build query text_id : %v : %w", textID, err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&length)
	if err != nil {
		r.lg.Errorf("can't fetch length from text %v", err)
		return 0, fmt.Errorf("no rows found with text_id : %v : %w", textID, err)
	}

	return length, nil
}

func (r *repo) RacerExist(ctx context.Context, id uuid.UUID) (bool, error) {

	var ex bool
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, _, err := sq.
		Select("EXISTS(SELECT 1 FROM racer WHERE id = ?)").
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query %v", err)
		return false, fmt.Errorf("failed to build query %w user_id=%v", err, id)
	}

	err = r.db.QueryRow(ctx, sql, id).Scan(&ex)
	if err != nil {
		r.lg.Errorf("can't execute use exist check %v", err)
		return false, fmt.Errorf("racer doesn't exist %w", err)
	}

	return ex, nil
}

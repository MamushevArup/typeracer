package single

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
)

type Starter interface {
	RacerInfo(ctx context.Context, id uuid.UUID) (models.RacerInfo, error)
	TextInfo(ctx context.Context) (models.TextInfo, error)

	EndSingleRace(ctx context.Context, req *models.RespEndSingle) error
	GetTextLen(ctx context.Context) (int, error)
	RacerExist(ctx context.Context, id uuid.UUID) (bool, error)
}

// This is responsible for the practice yourself section
type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

// New return instance of the implemented interface
func New(lg *logger.Logger, db *pgxpool.Pool) Starter {
	return &repo{
		db: db,
		lg: lg,
	}
}

type ids struct {
	raceId uuid.UUID
	textId uuid.UUID
}

var identifiers ids

func (r *repo) RacerInfo(ctx context.Context, id uuid.UUID) (models.RacerInfo, error) {
	var racer models.RacerInfo

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := sq.Select("username, avatar").From("racer").Where(squirrel.Eq{"id": id}).ToSql()
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

func (r *repo) randomText(ctx context.Context) uuid.UUID {
	var uuids []uuid.UUID

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, _, err := sq.Select("*").From("random_text").ToSql()
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

func (r *repo) EndSingleRace(ctx context.Context, resp *models.RespEndSingle) error {
	begin, err := r.db.Begin(ctx)
	if err != nil {
		r.lg.Errorf("unable to start transaction %v", err)
		return err
	}
	sgl := fmt.Sprintf("INSERT INTO single values ($1, $2, $3, $4, $5, $6, $7);")
	_, err = begin.Exec(ctx, sgl, identifiers.raceId, resp.Wpm, resp.Duration, resp.Accuracy, resp.StartedTime, resp.RacerId, identifiers.textId)
	if err != nil {
		r.lg.Errorf("error with inserting into single transaction failed %v", err)
		return err
	}
	raceH := fmt.Sprintf("INSERT INTO race_history VALUES ($1, $2, $3, $4);")
	_, err = begin.Exec(ctx, raceH, identifiers.raceId, resp.RacerId, identifiers.textId, 1)
	if err != nil {
		r.lg.Errorf("error with inserting into race_history %v", err)
		return err
	}

	if err = begin.Commit(ctx); err != nil {
		if err = begin.Rollback(ctx); err != nil {
			r.lg.Errorf("unable to rollback %v", err)
			return err
		}
		r.lg.Errorf("unable to commit %v", err)
		return err
	}
	return nil
}

func (r *repo) GetTextLen(ctx context.Context) (int, error) {
	var length int
	query := "SELECT length FROM text where id=$1"
	err := r.db.QueryRow(ctx, query, identifiers.textId).Scan(&length)
	if err != nil {
		r.lg.Errorf("can't fetch length from text %v", err)
		return 0, err
	}
	return length, nil
}

func (r *repo) RacerExist(ctx context.Context, id uuid.UUID) (bool, error) {
	var ex bool

	query := "SELECT EXISTS(SELECT 1 FROM racer WHERE id = $1)"

	err := r.db.QueryRow(ctx, query, id).Scan(&ex)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("racer doesn't exist: %w", err)
		}

		r.lg.Errorf("can't execute use exist check %v", err)
		return false, fmt.Errorf("%w", err)
	}

	return ex, nil
}

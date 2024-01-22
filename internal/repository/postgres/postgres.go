package postgres

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Starter interface {
	StartSingle(ctx context.Context, id uuid.UUID) (*models.Single, error)
}

// This is responsible for the practice yourself section
type repo struct {
	db *pgxpool.Pool
	lg *logger.Logger
}

// NewSingle return instance of the implemented interface
func NewSingle(lg *logger.Logger, db *pgxpool.Pool) Starter {
	return &repo{
		db: db,
		lg: lg,
	}
}

func (r *repo) StartSingle(ctx context.Context, userID uuid.UUID) (*models.Single, error) {
	// TODO implement me
	// Need to work with practice yourself section
	// Modify text insert single update user
	// prepare data to insert and then fetch the data from the database.
	//create new single race id
	var single models.Single
	newRaceID, err := uuid.NewUUID()
	if err != nil {
		r.lg.Errorf("can't create uuid for single race %v\n", err)
		return nil, err
	}
	single.ID = newRaceID
	single.RacerID = userID
	// using racer_id fetch the data related to the racer and
	// with random uuid fetch the text.
	begin, err := r.db.Begin(ctx)
	if err != nil {
		r.lg.Errorf("can't start a transaction %v", err)
		return nil, err
	}
	// fetch data from random_text table and return one random
	textUUID := r.randomText(ctx)

	// fetch data from racer table and place it in single
	racer := "SELECT username, avatar FROM racer WHERE id = $1"
	err = pgxscan.Get(ctx, r.db, &single, racer, userID)
	if err != nil {
		r.lg.Errorf("can't scan racer data %v", err)
		return nil, err
	}

	// fetch data from text and contributor table and place it in single
	text := fmt.Sprintf("SELECT content, length, contributor FROM text JOIN " +
		"contributor on text.contributor_id=contributor.user_id where text.id=$1 and contributor_id=$2")
	err = pgxscan.Get(ctx, r.db, &single, text, textUUID, userID)
	if err != nil {
		r.lg.Errorf("can't scan text data %v", err)
		return nil, err
	}

	fmt.Println(single)

	err = begin.Commit(ctx)
	if err != nil {
		err = begin.Rollback(ctx)
		if err != nil {
			r.lg.Errorf("can't rollback a transaction %v", err)
			return nil, err
		}
		r.lg.Errorf("can't commit a transaction %v", err)
		return nil, err
	}
	return &single, nil
}

func (r *repo) randomText(ctx context.Context) *uuid.UUID {
	var uuids []uuid.UUID
	getRandom := fmt.Sprintf("SELECT * FROM random_text")
	txUUIDS, err := r.db.Query(ctx, getRandom)
	if err != nil {
		r.lg.Errorf("can't exec query random uuids %v", err)
		return nil
	}
	for txUUIDS.Next() {
		var id uuid.UUID
		if err = txUUIDS.Scan(&id); err != nil {
			r.lg.Errorf("can't scan value in uuid's set %v", err)
		}
		uuids = append(uuids, id)
	}
	return &uuids[0]
}

package single

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
)

type Starter interface {
	StartSingle(ctx context.Context, id, racerID uuid.UUID) (*models.Single, error)
	EndSingleRace(ctx context.Context, req *models.RespEndSingle) error
	GetTextLen(ctx context.Context) (int, error)
	RacerExist(ctx context.Context, id uuid.UUID) (bool, error)
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

type ids struct {
	raceId uuid.UUID
	textId uuid.UUID
}

var identifiers ids

func (r *repo) StartSingle(ctx context.Context, userID, raceID uuid.UUID) (*models.Single, error) {
	// Need to work with practice yourself section
	// Modify text insert single update user
	// prepare data to insert and then fetch the data from the database.
	//create new single race id
	single := new(models.Single)

	single.ID = raceID
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
	single.TextID = textUUID
	identifiers.textId = textUUID
	identifiers.raceId = raceID
	// fetch data from racer table and place it in single
	racer := "SELECT username, avatar FROM racer WHERE id = $1"
	err = pgxscan.Get(ctx, begin, single, racer, userID)
	if err != nil {
		r.lg.Errorf("can't scan racer data %v", err)
		return nil, err
	}

	// fetch data from text and contributor table and place it in single
	text := fmt.Sprintf("SELECT content, length, contributor FROM text JOIN " +
		"contributor on text.contributor_id=contributor.user_id where text.id=$1 and contributor_id=$2")
	err = pgxscan.Get(ctx, begin, single, text, textUUID, userID)
	if err != nil {
		r.lg.Errorf("can't scan text data %v", err)
		return nil, err
	}

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
	return single, nil
}

func (r *repo) randomText(ctx context.Context) uuid.UUID {
	var uuids []uuid.UUID
	getRandom := fmt.Sprintf("SELECT * FROM random_text")
	txUUIDS, err := r.db.Query(ctx, getRandom)
	if err != nil {
		r.lg.Errorf("can't exec query random uuids %v", err)
		return [16]byte{}
	}
	for txUUIDS.Next() {
		var id uuid.UUID
		if err = txUUIDS.Scan(&id); err != nil {
			r.lg.Errorf("can't scan value in uuid's set %v", err)
		}
		uuids = append(uuids, id)
	}
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(uuids))))
	if err != nil {
		// Handle error (e.g., log, panic, return a default value)
		r.lg.Errorf("can't randomize due to %v", err)
		return [16]byte{}

	}

	// Return the UUID at the randomly selected index
	return uuids[randomIndex.Int64()]
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
		r.lg.Errorf("can't execute use exist check %v", err)
		return false, err
	}
	return ex, nil
}

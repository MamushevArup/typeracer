package single

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
)

const (
	guest = "guest"
)

type Practice interface {
	StartRace(ctx context.Context, userId string) (models.SingleResponse, error)
	EndRace(ctx context.Context, req models.ReqEndSingle) (models.RespEndSingle, error)
	RacerExists(ctx context.Context, racerId string) (bool, error)
	RealTimeCalc(ctx context.Context, currentSymbol, duration int) (int, error)
}

type service struct {
	repo *repository.Repo
	ids  identifiers
}

type identifiers struct {
	textUUID  uuid.UUID
	racerUUID uuid.UUID
	raceUUID  uuid.UUID
}

func NewPracticeY(repo *repository.Repo) Practice {
	return &service{
		repo: repo,
		ids:  identifiers{},
	}
}

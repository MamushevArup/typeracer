package race

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/gorilla/websocket"
	"sync"
)

type Racer interface {
	RandomText(ctx context.Context) (string, error)
	Join(token string, conn *websocket.Conn, link string) (*[]*websocket.Conn, models.RacerM, error)
	Timer(link string, cons *[]*websocket.Conn) (int, error)
	WhiteLine(ctx context.Context, link string) error
	CurrentSpeed(racer *models.RacerCurrentWpm, textLen int) (models.RacerSpeed, error)
	EndRace(raceReq models.RaceEndRequest, link, id string) (models.RaceResult, error)
}

type service struct {
	repo           *repository.Repo
	timers         map[string]int
	cfg            *config.Config
	racers         []string
	connections    []*websocket.Conn
	mu             sync.Mutex
	d              data
	finishedRacers map[string]int
}

func NewMultiple(repo *repository.Repo, cfg *config.Config) Racer {
	return &service{
		repo:           repo,
		cfg:            cfg,
		timers:         make(map[string]int),
		racers:         make([]string, 0, 5),
		connections:    make([]*websocket.Conn, 0, 5),
		d:              data{},
		finishedRacers: make(map[string]int),
	}
}

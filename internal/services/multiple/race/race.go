package race

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Racer interface {
	RandomText(ctx context.Context) (string, error)
	Join(token string, conn *websocket.Conn, link string) (*[]*websocket.Conn, error)
	Timer(link string, cons *[]*websocket.Conn) (int, error)
	WhiteLine(ctx context.Context, link string) error
}

type service struct {
	repo        *repository.Repo
	timers      map[string]int
	cfg         *config.Config
	racers      []string
	connections []*websocket.Conn
	mu          sync.Mutex
	d           data
}

// data struct hold values need to pass to the repo layer
type data struct {
	textID    uuid.UUID
	createdAt time.Time
}

func (s *service) Join(token string, conn *websocket.Conn, link string) (*[]*websocket.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maxRacers := s.cfg.Multiple.MaxRacer
	if len(s.racers) > maxRacers {
		return nil, errors.New("maximum 5 racers allowed")

	}
	if len(s.racers) == maxRacers {
		return &s.connections, nil
	}
	//body, err := utils.ValidateToken(token)
	//if err != nil {
	//	return err
	//}
	//id, err := uuid.Parse(body.ID)
	//if err != nil {
	//	return err
	//}

	s.racers = append(s.racers, token)
	s.connections = append(s.connections, conn)

	t := s.cfg.Multiple.Timer
	if _, exists := s.timers[link]; !exists {
		s.timers[link] = t
	}

	if len(s.racers) == 1 {
		s.d.createdAt = time.Now()
	}

	return &s.connections, nil
}

func (s *service) Timer(link string, cons *[]*websocket.Conn) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(*cons) < 2 {
		return -1, errors.New("waiting for other competitors")
	}

	if timer, ok := s.timers[link]; !ok || timer == 0 {
		delete(s.timers, link)
		return -1, errors.New("timer is over")
	}

	time.Sleep(1 * time.Second)
	s.timers[link]--

	return s.timers[link], nil
}

// WhiteLine means the line where racers starts the race. Where timer is zero and race started
func (s *service) WhiteLine(ctx context.Context, link string) error {
	var mlt models.MultipleRace

	l, err := uuid.Parse(link)
	if err != nil {
		log.Println(err.Error())
		return errors.New("link incorrect")
	}

	mlt.GeneratedLink = l
	mlt.CreatedAt = s.d.createdAt
	mlt.Racers = s.racers
	mlt.CreatorId = s.racers[0]
	mlt.Text = s.d.textID
	fmt.Println(mlt.Text, "MLT_TEXT")
	err = s.repo.Multiple.AddRacers(ctx, mlt)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (s *service) RandomText(ctx context.Context) (string, error) {
	ids, err := s.repo.Multiple.Texts(ctx)
	if err != nil {
		return "", errors.New("unable to get text")
	}
	id := randomize(ids)
	s.d.textID = id
	fmt.Println(s.d.textID, "SDTEXTID")
	text, err := s.repo.Multiple.Text(ctx, id)
	if err != nil {
		return "", errors.New("unable to get text")
	}
	return text, nil
}

func randomize(ids []uuid.UUID) uuid.UUID {
	if len(ids) == 0 {
		return uuid.Nil
	}
	return ids[rand.Intn(len(ids))] // Select a random UUID from the slice
}

func NewMultiple(repo *repository.Repo, cfg *config.Config) Racer {
	return &service{
		repo:        repo,
		cfg:         cfg,
		timers:      make(map[string]int),
		racers:      make([]string, 0, 5),
		connections: make([]*websocket.Conn, 0, 5),
		d:           data{},
	}
}

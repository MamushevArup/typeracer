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
	Join(token string, conn *websocket.Conn, link string) (*[]*websocket.Conn, models.RacerM, error)
	Timer(link string, cons *[]*websocket.Conn) (int, error)
	WhiteLine(ctx context.Context, link string) error
	CurrentSpeed(racer *models.RacerCurrentWpm) (models.RacerSpeed, error)
}

var ErrorWaitingRacers = errors.New("waiting for other competitors")

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

func (s *service) Join(id string, conn *websocket.Conn, link string) (*[]*websocket.Conn, models.RacerM, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxRacers := s.cfg.Multiple.MaxRacer
	if len(s.racers) > maxRacers {
		return nil, models.RacerM{}, errors.New("maximum 5 racers allowed")

	}
	if len(s.racers) == maxRacers {
		return &s.connections, models.RacerM{}, nil
	}

	s.racers = append(s.racers, id)
	s.connections = append(s.connections, conn)

	t := s.cfg.Multiple.Timer
	if _, exists := s.timers[link]; !exists {
		s.timers[link] = t
	}

	if len(s.racers) == 1 {
		s.d.createdAt = time.Now()
	}

	parseId, err := uuid.Parse(id)
	if err != nil {
		tempUid, err2 := uuid.NewUUID()
		if err2 != nil {
			log.Println(err2)
			return nil, models.RacerM{}, err2
		}
		return &s.connections, models.RacerM{
			Role:  "guest",
			Email: tempUid.String(),
		}, nil
	}

	user, err := s.repo.Multiple.User(ctx, parseId)
	if err != nil {
		return &s.connections, models.RacerM{}, err
	}

	log.Println("Successfully return value")

	return &s.connections, user, nil
}

func (s *service) CurrentSpeed(racer *models.RacerCurrentWpm) (models.RacerSpeed, error) {
	if racer.Email == "" {
		return models.RacerSpeed{}, errors.New("email is empty")
	}
	if racer.Duration <= 0 {
		return models.RacerSpeed{}, errors.New("duration is less than or equal to 0")
	}
	if racer.Index < 0 {
		return models.RacerSpeed{}, errors.New("current symbol is less than 0")

	}
	var racerSpeed models.RacerSpeed

	wpm := countWPM(racer.Index, racer.Duration)

	racerSpeed.Email = racer.Email
	racerSpeed.Wpm = int(wpm)

	return racerSpeed, nil
}

func countWPM(length, duration int) float64 {
	const averageWordLength = 5 // Assuming an average word length of 5 characters

	// Calculate total words in the text
	totalWords := length / averageWordLength

	// Calculate WPM
	wpm := float64(totalWords) / (float64(duration) / 60.0)

	return wpm
}

func (s *service) Timer(link string, cons *[]*websocket.Conn) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(*cons) < 2 {
		return -1, ErrorWaitingRacers
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

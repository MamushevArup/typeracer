package race

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var ErrorWaitingRacers = errors.New("waiting for other competitors")

const ctxTimeout = 3 * time.Second

func (s *service) Join(id string, conn *websocket.Conn, link string) (*[]*websocket.Conn, models.RacerM, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	maxRacers := s.cfg.Multiple.MaxRacer
	if len(s.racers) > maxRacers {
		return nil, models.RacerM{}, fmt.Errorf("max racers reached %v, track link %v, wants to join %v", maxRacers, link, id)

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

	// convert racer to uuid if it is guest generate temporary uuid for guest
	parseId, err := uuid.Parse(id)
	if err != nil {

		tempUid, err2 := uuid.NewUUID()
		if err2 != nil {
			log.Println(err2)
			return nil, models.RacerM{}, fmt.Errorf("unable to generate uuid for guest %w", err2)
		}

		ava, err := s.repo.Multiple.GuestAvatar(ctx)
		if err != nil {
			return nil, models.RacerM{}, fmt.Errorf("unable to get guest avatar due to %w", err)
		}

		return &s.connections, models.RacerM{
			Email:    tempUid.String(),
			Username: "guest",
			Avatar:   ava,
		}, nil

	}

	user, err := s.repo.Multiple.User(ctx, parseId)
	if err != nil {
		return &s.connections, models.RacerM{}, fmt.Errorf("fail to get user info user=%v, link=%v, err=%w", id, link, err)
	}

	log.Println("Successfully return value")

	return &s.connections, user, nil
}

// WhiteLine means the line where racers starts the race. Where timer is zero and race started
func (s *service) WhiteLine(link string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout+(1*time.Second))
	defer cancel()

	var mlt models.MultipleRace

	l, err := uuid.Parse(link)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to parse link %v, err=%w", link, err)
	}

	mlt = s.populateWhiteLine(mlt, l)

	err = s.repo.Multiple.AddRacers(ctx, mlt)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("fail to add racers link=%v, user=%v, err=%w", link, mlt.CreatorId, err)
	}

	return nil
}

func (s *service) populateWhiteLine(mlt models.MultipleRace, l uuid.UUID) models.MultipleRace {

	mlt.GeneratedLink = l
	mlt.CreatedAt = s.d.createdAt
	mlt.Racers = s.racers
	mlt.CreatorId = s.racers[0]
	mlt.Text = s.d.textID

	return mlt
}

func (s *service) RandomText(ctx context.Context, racerID string) (string, error) {

	textUUIDS, err := s.repo.Multiple.Texts(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get texts user:%v, err=%w", racerID, err)
	}

	textUUID := randomize(textUUIDS)

	s.d.textID = textUUID

	text, err := s.repo.Multiple.Text(ctx, textUUID)
	if err != nil {
		return "", fmt.Errorf("fail to get text text:%v, user:%v, err=%w", textUUID, racerID, err)
	}

	return text, nil
}

func (s *service) Timer(link string, cons *[]*websocket.Conn) (int, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(*cons) < 2 {
		return -1, ErrorWaitingRacers
	}

	if timer, ok := s.timers[link]; !ok || timer == 0 {
		delete(s.timers, link)
		return -1, fmt.Errorf("timer is over")
	}

	time.Sleep(1 * time.Second)
	s.timers[link]--

	return s.timers[link], nil
}

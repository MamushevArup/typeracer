package multiple

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Racer interface {
	Link(id string) (uuid.UUID, error)
	ConnectWS(conn *websocket.Conn) *models.RacerDTO
	CheckLink(ctx context.Context, link string) error
	KillLink(ticker *time.Ticker)
}

type service struct {
	repo *repository.Repo
	mlt  *models.MultipleRace
	hub  *models.Hub
	mu   sync.Mutex
}

func (s *service) KillLink(ticker *time.Ticker) {
	for range ticker.C {
		err := s.repo.Multiple.CleanLink(context.TODO(), time.Now())
		if err != nil {
			log.Printf("error cleaning expired links %v\n", err)
		}
	}
}

func (s *service) CheckLink(ctx context.Context, link string) error {
	l, err := uuid.Parse(link)
	if err != nil {
		return err
	}
	ex, err := s.repo.Multiple.Link(ctx, l)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New("link doesn't exist")
	}

	return nil
}

type client struct {
	*models.RacerDTO
	hub *models.Hub
}

var cl client

func (s *service) ConnectWS(conn *websocket.Conn) *models.RacerDTO {
	msgChan := make(chan *models.Message)
	userId := s.mlt.CreatorId.String()
	s.mu.Lock()
	s.hub.Racers[userId] = true
	s.mu.Unlock()
	msg := &models.Message{
		Text: "connected",
		Wpm:  0,
	}
	msgChan <- msg
	racer := &models.RacerDTO{
		Conn:     conn,
		Username: "Arup",
		Message:  msgChan,
	}
	s.hub.Broadcast <- msg
	s.hub.Register <- racer

	fmt.Println(s.hub.Broadcast)

	go cl.writeMessage()
	go cl.readMessage()

	go s.Run()

	return racer
}

func (s *client) writeMessage() {
	defer func() {
		s.Conn.Close()
	}()
	for {
		message, ok := <-s.Message
		if !ok {
			return
		}
		err := s.Conn.WriteJSON(message)
		if err != nil {
			return
		}
	}
}

func (s *client) readMessage() {
	defer func() {
		s.hub.UnRegister <- s.RacerDTO
		s.Conn.Close()
	}()

	for {
		_, m, err := s.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err, "readMessage UNEXPECTEDCLOSEERROR")
			}
			return
		}
		msg := &models.Message{
			Text: string(m),
			Wpm:  0,
		}
		s.hub.Broadcast <- msg
	}
}

func (s *service) Run() {
	for {
		select {
		case m := <-s.hub.Broadcast:
			s.mu.Lock()
			if _, ok := s.hub.Racers[s.mlt.GeneratedLink.String()]; ok {
				s.hub.RacerConn.Message <- m
			}
			s.mu.Unlock()
		}
	}
}

func (s *service) Link(id string) (uuid.UUID, error) {
	trackId := uuid.New()
	s.mlt.GeneratedLink = trackId
	uid, err := uuid.Parse(id)
	if err != nil {
		return [16]byte{}, err
	}
	s.mlt.CreatorId = uid
	s.mlt.CreatedAt = time.Now()

	err = s.repo.Multiple.InsertLink(context.TODO(), trackId, id, time.Now())
	if err != nil {
		return [16]byte{}, err
	}
	return trackId, nil
}

func NewMultiple(repo *repository.Repo) Racer {
	return &service{
		repo: repo,
		mlt:  new(models.MultipleRace),
		hub:  models.NewHub(),
	}
}

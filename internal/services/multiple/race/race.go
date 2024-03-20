package race

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"math/rand"
)

type WebSocketConnection interface {
	SendMessage(message *models.RacerDTO) error
	ReceiveMessage() (*models.RacerDTO, error)
}

type Racer interface {
	ConnectWS(conn WebSocketConnection) *models.RacerDTO
	RandomText(ctx context.Context) (string, error)
}

type WebSocketConn struct {
	Conn *websocket.Conn
}

func (wsc *WebSocketConn) SendMessage(message *models.RacerDTO) error {
	return wsc.Conn.WriteJSON(message)
}

func (wsc *WebSocketConn) ReceiveMessage() (*models.RacerDTO, error) {
	var message models.RacerDTO
	err := wsc.Conn.ReadJSON(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

type service struct {
	repo *repository.Repo
	hub  *models.Hub
}

func (s *service) ConnectWS(conn WebSocketConnection) *models.RacerDTO {
	// Use the conn object to send and receive messages
	// For example:
	message, err := conn.ReceiveMessage()
	if err != nil {
		// Handle error
	}
	err = conn.SendMessage(message)
	if err != nil {
		// Handle error
	}
	return nil
}

func (s *service) RandomText(ctx context.Context) (string, error) {
	ids, err := s.repo.Multiple.Texts(ctx)
	if err != nil {
		return "", errors.New("unable to get text")
	}
	id := randomize(ids)
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

func NewMultiple(repo *repository.Repo) Racer {
	return &service{
		repo: repo,
		hub:  models.NewHub(),
	}
}

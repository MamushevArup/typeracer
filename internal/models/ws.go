package models

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Racers     map[string]bool `json:"racer"`
	Broadcast  chan *Message   `json:"broadcast"`
	Register   chan *RacerDTO  `json:"register"`
	UnRegister chan *RacerDTO  `json:"un_register"`
	RacerConn  *RacerDTO       `json:"racer_conn"`
}

type RacerDTO struct {
	Conn     *websocket.Conn `json:"conn"`
	Username string          `json:"username"`
	Message  chan *Message   `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Wpm  int    `json:"wpm"`
}

func NewHub() *Hub {
	return &Hub{
		Racers:     make(map[string]bool),
		Broadcast:  make(chan *Message, 5),
		Register:   make(chan *RacerDTO),
		UnRegister: make(chan *RacerDTO),
	}
}

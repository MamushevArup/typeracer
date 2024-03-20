package models

type Hub struct {
	Racers     map[string]bool `json:"racer"`
	Broadcast  chan *Message   `json:"broadcast"`
	Register   chan *RacerDTO  `json:"register"`
	UnRegister chan *RacerDTO  `json:"un_register"`
	RacerConn  *RacerDTO       `json:"racer_conn"`
}

type RacerDTO struct {
	Username string `json:"username"`
	Duration int    `json:"duration"`
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

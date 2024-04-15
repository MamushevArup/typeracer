package models

import (
	"github.com/google/uuid"
	"time"
)

// Response for start Multiple Race
type MultipleRace struct {
	GeneratedLink uuid.UUID `json:"generated_link"`
	CreatorId     string    `json:"racer_id"`
	TrackName     string    `json:"track_name"`
	CreatedAt     time.Time `json:"created_at"`
	Racers        []string  `json:"racers"`
	Text          uuid.UUID `json:"text"`
}

type IncomingMessage struct {
	Data interface{}
}

type RacerM struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

type RacerSpeed struct {
	Email string `json:"email"`
	Wpm   int    `json:"wpm"`
}

type RacerCurrentWpm struct {
	Email    string `json:"email"`
	Duration int    `json:"duration"`
	Index    int    `json:"index"`
}

type RaceResult struct {
	RacerId  uuid.UUID `json:"racer_id"`
	Email    string    `json:"email"`
	Place    int       `json:"place"`
	Accuracy int       `json:"accuracy"`
	Time     int       `json:"time"`
	Wpm      int       `json:"wpm"`
}

type RaceEndRequest struct {
	Type     int    `json:"type"`
	Errors   int    `json:"errors"`
	Length   int    `json:"length"`
	Email    string `json:"email"`
	Duration int    `json:"duration"`
}

// RacerRepoM model between service and repository
type RacerRepoM struct {
	GeneratedLink uuid.UUID `db:"generated_link"`
	RacerId       uuid.UUID `db:"racer_id"`
	Duration      int       `db:"duration"`
	Wpm           int       `db:"wpm"`
	Accuracy      int       `db:"accuracy"`
	StartTime     time.Time `db:"start_time"`
	// winner is string due to the fact that we can have a guest racer which win the race
	Winner    string `db:"winner"`
	Place     int    `db:"place"`
	TrackSize int    `db:"track_size"`
}

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

type MultipleRaceDTO struct {
	GeneratedLink uuid.UUID   `json:"generated_link"`
	TrackName     string      `json:"track_name"`
	Racers        []uuid.UUID `json:"racers"`
	Text          string      `json:"text"`
}

type MultipleSession struct {
	GeneratedLink uuid.UUID `json:"generated_link"`
	RacerId       uuid.UUID `json:"racer_id"`
	StartTime     time.Time `json:"start_time"`
	Duration      int       `json:"duration"`
	Wpm           int       `json:"wpm"`
	Accuracy      float64   `json:"accuracy"`
	Winner        string    `json:"winner"`
	Place         int       `json:"place"`
	TrackSize     int       `json:"track_size"`
}

type MultipleSessionDTO struct {
	Duration int     `json:"duration"`
	Wpm      int     `json:"wpm"`
	Accuracy float64 `json:"accuracy"`
	Place    int     `json:"place"`
}

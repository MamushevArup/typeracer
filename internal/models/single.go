package models

import (
	"github.com/google/uuid"
	"time"
)

// Single represents the single table in PostgreSQL.
type Single struct {
	ID      uuid.UUID `json:"-" db:"id"`
	RacerID uuid.UUID `json:"racer_id"`

	Text            string `json:"text" db:"content"`
	TextLen         int    `json:"text_len" db:"length"`
	TextAuthor      string `json:"text_author" db:"author"`
	ContributorName string `json:"contributor_name" db:"contributor"`

	RacerName string `json:"racer_name" db:"username"`
	Avatar    string `json:"avatar" db:"avatar"`
}

type EndSingle struct {
	TextId    uuid.UUID
	RacerId   uuid.UUID
	RaceId    uuid.UUID
	Speed     int       `json:"speed" db:"speed"`
	Duration  int       `json:"duration" db:"duration"`
	Accuracy  float64   `json:"accuracy" db:"accuracy"`
	StartTime time.Time `db:"start_time"`
}

// RaceHistory represents the race_history table in PostgreSQL.
type RaceHistory struct {
	SingleID uuid.UUID `json:"single_id" db:"single_id"`
	RacerID  uuid.UUID `json:"racer_id" db:"racer_id"`
	TextID   uuid.UUID `json:"text_id" db:"text_id"`
}

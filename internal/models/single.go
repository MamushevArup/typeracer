package models

import (
	"github.com/google/uuid"
	"time"
)

// Single represents the single table in PostgreSQL.
type Single struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Speed     int       `json:"speed" db:"speed"`
	Duration  int       `json:"duration" db:"duration"`
	Accuracy  float64   `json:"accuracy" db:"accuracy"`
	StartTime time.Time `db:"start_time"`
	RacerID   uuid.UUID `json:"racer_id" db:"racer_id"`
	TextID    uuid.UUID `json:"text_id" db:"text_id"`
}

// RaceHistory represents the race_history table in PostgreSQL.
type RaceHistory struct {
	SingleID uuid.UUID `json:"single_id" db:"single_id"`
	RacerID  uuid.UUID `json:"racer_id" db:"racer_id"`
	TextID   uuid.UUID `json:"text_id" db:"text_id"`
}

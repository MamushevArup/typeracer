package models

import (
	"github.com/google/uuid"
	"time"
)

// Single represents the single table in PostgreSQL.
type Single struct {
	ID              uuid.UUID `json:"-" db:"id"`
	RacerID         uuid.UUID `json:"-"`
	TextID          uuid.UUID `json:"-"`
	Text            string    `json:"text" db:"content"`
	TextLen         int       `json:"text_len" db:"length"`
	TextAuthor      string    `json:"text_author" db:"author"`
	ContributorName string    `json:"contributor_name" db:"contributor_name"`
	RacerName       string    `json:"racer_name" db:"username"`
	Avatar          string    `json:"avatar" db:"avatar"`
}

type RespEndSingle struct {
	RacerId     uuid.UUID `json:"-"`
	Wpm         int       `json:"wpm" db:"speed"`
	Accuracy    float64   `json:"accuracy" db:"accuracy"`
	Duration    int       `json:"duration" db:"duration"`
	StartedTime time.Time `json:"-" db:"start_time"`
}

type ReqEndSingle struct {
	RacerId  uuid.UUID `json:"-"`
	Duration int       `json:"duration"`
	Errors   int       `json:"errors"`
}

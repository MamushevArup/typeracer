package models

import (
	"github.com/google/uuid"
	"time"
)

// Text represents the text table in PostgreSQL.
type Text struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Content       string    `json:"content" db:"content"`
	Author        string    `json:"author" db:"author"`
	Occurrence    int       `json:"occurrence" db:"occurrence"`
	AcceptedAt    time.Time `db:"accepted_at"`
	Length        int       `json:"length" db:"length"`
	AvgSpeed      int       `json:"avg_speed" db:"avg_speed"`
	ContributorID uuid.UUID `json:"contributor_id" db:"contributor_id"`
}

// RandomText represents the random_text table in PostgreSQL.
type RandomText struct {
	TextID uuid.UUID `db:"text_id"`
}

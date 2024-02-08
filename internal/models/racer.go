package models

import (
	"github.com/google/uuid"
	"time"
)

type Racer struct {
	Id            uuid.UUID `json:"-" db:"id"`
	Email         string    `json:"email" db:"email"`
	Password      string    `db:"password"`
	Username      string    `json:"username" db:"username"`
	Avatar        string    `json:"avatar" db:"avatar"`
	Country       string    `json:"country" db:"country"`
	CreatedAt     time.Time `db:"created_at"`
	LastLogin     time.Time `db:"last_login"`
	Races         int       `json:"races" db:"races"`
	AvgSpeed      int       `json:"avg_speed" db:"avg_speed"`
	LastRaceSpeed int       `json:"last_race_speed" db:"last_race_speed"`
	BestSpeed     int       `json:"best_speed" db:"best_speed"`
	Theme         bool      `json:"theme" db:"theme"`
	RefreshToken  string    `json:"refresh_token" db:"refresh_token"`
	Role          string    `json:"role" db:"role"`
}

type RacerAuth struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	Username     string    `db:"username"`
	RefreshToken string    `db:"refresh_token"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	LastLogin    time.Time `db:"last_login"`
	Fingerprint  string
}

// Contributor represents the contributor table in PostgreSQL.
type Contributor struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	SentAt time.Time `json:"sent_at" db:"sent_at"`
	TextID uuid.UUID `json:"text_id" db:"text_id"`
}
type ContributeText struct {
	RacerID     uuid.UUID `json:"racer_id"`
	TextID      uuid.UUID `json:"-"`
	Content     string    `json:"content"`
	Length      int       `json:"-"`
	Author      string    `json:"author"`
	Source      string    `json:"source"`
	SourceTitle string    `json:"source_title"`
	SentAt      time.Time `json:"sent_at"`
}

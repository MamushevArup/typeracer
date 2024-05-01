package models

import (
	"github.com/google/uuid"
	"time"
)

type RacerHandler struct {
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	CreatedAt     string `json:"created_at"`
	AvgSpeed      int    `json:"avg_speed"`
	LastRaceSpeed int    `json:"last_race_speed"`
	BestSpeed     int    `json:"best_speed"`
	Races         int    `json:"races"`
}

type RacerRepository struct {
	Username      string    `db:"username"`
	Avatar        string    `db:"avatar"`
	CreatedAt     time.Time `db:"created_at"`
	AvgSpeed      int       `db:"avg_speed"`
	LastRaceSpeed int       `db:"last_race_speed"`
	BestSpeed     int       `db:"best_speed"`
	Races         int       `db:"races"`
}

type Avatar struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type AvatarUpdateRepo struct {
	Id      int
	RacerID uuid.UUID
}

type AvatarUpdate struct {
	Id      int    `json:"id"`
	RacerId string `json:"-"`
}

type RacerUpdate struct {
	Id       string
	Email    string `json:"email" valid:"email,optional"`
	Username string `json:"username" valid:"optional~Username is required,matches(^[a-zA-Z0-9]+$)~Username must consist only of ASCII letters and digits"`
}

type RacerUpdateRepo struct {
	Id       uuid.UUID
	Email    string
	Username string
}

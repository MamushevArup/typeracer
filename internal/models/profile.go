package models

import "time"

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
	Url string `db:"url" json:"url"`
}

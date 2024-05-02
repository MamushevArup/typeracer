package models

import (
	"github.com/google/uuid"
	"time"
)

type SingleHistory struct {
	SingleID  uuid.UUID `json:"id"`
	Speed     int       `json:"speed"`
	Duration  int       `json:"duration"`
	Accuracy  float64   `json:"accuracy"`
	StartedAt time.Time `json:"started_at"`
}

type SingleHistoryHandler struct {
	SingleID  uuid.UUID `json:"id"`
	Speed     int       `json:"speed"`
	Duration  int       `json:"duration"`
	Accuracy  float64   `json:"accuracy"`
	StartedAt string    `json:"started_at"`
}

type SingleHistoryText struct {
	Content     string `json:"content"`
	Author      string `json:"author"`
	Source      string `json:"source"`
	SourceTitle string `json:"source_title"`
	Contributor string `json:"contributor"`
}

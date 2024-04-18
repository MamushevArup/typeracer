package models

import (
	"github.com/google/uuid"
	"time"
)

type ContributeHandlerRequest struct {
	RacerID     string `json:"-"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	Source      string `json:"source"`
	SourceTitle string `json:"source_title"`
}

type ContributeServiceRequest struct {
	ModerationID uuid.UUID
	RacerID      uuid.UUID
	Content      string
	Author       string
	Source       string
	SourceTitle  string
	SentAt       time.Time
	Status       int
	Length       int
	ContentHash  string
}

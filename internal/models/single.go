package models

import (
	"github.com/google/uuid"
	"time"
)

// EndRace models

type RespEndSingle struct {
	RacerId     uuid.UUID `json:"-"`
	Wpm         int       `json:"wpm" db:"speed"`
	Accuracy    float64   `json:"accuracy" db:"accuracy"`
	Duration    int       `json:"duration" db:"duration"`
	StartedTime time.Time `json:"-" db:"start_time"`
	RaceId      uuid.UUID `json:"-" db:"id"`
	TextId      uuid.UUID `json:"-" db:"text_id"`
}

type ReqEndSingle struct {
	RacerId  uuid.UUID `json:"-"`
	Duration int       `json:"duration"`
	Errors   int       `json:"errors"`
	Length   int       `json:"length"`
}

// SingleResponse starts race accumulates two structs RacerInfo and TextInfo
type SingleResponse struct {
	RacerInfo `json:"racer"`
	TextInfo  `json:"text"`
}

type TextInfo struct {
	TextID uuid.UUID `json:"-" db:"id"`
	// content is actual text
	Content string `json:"content" db:"content"`
	// source from what it is coming. It can be from a book, article, etc.
	Source string `json:"source" db:"source"`
	// header the title of the text. Ex name of the book, song
	Header string `json:"header" db:"source_title"`
	// who wrote the content
	Author string `json:"text_author" db:"author"`
	// who contributed the content
	ContributorName string `json:"contributor_name" db:"contributor_name"`
}

type RacerInfo struct {
	Username string `json:"username" db:"username"`
	Avatar   string `json:"avatar" db:"avatar"`
}

type CountWpm struct {
	CurrIdx  int `json:"index"`
	Duration int `json:"duration"`
}

type Speed struct {
	Wpm int `json:"wpm"`
}

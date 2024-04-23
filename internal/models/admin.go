package models

import (
	"github.com/google/uuid"
	"time"
)

type AdminSignIn struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type AdminSignInResponse struct {
	Access string `json:"access"`
}

type ModerationRepoResponse struct {
	ModerationID uuid.UUID `json:"moderation_id"`
	Username     string    `json:"username"`
	SentAt       time.Time `json:"sent_at"`
}

type ModerationServiceResponse struct {
	ModerationID    uuid.UUID `json:"moderation_id"`
	ContributorName string    `json:"contributor_name"`
	SentAt          string    `json:"sent_at"`
}

type ModerationTextDetails struct {
	ModerationID uuid.UUID `json:"moderation_id"`
	Content      string    `json:"content"`
	Author       string    `json:"author"`
	Source       string    `json:"source"`
	SourceTitle  string    `json:"source_title"`
}

type ModerationApprove struct {
	ModerationID uuid.UUID `json:"moderation_id"`
	RacerID      uuid.UUID `json:"racer_id"`
	Content      string    `json:"content"`
	Author       string    `json:"author"`
	Length       int       `json:"length"`
	Source       string    `json:"source"`
	SourceTitle  string    `json:"source_title"`
	SentAt       time.Time `json:"sent_at"`
}

type ApproveToText struct {
	TextID        uuid.UUID `db:"id"`
	ContributorID uuid.UUID `db:"contributor_id"`
	Content       string    `db:"content"`
	Author        string    `db:"author"`
	AcceptedAt    time.Time `db:"accepted_at"`
	Length        int       `db:"length"`
	Source        string    `db:"source"`
	SourceTitle   string    `db:"source_title"`
}

type ApproveToContributor struct {
	ContributorID uuid.UUID `db:"user_id"`
	SentAt        time.Time `db:"sent_at"`
	TextID        uuid.UUID `db:"text_id"`
}

type TextAcceptTransaction struct {
	Text        ApproveToText
	Contributor ApproveToContributor
}

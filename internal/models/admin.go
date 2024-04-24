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

type AdminSignInRefresh struct {
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
	TextID        uuid.UUID
	ContributorID uuid.UUID
	Content       string
	Author        string
	AcceptedAt    time.Time
	Length        int
	Source        string
	SourceTitle   string
}

type ApproveToContributor struct {
	ContributorID uuid.UUID
	SentAt        time.Time
	TextID        uuid.UUID
}

type TextAcceptTransaction struct {
	Text        ApproveToText
	Contributor ApproveToContributor
}

type ModerationRejectToService struct {
	ModerationID string
	Reason       string `json:"reason"`
}

type ModerationRejectToRepo struct {
	ModerationID uuid.UUID
	Reason       string
}

package models

import (
	"github.com/google/uuid"
	"time"
)

type RacerAuth struct {
	ID           uuid.UUID `db:"user_id"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	Username     string    `db:"username"`
	RefreshToken string    `db:"refresh_token"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	LastLogin    time.Time `db:"last_login"`
	Fingerprint  string
}

type AuthResponse struct {
	Access string `json:"access"`
}

type SignIn struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type SignUp struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type RefreshS struct {
	Fingerprint string `json:"fingerprint"`
}

type SignUpService struct {
	Access  string
	Avatar  string
	Refresh string
}

type SignUpHandler struct {
	Access string `json:"access"`
	Avatar string `json:"avatar"`
}

type SignInService struct {
	Access   string
	Refresh  string
	Username string
	Avatar   string
}

type SignInHandler struct {
	Access   string `json:"access"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

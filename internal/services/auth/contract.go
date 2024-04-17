package auth

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	accessTokenExpiration  = 30 * time.Minute
	refreshTokenExpiration = 60 * 24 * time.Hour // 60 days
	role                   = "racer"
)

var signingMethod = jwt.SigningMethodHS256

type Auth interface {
	SignIn(ctx context.Context, email, password, fingerprint string) (string, string, error)
	SignUp(ctx context.Context, email, username, password, fingerprint string) (string, string, error)
	CheckUserSignUp(ctx context.Context, email, password string) error
	RefreshToken(ctx context.Context, refresh, fingerprint string) (string, string, error)
	Logout(ctx context.Context, refresh string) error
}

type auth struct {
	repo *repository.Repo
}

func NewAuth(repo *repository.Repo) Auth {
	return &auth{repo: repo}
}

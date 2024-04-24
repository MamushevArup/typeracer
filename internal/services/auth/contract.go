package auth

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/models"
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
	SignIn(ctx context.Context, email, password, fingerprint string) (*models.SignInService, error)
	SignUp(ctx context.Context, email, username, password, fingerprint string) (models.SignUpService, error)
	CheckUserSignUp(ctx context.Context, email string) error
	RefreshToken(ctx context.Context, refresh, fingerprint string) (string, string, error)
	Logout(ctx context.Context, refresh string) error
	AdminRefresh(ctx context.Context, refresh string) (string, string, error)
}

type auth struct {
	repo *repository.Repo
	cfg  *config.Config
}

func NewAuth(repo *repository.Repo, cfg *config.Config) Auth {
	return &auth{repo: repo, cfg: cfg}
}

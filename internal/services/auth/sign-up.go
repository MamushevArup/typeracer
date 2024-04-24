package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

type tokenClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (a *auth) SignUp(ctx context.Context, email, username, password, fingerprint string) (models.SignUpService, error) {
	// Sign up accept email, username, password after pass validation and check user existence if user doesn't exist
	// we provide endpoint/refresh token with role and id inside
	// endpoint token live 30 minute refresh token live 60 day
	// refresh token should store in the cookies
	// in the service we should hash the password and compare it
	var racer models.RacerAuth

	userId, err := uuid.NewUUID()
	if err != nil {
		return models.SignUpService{}, fmt.Errorf("unable to generate uuid %w", err)
	}

	email = strings.ToLower(email)
	username = strings.ToLower(username)
	password = strings.ToLower(password)

	// cipher password using bcrypt
	hashPasswd, err := a.generateHashPassword(password)
	if err != nil {
		return models.SignUpService{}, fmt.Errorf("%w", err)
	}

	access, err := a.generateAccessToken(userId.String(), role)
	if err != nil {
		return models.SignUpService{}, fmt.Errorf("%w", err)
	}

	refresh, err := a.generateRefreshToken()
	if err != nil {
		return models.SignUpService{}, fmt.Errorf("%w", err)
	}

	racer.ID = userId
	racer.Email = email
	racer.Username = username
	racer.Password = hashPasswd
	racer.RefreshToken = refresh
	racer.Role = role
	racer.CreatedAt = time.Now()
	racer.LastLogin = time.Now()
	racer.Fingerprint = fingerprint

	err = a.repo.Auth.InsertUser(ctx, racer)
	if err != nil {
		return models.SignUpService{}, fmt.Errorf("%w", err)
	}

	avatarImage, err := a.repo.Auth.UserAvatar(ctx, email)
	if err != nil {
		return models.SignUpService{}, fmt.Errorf("%w", err)
	}

	res := models.SignUpService{
		Access:  access,
		Avatar:  avatarImage,
		Refresh: refresh,
	}

	return res, nil
}

func (a *auth) CheckUserSignUp(ctx context.Context, email string) error {
	// check user exist if so redirect to the sign in else generate token
	email = strings.ToLower(email)

	byEmail, err := a.repo.Auth.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%w", err)
	}

	if byEmail == nil {
		return fmt.Errorf("account with this email already created use sign-in")
	}

	return nil
}

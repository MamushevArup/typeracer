package auth

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

func (a *auth) SignIn(ctx context.Context, email, password, fingerprint string) (string, string, error) {

	var r models.RacerAuth

	// to generalize all password and email will store in lowercase
	password = strings.ToLower(password)
	email = strings.ToLower(email)

	bytePass := []byte(password)

	byEmail, err := a.repo.Auth.UserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("unable find user by email=%v, err=%w", email, err)
	}

	if !byEmail {
		return "", "", fmt.Errorf("user do not exist sign up first")
	}

	// here I need to make repo call to get password by
	id, token, pwd, err := a.repo.Auth.GetUserPasswordByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("internal error due to password checking err=%w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwd), bytePass)
	if err != nil {
		return "", "", fmt.Errorf("password does not match err=%w", err)
	}

	isActive, err := a.repo.Auth.UserSession(ctx, token, id)
	if err != nil {
		return "", "", fmt.Errorf("unable find user session err=%w", err)
	}

	r.Fingerprint = fingerprint

	if isActive {
		err = a.repo.Auth.DeleteSession(ctx, r.Fingerprint, token)
		if err != nil {
			return "", "", fmt.Errorf("%w", err)
		}
	}

	access, err := a.generateAccessToken(id.String(), role)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	refresh, err := a.generateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	r.ID = id
	r.Role = role
	r.Email = email
	r.Password = pwd
	r.RefreshToken = refresh
	r.LastLogin = time.Now()

	err = a.repo.Auth.InsertSession(ctx, r)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	return access, refresh, nil
}

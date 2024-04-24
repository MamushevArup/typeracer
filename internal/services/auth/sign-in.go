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

func (a *auth) SignIn(ctx context.Context, email, password, fingerprint string) (*models.SignInService, error) {

	adminUsername := a.cfg.Admin.Username
	adminPassword := a.cfg.Admin.Password

	if email == adminUsername && password == adminPassword {
		access, err := a.generateAccessToken("1", "admin")
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		refresh, err := a.generateRefreshToken()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		err = a.repo.Auth.UpdateAdmin(ctx, adminUsername, refresh)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		adminInfo := &models.SignInService{
			Access:   access,
			Refresh:  refresh,
			Username: "admin",
			Avatar:   "",
		}

		return adminInfo, nil
	}

	var r models.RacerAuth
	// to generalize all password and email will store in lowercase
	password = strings.ToLower(password)
	email = strings.ToLower(email)

	bytePass := []byte(password)

	racerInfo, err := a.repo.Auth.UserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("unable find user by email=%v, err=%w", email, err)
	}

	if racerInfo == nil {
		return nil, fmt.Errorf("user do not exist sign up first")
	}

	// here I need to make repo call to get password by
	id, token, pwd, err := a.repo.Auth.GetUserPasswordByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("internal error due to password checking err=%w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwd), bytePass)
	if err != nil {
		return nil, fmt.Errorf("password does not match err=%w", err)
	}

	isActive, err := a.repo.Auth.UserSession(ctx, token, id)
	if err != nil {
		return nil, fmt.Errorf("unable find user session err=%w", err)
	}

	r.Fingerprint = fingerprint

	if isActive {
		err = a.repo.Auth.DeleteSession(ctx, r.Fingerprint, token)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	access, err := a.generateAccessToken(id.String(), role)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	refresh, err := a.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	r.ID = id
	r.Role = role
	r.Email = email
	r.Password = pwd
	r.RefreshToken = refresh
	r.LastLogin = time.Now()

	err = a.repo.Auth.InsertSession(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	racerInfo.Access = access
	racerInfo.Refresh = refresh

	return racerInfo, nil
}

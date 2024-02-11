package auth

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
	"time"
)

type Auth interface {
	SignIn(ctx context.Context, email, password string) (string, string, error)
	SignUp(ctx context.Context, email, username, password string) (string, string, error)
	CheckUserSignUp(ctx context.Context, email, password string) error
	ValidateToken(tokenString string) (string, error)
	RefreshToken(ctx context.Context, refresh, fingerprint string) (string, string, error)
}

type auth struct {
	repo *repository.Repo
}

type tokenClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuth(repo *repository.Repo) Auth {
	return &auth{repo: repo}
}

const (
	accessTokenExpiration  = 30 * time.Minute
	refreshTokenExpiration = 60 * 24 * time.Hour // 60 days
	refreshTokenCookieName = "refresh_token"
	role                   = "racer"
)

func (a *auth) ValidateToken(tokenString string) (string, error) {
	claim := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}
	if claim.ExpiresAt.Before(time.Now().UTC()) {
		return "", jwt.ErrTokenExpired
	}
	return claims.Role, nil
}

func (a *auth) RefreshToken(ctx context.Context, refresh, fingerprint string) (string, string, error) {
	// user_id, role, fingerprint
	racer, err := a.repo.Auth.Fingerprint(ctx, fingerprint, refresh)
	if err != nil {
		return "", "", err
	}
	// here is the repo layer to remove from session
	err = a.repo.Auth.DeleteSession(ctx, fingerprint, refresh)
	if err != nil {
		return "", "", errors.New("can't remove session " + err.Error())
	}
	// check the token expires or not
	if err = parseRefreshToken(refresh); err != nil {
		return "", "", err
	}
	// check the fingerprint are the same
	if racer.Fingerprint != fingerprint {
		return "", "", errors.New("fingerprint does not match ")
	}
	racer.LastLogin = time.Now()
	racer.RefreshToken, err = a.generateRefreshToken()
	if err != nil {
		return "", "", err
	}
	// or new refresh and store in the session
	//update racer table refresh and fingerprint
	if err = a.repo.Auth.InsertSession(ctx, racer); err != nil {
		return "", "", err
	}
	// create access token and create refresh token and send back to the frontend
	access, err := a.generateAccessToken(racer.ID.String(), racer.Role)
	return access, racer.RefreshToken, nil
}

func parseRefreshToken(refresh string) error {
	parsedToken, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return err
	}
	if !parsedToken.Valid {
		return jwt.ErrSignatureInvalid
	}
	return nil
}

func (a *auth) SignUp(ctx context.Context, email, username, password string) (string, string, error) {
	// Sign up accept email, username, password after pass validation and check user existence if user doesn't exist
	// we provide access/refresh token with role and id inside
	// access token live 30 minute refresh token live 60 day
	// refresh token should store in the cookies
	// in the service we should hash the password and compare it
	var racer models.RacerAuth
	userId := uuid.New()

	email = strings.ToLower(email)
	username = strings.ToLower(username)
	password = strings.ToLower(password)

	// cipher password using bcrypt
	hashPasswd, err := a.generateHashPassword(password)
	if err != nil {
		return "no token", "", err
	}
	access, err := a.generateAccessToken(userId.String(), role)
	if err != nil {
		return "no token", "", err
	}
	refresh, err := a.generateRefreshToken()
	if err != nil {
		return "no token", "", err
	}

	racer.ID = userId
	racer.Email = email
	racer.Username = username
	racer.Password = hashPasswd
	racer.RefreshToken = refresh
	racer.Role = role
	racer.CreatedAt = time.Now()
	racer.LastLogin = time.Now()
	racer.Fingerprint = "1234567890"

	err = a.repo.Auth.InsertUser(ctx, racer)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (a *auth) generateAccessToken(id, role string) (string, error) {
	t := tokenClaims{
		id,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "no token", err
	}
	return "Bearer " + accessToken, nil
}

func (a *auth) generateRefreshToken() (string, error) {
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (a *auth) generateHashPassword(password string) (string, error) {
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPasswd), nil
}

func (a *auth) CheckUserSignUp(ctx context.Context, email, password string) error {
	// check user exist if so redirect to the sign in else generate token
	password = strings.ToLower(password)
	email = strings.ToLower(email)
	byEmail, err := a.repo.Auth.UserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if byEmail {
		return errors.New("user exist use sign in")
	}
	return nil
}

func (a *auth) SignIn(ctx context.Context, email, password string) (string, string, error) {
	var r models.RacerAuth
	password = strings.ToLower(password)
	email = strings.ToLower(email)
	bytePass := []byte(password)

	byEmail, err := a.repo.Auth.UserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if !byEmail {
		return "", "", errors.New("user doesn't exist use sign up")
	}

	// here I need to make repo call to get password by
	id, token, pswd, err := a.repo.Auth.GetUserPasswordByEmail(ctx, email)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(pswd), bytePass)
	if err != nil {
		return "", "", err
	}

	isActive, err := a.repo.Auth.UserSession(ctx, token, id)
	if err != nil {
		return "", "", err
	}
	r.Fingerprint = "1245432134543"

	if isActive {
		err = a.repo.Auth.DeleteSession(ctx, r.Fingerprint, token)
		if err != nil {
			return "", "", err
		}
	}

	access, err := a.generateAccessToken(id.String(), role)
	if err != nil {
		return "", "", err
	}

	refresh, err := a.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	r.ID = id
	r.Role = role
	r.Email = email
	r.Password = pswd
	r.RefreshToken = refresh
	r.LastLogin = time.Now()
	err = a.repo.Auth.InsertSession(ctx, r)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

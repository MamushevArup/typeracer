package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (a *auth) generateAccessToken(id, role string) (string, error) {

	t := tokenClaims{
		id,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(signingMethod, t)

	accessToken, err := token.SignedString([]byte(a.cfg.Jwt.SecretKey))
	if err != nil {
		return "", fmt.Errorf("fail to generate token %w", err)
	}

	return "Bearer " + accessToken, nil
}

func (a *auth) generateRefreshToken() (string, error) {

	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(signingMethod, claim)

	refreshToken, err := token.SignedString([]byte(a.cfg.Jwt.SecretKey))
	if err != nil {
		return "", fmt.Errorf("unable to generate token %w", err)
	}

	return refreshToken, nil
}

func (a *auth) generateHashPassword(password string) (string, error) {

	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to process password %w", err)
	}

	return string(hashPasswd), nil
}

func (a *auth) parseRefreshToken(refresh string) error {

	parsedToken, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.Jwt.SecretKey), nil
	})
	if err != nil {
		return fmt.Errorf("unable to parse token %w", err)
	}

	if !parsedToken.Valid {
		return fmt.Errorf("%w", jwt.ErrSignatureInvalid)
	}

	return nil
}

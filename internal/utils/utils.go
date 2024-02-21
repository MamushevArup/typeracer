package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type tokenClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type claimsBody struct {
	ID   string
	Role string
}

func ValidateToken(tokenString string) (*claimsBody, error) {
	claim := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claim.ExpiresAt.Before(time.Now().UTC()) {
		return nil, jwt.ErrTokenExpired
	}
	return &claimsBody{
		ID:   claims.ID,
		Role: claims.Role,
	}, nil
}

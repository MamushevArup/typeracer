package token_validation

import (
	"errors"
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

var (
	guest      = "guest"
	authHeader = errors.New("auth header empty")
)

type tokenClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func TokenInspector(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		claim := &tokenClaims{}

		token, err := extractToken(c)
		if err != nil {
			if errors.Is(err, authHeader) {
				claim.Role = guest
				c.Set("Role", claim.Role)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, authHeader)
			return
		}

		t, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Jwt.SecretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		if !t.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "token is invalid"})
			return
		}
		claims, ok := t.Claims.(*tokenClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred try again later"})
			return
		}
		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": jwt.ErrTokenExpired})
			return
		}

		c.Set("ID", claims.ID)
		c.Set("Role", claims.Role)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return "", authHeader
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("auth header is not well formatted expect Bearer <token>")
	}
	return parts[1], nil
}

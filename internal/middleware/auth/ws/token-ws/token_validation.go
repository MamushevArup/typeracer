package token_ws

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

var (
	guest = "guest"
)

type tokenClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func TokenVerifier() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Query("token")
		if token == "" {
			c.Set("Role", guest)
			c.Next()
			return
		}

		claim := &tokenClaims{}

		t, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
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

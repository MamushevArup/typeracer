package middleware

import (
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	guest = "guest"
	racer = "racer"
	// soon add admin
)

type Middleware struct {
	service *services.Service
}

func NewMiddleware(service *services.Service) *Middleware {
	return &Middleware{service}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		// Bearer <token> Get only token
		token := strings.Split(auth, " ")[0]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "auth header empty"})
		}
		// this method check for token valid and expiry time and return role and nil if no error
		role, err := m.service.Auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
		}
		if role == guest {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Guest are not allowed to go here"})
			return
		} else if role == racer {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unknown role"})
			return
		}
	}
}

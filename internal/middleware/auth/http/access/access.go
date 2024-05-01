package access

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OnlyGuestOrRacer() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("Role")
		if role != "guest" && role != "racer" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "unauthorized access"})
			return
		}

		c.Next()
	}
}

func OnlyRacer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.MustGet("Role") != "racer" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "unauthorized access"})
			return
		}
	}
}

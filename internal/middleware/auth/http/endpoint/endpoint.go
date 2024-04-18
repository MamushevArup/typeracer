package endpoint

import (
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

const racerRole = "racer"

func Access(svc *services.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("Role")
		if role != racerRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "unauthorized endpoint"})
			return
		}

		ex, err := svc.Single.RacerExists(c, c.MustGet("ID").(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if !ex {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "racer doesn't exist"})
			return
		}

		c.Next()
	}
}

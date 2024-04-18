package admin

import "github.com/gin-gonic/gin"

const admin = "admin"

func Confirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("Role")
		if role != admin {
			c.AbortWithStatusJSON(403, gin.H{"message": "unauthorized access"})
			return
		}
		c.Next()
	}
}

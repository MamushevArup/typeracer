package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Println("error in the handler layer")
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}

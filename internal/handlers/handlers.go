package handlers

import (
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-gonic/gin"
)

// all validation will make here

type Handler interface {
	InitRoutes() *gin.Engine
}

type handler struct {
	service *services.Service
}

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	sgl := router.Group("/single")
	{
		sgl.POST("/race", h.startRace)
		sgl.POST("/end-race", h.endRace)
		sgl.POST("/curr-wpm", h.currWPM)
	}

	return router
}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

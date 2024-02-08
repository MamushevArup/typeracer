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
	router.POST("/contribute", h.contribute)
	router.GET("/moderation/:id", h.moderation)

	auth := router.Group("/api/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/logout", h.logOut)
		auth.POST("/refresh")
	}

	return router
}

func (h *handler) signIn(context *gin.Context) {

}

func (h *handler) logOut(context *gin.Context) {

}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

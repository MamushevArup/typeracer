package handlers

import (
	"github.com/MamushevArup/typeracer/internal/middleware"
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

	// this middleware check for jwt token valid and expiry
	mdl := middleware.NewMiddleware(h.service)

	router.Use(mdl.AuthMiddleware())
	sgl := router.Group("/single")
	{
		sgl.GET("/race", h.startRace)
		sgl.POST("/end-race", h.endRace)
		sgl.POST("/curr-wpm", h.currWPM)
	}
	router.POST("/contribute", h.contribute)
	router.GET("/moderation/:id", h.moderation)

	auth := router.Group("/api/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.DELETE("/logout", h.logOut)
		auth.POST("/refresh", h.refresh)
	}
	// this route stands for create racetrack and start a multiple race
	mlt := router.Group("/track")
	{
		mlt.POST("/link", h.createLink)
		mlt.GET("/race/:id")
	}

	return router
}

func (h *handler) createLink(c *gin.Context) {
	// I get access token and parse it and do not check for the role
	// If settings will required then wait for number of racers
	//
}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

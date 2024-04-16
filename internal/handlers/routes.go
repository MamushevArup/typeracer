package handlers

import (
	validation "github.com/MamushevArup/typeracer/internal/middleware/auth/http/token-validation"
	validationWs "github.com/MamushevArup/typeracer/internal/middleware/auth/ws/token-ws"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

type Handler interface {
	InitRoutes() *gin.Engine
}

type handler struct {
	service *services.Service
}

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Use(validation.TokenInspector())

	sgl := router.Group("/single")
	{
		sgl.GET("/race", h.startRace)
		sgl.POST("/end-race", h.endRace)
		sgl.POST("/curr-wpm", h.currWPM)
	}

	//router.POST("/contribute", h.contribute)
	//router.GET("/moderation/:id", h.moderation)

	auth := router.Group("/api/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.DELETE("/logout", h.logOut)
		auth.POST("/refresh", h.refresh)
	}
	// this route stands for create racetrack and start a multiple race
	mlt := router.Group("/track")
	mlt.Use(validationWs.TokenVerifier())
	{
		mlt.POST("/link", h.createLink)
		// this racetrack will look like this. /race/link?t=<access token>
		mlt.GET("/race/:link", h.raceTrack)
		mlt.DELETE("/race/finish/:id")
	}

	return router
}

package handlers

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/middleware"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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

	router.Use(cors.Default())
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
		// this racetrack will look like this. /race/link?t=<access token>
		mlt.GET("/race/:link", h.raceTrack)
		mlt.DELETE("/race/finish/:id")
	}

	return router
}

func (h *handler) createLink(c *gin.Context) {
	// I get access token and parse it and do not check for the role
	// If settings will required then wait for number of racers
	// here we will create a websocket channel and create a link to join
	// we will redirect to the /race/:id route right there
	// create a unique identifier for race then redirect to the racetrack to start race
	/*
		Header : access token with role, id
	*/
	id, exists := c.Get("id")
	if !exists {
		newErrorResponse(c, http.StatusForbidden, "can't get value from context")
		return
	}
	link, err := h.service.Multiple.Link(id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "/track/race/"+link.String())
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *handler) raceTrack(c *gin.Context) {
	// check link correctness

	// here with created links I will create a room for max 5 racer and open websocket connection
	link := c.Param("link")

	err := h.service.Multiple.CheckLink(context.TODO(), link)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	racer := h.service.Multiple.ConnectWS(conn)
	fmt.Println(racer.Message)
	c.JSON(200, gin.H{"link": "://" + c.Request.Host + c.Request.URL.String()})
}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

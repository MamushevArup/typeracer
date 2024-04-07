package handlers

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/middleware"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

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
	id, _ := authHeader(c)

	link, err := h.service.Link.Create(id.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	text, err := h.service.Multiple.RandomText(context.TODO())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"link":    link,
		"content": text,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *handler) raceTrack(c *gin.Context) {
	// check link correctness

	// here with created links I will create a room for max 5 racer and open websocket connection
	link := c.Param("link")
	id := c.Query("id")
	err := h.service.Link.Check(context.TODO(), link)
	if err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	connections, err := h.service.Multiple.Join(id, conn, link)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {

			t, err := h.service.Multiple.Timer(link, connections)
			if err != nil {
				if errors.Is(err, errors.New("waiting for other competitors")) {
					continue
				}
				err2 := h.service.Multiple.WhiteLine(context.TODO(), link)
				if err2 != nil {
					log.Println(err2.Error())
					return
				}
				log.Println(err.Error())
				return
			}
			sentClients := make(map[*websocket.Conn]bool)
			for _, clientConn := range *connections {
				if _, sent := sentClients[clientConn]; !sent {
					// Send the message to the writer goroutine
					err = clientConn.WriteJSON(t)
					if err != nil {
						log.Println(err)
						return
					}
					sentClients[clientConn] = true
				}
			}
		}
	}()

	//cl.connections = append(cl.connections, conn) // Add the new connection to the slice
	//fmt.Println(cl.connections)
	//// Join a new goroutine for receiving messages
	//messages := make(chan *models.RacerDTO)
	//cl.timer = make(chan int)
	//// Join a new goroutine for receiving messages
	//go func() {
	//	for {
	//		var msg models.RacerDTO
	//		err := conn.ReadJSON(<-cl.timer)
	//		if err != nil {
	//			log.Println(err)
	//			return
	//		}
	//
	//		// Log the received message
	//		log.Println("Received message:", msg)
	//
	//		// Pass the received message to the writing goroutine
	//		messages <- &msg
	//	}
	//}()
	//
	//// Join a new goroutine for sending messages
	//go func() {
	//	for {
	//		// Wait for a message from the reading goroutine
	//		//msg := <-messages
	//		startRace := <-cl.timer
	//		// Broadcast the message to all connected clients
	//		for _, clientConn := range cl.connections {
	//			err := clientConn.WriteJSON(startRace)
	//			if err != nil {
	//				log.Println(err)
	//				return
	//			}
	//		}
	//	}
	//}()
}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

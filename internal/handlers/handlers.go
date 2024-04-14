package handlers

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/MamushevArup/typeracer/internal/middleware/auth/http/token-validation"
	validationWs "github.com/MamushevArup/typeracer/internal/middleware/auth/ws/token-ws"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/MamushevArup/typeracer/internal/services/multiple/race"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
)

// all validation will make here
const (
	currSpeed = iota + 1
	endRace
	leaveRace
)

type Handler interface {
	InitRoutes() *gin.Engine
}

type handler struct {
	service *services.Service
}

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// this middleware check for jwt token valid and expiry

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
	mlt.Use(validationWs.TokenVerifier())
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

	role := c.MustGet("Role")
	id, ex := c.Get("ID")
	if !ex {
		id = role
	}

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

	connections, racer := h.joinRacers(conn, link, id.(string))
	fmt.Println(racer, "RACER")
	currSpeedCh := make(chan models.RacerSpeed)
	errorCh := make(chan error)

	go func() {
		for {
			var msgType models.IncomingMessage
			err = conn.ReadJSON(&msgType.Data)
			if err != nil {
				log.Println(err)
				return
			}
			data, ok := msgType.Data.(map[string]interface{})
			if !ok {
				log.Println("Invalid data format")
				return
			}

			// Retrieve the type value from the map
			typeValue, ok := data["type"].(float64) // JSON numbers are float64
			if !ok {
				log.Println("Invalid type value")
				return
			}

			switch int(typeValue) {
			case currSpeed:

				var racerSpeed models.RacerCurrentWpm
				err = mapstructure.Decode(data, &racerSpeed)
				if err != nil {
					errorCh <- err
					return
				}
				currWpm, err := h.service.Multiple.CurrentSpeed(&racerSpeed)
				if err != nil {
					errorCh <- err
					close(errorCh)
					return
				}
				currSpeedCh <- currWpm
			}

		}
	}()

	go func() {
		for m := range currSpeedCh {
			writeMessage(connections, m)
		}
		for err := range errorCh {
			writeMessage(connections, map[string]interface{}{"error": err.Error()})
		}
	}()
	h.timerSender(link, connections)

}

func (h *handler) joinRacers(conn *websocket.Conn, link, id string) (*[]*websocket.Conn, models.RacerM) {
	connections, racer, err := h.service.Multiple.Join(id, conn, link)
	if err != nil {
		writeMessage(connections, map[string]interface{}{"error": err.Error()})
		log.Println(err)
		return connections, racer
	}
	writeMessage(connections, racer)
	return connections, racer
}

func (h *handler) timerSender(link string, connections *[]*websocket.Conn) chan struct{} {

	signal := make(chan struct{})

	go func() {
		for {

			t, err := h.service.Multiple.Timer(link, connections)
			if err != nil {
				if errors.Is(err, race.ErrorWaitingRacers) {
					continue
				}

				err2 := h.service.Multiple.WhiteLine(context.TODO(), link)
				if err2 != nil {
					log.Println(err2.Error())
					return
				}
				// when timer is over, and we save it in the database
				signal <- struct{}{}
				log.Println(err.Error())
				return
			}

			writeMessage(connections, map[string]interface{}{"timer": t, "type": "timer"})
		}
	}()

	return signal
}

func writeMessage(connections *[]*websocket.Conn, message interface{}) {
	if *connections == nil {
		return
	}
	sentClients := make(map[*websocket.Conn]bool)
	for _, clientConn := range *connections {
		if _, sent := sentClients[clientConn]; !sent {
			// Send the message to the writer goroutine
			err := clientConn.WriteJSON(message)
			if err != nil {
				log.Println(err)
				return
			}
			sentClients[clientConn] = true
		}
	}
}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

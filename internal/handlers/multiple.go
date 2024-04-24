package handlers

import (
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/services/multiple/race"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
)

const (
	// switch case to define a type messages I get during ws connection. Used to calculate the current speed of the racer
	currSpeed = iota + 1
	// switch case to define type messages I get during ws connection. Used to finish the race
	endRace
)

// createLink unique uuid link for the racetrack. Can lift up to 5 racer for now
//
//	@Summary		Create a racetrack
//	@Tags			multiple
//	@Description	This endpoint is used to create a racetrack. It generates a unique link for the racetrack and returns it to the user.
//	@ID				create-racetrack
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	models.LinkCreation
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Security		ApiKeyAuth
//	@Router			/track/link [post]
func (h *handler) createLink(c *gin.Context) {

	// get value from global auth middleware parsing endpoint token and if token empty user is guest
	id, ex := c.Get("ID")
	role := c.MustGet("Role")
	if !ex {
		id = role
	}

	link, err := h.service.Link.Create(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	text, err := h.service.Multiple.RandomText(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Set("text_len", len(text))

	linkResult := models.LinkCreation{
		Link:    link,
		Content: text,
	}

	c.JSON(http.StatusCreated, linkResult)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// @Summary		Join a racetrack
// @Tags			multiple
// @Description	This endpoint is used to join a racetrack. It upgrades the HTTP connection to a WebSocket connection. The server sends messages with the current race status to the client over the WebSocket connection.
// @ID				racetrack
// @Accept			json
// @Produce		json
// @Param			link	path		string	true	"Race Link"
// @Success		200		{object}	models.RacerM
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Security		ApiKeyAuth
// @Router			/track/race/{link} [get]
func (h *handler) raceTrack(c *gin.Context) {

	link := c.Param("link")

	role := c.MustGet("Role")
	id, ex := c.Get("ID")
	if !ex {
		id = role
	}

	err := h.service.Link.Check(c, link)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	connections, _ := h.joinRacers(conn, link, id.(string))

	// channels to receive messages from the client to calculate the current speed of the racer
	currSpeedCh := make(chan models.RacerSpeed)
	// channels to receive messages from the client to finish race
	endRaceResult := make(chan models.RaceResult)
	// channel to receive errors
	errorCh := make(chan error)

	go func() {
		for {
			var msgType models.IncomingMessage

			err = conn.ReadJSON(&msgType.Data)
			if err != nil {
				log.Println(err)
				errorCh <- fmt.Errorf("error during parse body try again %w", err)
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
					errorCh <- fmt.Errorf("error during parse body try again %w", err)
				}

				currWpm, err := h.service.Multiple.CurrentSpeed(&racerSpeed, c.GetInt("text_len"))
				if err != nil {
					errorCh <- fmt.Errorf("fail to calculate current speed for user %v, err=%w", racerSpeed.Email, err)
				} else {
					currSpeedCh <- currWpm
				}

			case endRace:

				var raceEnd models.RaceEndRequest

				err = mapstructure.Decode(data, &raceEnd)
				if err != nil {
					errorCh <- fmt.Errorf("error during parse body try again %w", err)
				}

				raceResult, err := h.service.Multiple.EndRace(raceEnd, link, id.(string))
				if err != nil {
					errorCh <- fmt.Errorf("fail to calculate current speed for user %v, err=%w", raceResult.Email, err)
				} else {
					endRaceResult <- raceResult
				}
				func(conn *websocket.Conn) {
					defer func(conn *websocket.Conn) {
						err = conn.Close()
						if err != nil {
							log.Printf("error closing connection %v", err)
						}
					}(conn)

					defer close(endRaceResult)
					defer close(currSpeedCh)
					defer close(errorCh)
				}(conn)
			default:
				log.Println("Invalid type value")
				errorCh <- fmt.Errorf("invalid type value %v user=%v", typeValue, id.(string))
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case v, ex := <-currSpeedCh:
				if ex {
					writeMessage(connections, v)
				}
			case v, ex := <-endRaceResult:
				if ex {
					writeMessage(connections, v)
				}
			case v, ex := <-errorCh:
				if ex {
					writeMessage(connections, map[string]interface{}{"error": v.Error()})
				}
			}
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

	// sent signal to the main goroutine when the timer is over to start race
	signal := make(chan struct{})

	go func() {
		for {

			t, err := h.service.Multiple.Timer(link, connections)
			if err != nil {
				if errors.Is(err, race.ErrorWaitingRacers) {
					continue
				}

				err2 := h.service.Multiple.WhiteLine(link)
				if err2 != nil {
					log.Println(err2.Error())
					return
				}

				// when timer is over, and we save it in the database
				signal <- struct{}{}
				log.Println(err.Error())
				return
			}

			writeMessage(connections, map[string]interface{}{"timer": t, "type": 0})
		}
	}()

	return signal
}

func writeMessage(connections *[]*websocket.Conn, message interface{}) {

	if connections == nil {
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

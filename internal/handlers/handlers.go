package handlers

import (
	"context"
	"errors"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/services/multiple/race"
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

func (h *handler) createLink(c *gin.Context) {
	//id, _ := authHeader(c)
	//
	//link, err := h.service.Link.Create(id.String())
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//text, err := h.service.Multiple.RandomText(c)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//c.Set("text_len", len(text))
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"link":    link,
	//	"content": text,
	//})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *handler) raceTrack(c *gin.Context) {

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

	connections, _ := h.joinRacers(conn, link, id.(string))

	currSpeedCh := make(chan models.RacerSpeed)
	endRaceResult := make(chan models.RaceResult)
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
				}

				currWpm, err := h.service.Multiple.CurrentSpeed(&racerSpeed, c.GetInt("text_len"))
				if err != nil {
					errorCh <- err
				} else {
					currSpeedCh <- currWpm
				}

			case endRace:

				var raceEnd models.RaceEndRequest

				err = mapstructure.Decode(data, &raceEnd)
				if err != nil {
					errorCh <- err
				}

				raceResult, err := h.service.Multiple.EndRace(raceEnd, link, id.(string))
				if err != nil {
					errorCh <- err
				} else {
					endRaceResult <- raceResult
				}

			default:
				log.Println("Invalid type value")
				writeMessage(connections, map[string]interface{}{"error": "Invalid type value"})
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
			//for m := range currSpeedCh {
			//	writeMessage(connections, m)
			//}
			//for err := range errorCh {
			//	writeMessage(connections, map[string]interface{}{"error": err.Error()})
			//}
			//for r := range endRaceResult {
			//	writeMessage(connections, r)
			//}
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

package handlers

import (
	"context"
	"encoding/json"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/MamushevArup/typeracer/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func authHeader(c *gin.Context) (uuid.UUID, string) {
	auth := c.GetHeader("Authorization")
	// as usual cut bearer prefix
	token := strings.Split(auth, " ")[0]
	validateToken, err := utils.ValidateToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return [16]byte{}, ""
	}

	racerUUID, err := uuid.Parse(validateToken.ID)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return [16]byte{}, ""
	}
	return racerUUID, validateToken.Role
}

func (h *handler) startRace(c *gin.Context) {

	racerUUID, _ := authHeader(c)

	ex, err := h.service.PracticeY.RacerExists(context.Background(), racerUUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}

	race, err := h.service.PracticeY.StartRace(context.Background(), racerUUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var stRace struct {
		ID              uuid.UUID `json:"-" db:"id"`
		RacerID         uuid.UUID `json:"racer_id"`
		TextID          uuid.UUID `json:"-"`
		Text            string    `json:"content" db:"content"`
		TextLen         int       `json:"length" db:"length"`
		TextAuthor      string    `json:"text_author" db:"author"`
		ContributorName string    `json:"contributor_name" db:"contributor"`
		RacerName       string    `json:"racer_name" db:"username"`
		Avatar          string    `json:"avatar" db:"avatar"`
	}
	marshal, err := json.Marshal(race)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(marshal, &stRace)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, &stRace)
}

func (h *handler) currWPM(c *gin.Context) {
	/*
		{
			"current index", "duration"
		}
	*/

	racerUUID, _ := authHeader(c)

	var curr struct {
		CurrIdx  int `json:"current_index"`
		Duration int `json:"duration"`
	}
	var resp struct {
		Wpm int `json:"wpm"`
	}
	if err := c.BindJSON(&curr); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
	}

	// Perform other validations as needed
	if curr.Duration < 0 {
		newErrorResponse(c, http.StatusBadRequest, "duration must be non-negative")
		return
	}

	if curr.CurrIdx <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "index must be non-negative")
		return
	}

	ex, err := h.service.PracticeY.RacerExists(context.Background(), racerUUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}
	calc, err := h.service.PracticeY.RealTimeCalc(context.Background(), curr.CurrIdx, curr.Duration)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.Wpm = calc
	c.JSON(http.StatusCreated, &resp)
}

func (h *handler) endRace(c *gin.Context) {

	racerUUID, _ := authHeader(c)

	var req *models.ReqEndSingle

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid message body")
		return
	}

	// Perform other validations as needed
	if req.Duration < 0 {
		newErrorResponse(c, http.StatusBadRequest, "duration must be non-negative")
		return
	}

	if req.Errors < 0 {
		newErrorResponse(c, http.StatusBadRequest, "errors must be non-negative")
		return
	}

	ex, err := h.service.PracticeY.RacerExists(context.Background(), racerUUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}
	req.RacerId = racerUUID
	race, err := h.service.PracticeY.EndRace(c, req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var endRace struct {
		Wpm      int     `json:"wpm" db:"speed"`
		Accuracy float64 `json:"accuracy" db:"accuracy"`
		Duration int     `json:"duration" db:"duration"`
	}

	end, err := json.Marshal(race)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(end, &endRace)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, &endRace)
}

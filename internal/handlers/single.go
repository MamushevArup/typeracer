package handlers

import (
	"context"
	"encoding/json"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// TODO implement swagger docs
func (h *handler) startRace(c *gin.Context) {

	id, ex := c.Get("ID")
	role := c.MustGet("Role")

	if !ex {
		id = role
	}

	ex, err := h.service.Single.RacerExists(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}

	race, err := h.service.Single.StartRace(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, race)
}

// TODO implement swagger docs
func (h *handler) currWPM(c *gin.Context) {

	id, ex := c.Get("ID")
	role := c.MustGet("Role")

	if !ex {
		id = role
	}

	var curr struct {
		CurrIdx  int `json:"index"`
		Duration int `json:"duration"`
	}
	var resp struct {
		Wpm int `json:"wpm"`
	}

	if err := c.BindJSON(&curr); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// Perform other validations as needed
	if curr.Duration < 0 {
		newErrorResponse(c, http.StatusBadRequest, "duration must be non-negative")
		return
	}

	if curr.CurrIdx <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "index must be non-negative or not zero")
		return
	}

	ex, err := h.service.Single.RacerExists(context.Background(), id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}

	calc, err := h.service.Single.RealTimeCalc(context.Background(), curr.CurrIdx, curr.Duration)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	resp.Wpm = calc
	c.JSON(http.StatusCreated, &resp)
}

func (h *handler) endRace(c *gin.Context) {

	//racerUUID, _ := authHeader(c)

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

	ex, err := h.service.Single.RacerExists(context.Background(), uuid.Nil.String())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}
	req.RacerId = uuid.Nil
	race, err := h.service.Single.EndRace(c, req)
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

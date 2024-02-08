package handlers

import (
	"context"
	"encoding/json"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *handler) startRace(c *gin.Context) {
	// TODO check if racer with this id exist and correctness of type race
	/*
		{
			"racer_id" : uuid.UUID,
		}
	*/
	var racer struct {
		RacerID uuid.UUID `json:"racer_id"`
	}

	if err := c.BindJSON(&racer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if _, err := uuid.Parse(racer.RacerID.String()); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid UUID format for racer_id")
		return
	}

	ex, err := h.service.PracticeY.RacerExists(context.Background(), racer.RacerID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}

	race, err := h.service.PracticeY.StartRace(context.Background(), racer.RacerID)
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
			"racer_id", "current index", "duration"
		}
	*/
	var curr struct {
		RacerId  uuid.UUID `json:"racer_id"`
		CurrIdx  int       `json:"current_index"`
		Duration int       `json:"duration"`
	}
	var resp struct {
		RacerId uuid.UUID `json:"racer_id"`
		Wpm     int       `json:"wpm"`
	}
	if err := c.BindJSON(&curr); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
	}

	// Check if the RacerID is a valid UUID
	if _, err := uuid.Parse(curr.RacerId.String()); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid UUID format for racer_id")
		return
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

	ex, err := h.service.PracticeY.RacerExists(context.Background(), curr.RacerId)
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
	resp.RacerId = curr.RacerId
	resp.Wpm = calc
	c.JSON(http.StatusCreated, &resp)
}

func (h *handler) endRace(c *gin.Context) {
	var req *models.ReqEndSingle

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid message body")
		return
	}

	// Check if the RacerID is a valid UUID
	if _, err := uuid.Parse(req.RacerID.String()); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid UUID format for racer_id")
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

	ex, err := h.service.PracticeY.RacerExists(context.Background(), req.RacerID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}
	race, err := h.service.PracticeY.EndRace(c, req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var endRace struct {
		RacerId  uuid.UUID `json:"racer_id"`
		Wpm      int       `json:"wpm" db:"speed"`
		Accuracy float64   `json:"accuracy" db:"accuracy"`
		Duration int       `json:"duration" db:"duration"`
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

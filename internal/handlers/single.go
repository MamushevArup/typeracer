package handlers

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) startRace(c *gin.Context) {
	// TODO check if racer with this id exist and correctness of type race
	/*
		{
			"racer_id" : uuid.UUID,
		}
	*/
	var single models.Single

	if err := c.BindJSON(&single); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	race, err := h.service.PracticeY.StartRace(context.Background(), single.RacerID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, &race)
}

func (h *handler) endRace(c *gin.Context) {
	var req models.ReqEndSingle

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid message body")
		return
	}
	race, err := h.service.PracticeY.EndRace(c, req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, &race)
}

package handlers

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *handler) moderation(c *gin.Context) {
}

func (h *handler) contribute(c *gin.Context) {
	var ctr models.ContributeText
	entryValidation(c, h, &ctr)
	err := h.service.Contribute.ContributeText(ctr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "text successfully sent to moderation",
	})
}

func entryValidation(c *gin.Context, h *handler, ctr *models.ContributeText) {

	var (
		minLenContent = 20
		maxLenContent = 300
	)

	if err := c.BindJSON(&ctr); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	if _, err := uuid.Parse(ctr.RacerID.String()); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid UUID format for racer_id")
		return
	}
	ex, err := h.service.Contribute.RacerExists(context.Background(), ctr.RacerID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}
	if ctr.Content == "" || ctr.Author == "" || ctr.Source == "" || ctr.SourceTitle == "" {
		newErrorResponse(c, http.StatusBadRequest, "must fill all entries")
		return
	}
	if len(ctr.Content) < minLenContent || len(ctr.Content) > maxLenContent {
		newErrorResponse(c, http.StatusBadRequest, "length should be in range 20 to 300")
		return
	}
}

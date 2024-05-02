package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var ErrRacerNotFound = "racer not found"

func (h *handler) profileInfo(c *gin.Context) {
	id := c.MustGet("ID")

	ex, err := h.service.Single.RacerExists(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, ErrRacerNotFound)
		return
	}

	racerInfo, err := h.service.Racer.Details(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, racerInfo)
}

func (h *handler) updateProfile(c *gin.Context) {
	id := c.MustGet("ID")

	var updateRacer models.RacerUpdate
	if err := c.BindJSON(&updateRacer); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	updateRacer.Id = id.(string)

	_, err := govalidator.ValidateStruct(updateRacer)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Racer.UpdateRacerInfo(c, updateRacer)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handler) singleHistory(c *gin.Context) {
	id := c.MustGet("ID")

	limit := c.Query("limit")
	offset := c.Query("offset")

	sglHistory, err := h.service.Racer.SingleHistory(c.Request.Context(), id.(string), limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, sglHistory)
}

func (h *handler) singleHistoryText(c *gin.Context) {
	id := c.MustGet("ID")
	sId := c.Param("single_id")

	sUUID, err := uidVerifyAndConvert(sId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	sglHistoryText, err := h.service.Racer.HistorySingleText(c.Request.Context(), id.(string), sUUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, sglHistoryText)
}

func uidVerifyAndConvert(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

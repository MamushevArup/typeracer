package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
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

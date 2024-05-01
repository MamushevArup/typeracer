package handlers

import (
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

func (h *handler) avatars(c *gin.Context) {
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

	avatars, err := h.service.Racer.Avatars(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, avatars)
}

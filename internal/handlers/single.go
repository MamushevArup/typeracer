package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @Summary		Start a new race
// @Tags			single
// @Description	This endpoint is used to start a new race for a racer.
// @ID				start-race
// @Accept			json
// @Produce		json
// @Success		201	{object}	models.SingleResponse
// @Failure		400	{object}	errorResponse
// @Failure		404	{object}	errorResponse
// @Failure		500	{object}	errorResponse
// @Security		Bearer
// @Router			/single/race [get]
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

// @Summary		Calculate current Words Per Minute (WPM)
// @Tags			single
// @Description	This endpoint is used to calculate the current WPM for a racer.
// @ID				curr-wpm
// @Accept			json
// @Produce		json
// @Param			models.CountWpm	body		models.CountWpm	true	"Wpm calculation"
// @Success		201				{object}	models.Speed
// @Failure		400				{object}	errorResponse
// @Failure		404				{object}	errorResponse
// @Failure		500				{object}	errorResponse
// @Security		Bearer
// @Router			/single/curr-wpm [post]
func (h *handler) currWPM(c *gin.Context) {

	id, ex := c.Get("ID")
	role := c.MustGet("Role")

	if !ex {
		id = role
	}

	var wpmCounter models.CountWpm

	if err := c.BindJSON(&wpmCounter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// Perform other validations as needed
	if wpmCounter.Duration < 0 {
		newErrorResponse(c, http.StatusBadRequest, "duration must be non-negative")
		return
	}

	if wpmCounter.CurrIdx <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "index must be non-negative or not zero")
		return
	}

	ex, err := h.service.Single.RacerExists(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}

	calc, err := h.service.Single.RealTimeCalc(c, wpmCounter.CurrIdx, wpmCounter.Duration)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var count models.Speed

	count.Wpm = calc

	c.JSON(http.StatusCreated, count)
}

// @Summary		End a race
// @Tags			single
// @Description	This endpoint is used to end a race for a racer.
// @ID				end-race
// @Accept			json
// @Produce		json
// @Param			models.ReqEndSingle	body		models.ReqEndSingle	true	"End Race"
// @Success		201					{object}	models.RespEndSingle
// @Failure		400					{object}	errorResponse
// @Failure		404					{object}	errorResponse
// @Failure		500					{object}	errorResponse
// @Security		Bearer
// @Router			/single/end-race [post]
func (h *handler) endRace(c *gin.Context) {

	var req models.ReqEndSingle

	id, ex := c.Get("ID")
	role := c.MustGet("Role")

	if !ex {
		id = role
	}

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

	ex, err := h.service.Single.RacerExists(c, id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !ex {
		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
		return
	}

	uid, err := uuid.Parse(id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	req.RacerId = uid

	race, err := h.service.Single.EndRace(c, req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, &race)
}

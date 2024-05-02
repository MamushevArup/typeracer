package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var ErrRacerNotFound = "racer not found"

// @Summary Get profile info
// @Tags profile
// @Description get profile info
// @ID get-profile-info
// @Produce json
// @Security Bearer
// @Success 200 {object} models.RacerHandler
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile/info [get]
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

// @Summary update profile
// @Tags profile
// @Description update profile
// @ID update-profile
// @Accept json
// @Produce json
// @Param models.RacerUpdate body models.RacerUpdate true "racer update info"
// @Security Bearer
// @Success 200
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile/update [put]
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

// @Summary single history info details
// @Tags history
// @Description single history info
// @ID single-history
// @Produce json
// @Security Bearer
// @Success 200 {array} models.SingleHistoryHandler
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile/history/single/ [get]
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

// @Summary single history text details
// @Tags history
// @Description single history text details
// @ID single-history-text
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SingleHistoryText
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile/history/single/{single_id} [get]
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

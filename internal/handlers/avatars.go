package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Select avatar
// @Tags profile
// @Description select avatar
// @ID select-avatar
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Avatar
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile/avatars [get]
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

// @Summary update avatar
// @Tags profile
// @Description update avatar
// @ID update avatar
// @Accept json
// @Produce json
// @Param models.AvatarUpdate body models.AvatarUpdate true "avatar update info"
// @Security Bearer
// @Success 200
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /profile/avatar [put]
func (h *handler) updateAvatar(c *gin.Context) {
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

	var avatarInfo models.AvatarUpdate
	if err = c.BindJSON(&avatarInfo); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	avatarInfo.RacerId = id.(string)

	err = h.service.Racer.UpdateAvatar(c.Request.Context(), avatarInfo)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

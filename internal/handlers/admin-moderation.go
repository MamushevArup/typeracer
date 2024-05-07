package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary		Get all text to see in pending mode
// @Tags		moderation
// @Description	This endpoint is used for see content in moderation
// @ID				moderation-all
// @Accept			json
// @Produce		json
// @Param 			limit query string false "limit"
// @Param 			offset query string false "offset"
// @Param 			sort query string false "sort"
// @Success		201				{object}	[]models.ModerationServiceResponse
// @Failure		400				{object}	errorResponse
// @Failure		500				{object}	errorResponse
//
//	@Security		Bearer
//
// @Router			/admin/moderation/all [get]
func (h *handler) showContentToModerate(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")

	sort := c.Query("sort")

	allEntries, err := h.service.Admin.ShowContentToModerate(c.Request.Context(), limit, offset, sort)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, allEntries)
}

// @Summary Moderation details
// @Description Get details of a specific moderation item
// @ID moderation-content
// @Tags moderation
// @Accept  json
// @Produce  json
// @Param moderation_id path string true "Moderation ID"
// @Success 200 {object} models.ModerationTextDetails
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security Bearer
// @Router /admin/moderation/{moderation_id} [get]
func (h *handler) moderationText(c *gin.Context) {
	modId := c.Param("moderation_id")

	details, err := h.service.Admin.TextDetails(c.Request.Context(), modId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, details)
}

// @Summary Approve provided text content
// @Description Admin can approve this content and it appear in global text storage
// @ID moderation-approve
// @Tags moderation
// @Accept  json
// @Produce  json
// @Param moderation_id path string true "Moderation ID"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security Bearer
// @Router /admin/moderation/content/{moderation_id}/approve [POST]
func (h *handler) approveContent(c *gin.Context) {
	modId := c.Param("moderation_id")

	err := h.service.Admin.ApproveContent(c.Request.Context(), modId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Reject provided content
// @Description Admin can reject some content because of problem in content
// @ID moderation-reject
// @Tags moderation
// @Accept  json
// @Produce  json
// @Param moderation_id path string true "Moderation ID"
// @Param models.ModerationRejectToService body models.ModerationRejectToService true "Reject"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security Bearer
// @Router /admin/moderation/content/{moderation_id}/reject [POST]
func (h *handler) rejectContent(c *gin.Context) {
	modId := c.Param("moderation_id")

	var reject models.ModerationRejectToService

	if err := c.ShouldBindJSON(&reject); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	reject.ModerationID = modId

	err := h.service.Admin.RejectContent(c.Request.Context(), reject)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

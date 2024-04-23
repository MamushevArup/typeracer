package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (h *handler) moderationText(c *gin.Context) {
	modId := c.Param("moderation_id")

	details, err := h.service.Admin.TextDetails(c.Request.Context(), modId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, details)
	return
}

func (h *handler) approveContent(c *gin.Context) {
	modId := c.Param("moderation_id")

	err := h.service.Admin.ApproveContent(c.Request.Context(), modId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
	return
}

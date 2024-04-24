package handlers

import (
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	maxContentLength  = 700
	minContentLength  = 100
	authorMaxLen      = 20
	sourceMaxLen      = 15
	sourceTitleMaxLen = 30
)

// @Summary		Contribute text
// @Tags			content
// @Description	Endpoint related to contribute text to the general text set
// @ID				contribute
// @Accept			json
// @Produce		json
// @Param			models.ContributeHandlerRequest	body		models.ContributeHandlerRequest	true	"Contribute"
// @Success		201
// @Failure		400
// @Failure		500
// @Security		Bearer
// @Router			/content/contribute [post]
func (h *handler) contribute(c *gin.Context) {
	var review models.ContributeHandlerRequest

	if err := c.ShouldBindJSON(&review); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validModerationInput(review); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, ex := c.Get("ID")
	if !ex {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	review.RacerID = id.(string)

	if err := h.service.Contribute.ContributeText(c, review); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func validModerationInput(r models.ContributeHandlerRequest) error {
	if isEqual(r.Author, "") || isEqual(r.Content, "") || isEqual(r.Source, "") || isEqual(r.SourceTitle, "") {
		return fmt.Errorf("empty fields are not allowed")
	}

	if len(r.Content) < minContentLength || len(r.Content) > maxContentLength {
		return fmt.Errorf("content length should be between %d and %d", minContentLength, maxContentLength)
	}

	if len(r.Author) > authorMaxLen {
		return fmt.Errorf("author length should be less than %d", authorMaxLen)
	}

	if len(r.Source) > sourceMaxLen {
		return fmt.Errorf("source length should be less than %d", sourceMaxLen)
	}

	if len(r.SourceTitle) > sourceTitleMaxLen {
		return fmt.Errorf("source title length should be less than %d", sourceTitleMaxLen)
	}

	return nil
}

func isEqual(str, str2 string) bool {
	return str == str2
}

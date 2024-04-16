package handlers

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/MamushevArup/typeracer/internal/models"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"net/http"
//	"strings"
//)
//
//func (h *handler) moderation(c *gin.Context) {
//}
//
//func (h *handler) contribute(c *gin.Context) {
//	var ctr models.ContributeText
//	code, err := entryValidation(c, h, &ctr)
//	if err != nil {
//		newErrorResponse(c, code, err.Error())
//		return
//	}
//	err = h.service.Contribute.ContributeText(ctr)
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return
//	}
//	c.JSON(http.StatusCreated, gin.H{
//		"message": "text successfully sent to moderation",
//	})
//}
//
//func entryValidation(c *gin.Context, h *handler, ctr *models.ContributeText) (int, error) {
//
//	var (
//		minLenContent = 20
//		maxLenContent = 300
//	)
//
//	if err := c.BindJSON(&ctr); err != nil {
//		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
//		return http.StatusBadRequest, err
//	}
//	accessT := c.GetHeader("Authorization")
//	slice := strings.Split(accessT, " ")
//	if len(slice) != 2 {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "auth header is not well formatted expect Bearer <token>"})
//		return http.StatusBadRequest, errors.New("auth header incorrect")
//	}
//	token := slice[1]
//	access, err := utils.ValidateToken(token)
//	if err != nil {
//		newErrorResponse(c, http.StatusForbidden, err.Error())
//		return http.StatusForbidden, err
//	}
//	uAcess, err := uuid.Parse(access.ID)
//	if err != nil {
//		newErrorResponse(c, http.StatusBadRequest, "invalid UUID format for racer_id")
//		return http.StatusBadRequest, err
//	}
//	ctr.RacerID = uAcess
//	fmt.Println(ctr.RacerID, "RACERID")
//	ex, err := h.service.Contribute.RacerExists(context.Background(), ctr.RacerID)
//	if err != nil {
//		newErrorResponse(c, http.StatusInternalServerError, err.Error())
//		return http.StatusInternalServerError, err
//	}
//	if !ex {
//		newErrorResponse(c, http.StatusNotFound, "racer doesn't exist")
//		return http.StatusNotFound, errors.New("racer doesn't exist")
//	}
//	if ctr.Content == "" || ctr.Author == "" || ctr.Source == "" || ctr.SourceTitle == "" {
//		newErrorResponse(c, http.StatusBadRequest, "must fill all entries")
//		return http.StatusBadRequest, errors.New("entries not filled")
//	}
//	if len(ctr.Content) < minLenContent || len(ctr.Content) > maxLenContent {
//		newErrorResponse(c, http.StatusBadRequest, "content for race is small at 20 symbols")
//		return http.StatusBadRequest, errors.New("length not acceptable")
//	}
//	return http.StatusOK, nil
//}

package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (h *handler) adminSignIn(c *gin.Context) {
	var sign models.AdminSignIn

	if err := c.BindJSON(&sign); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if sign.Username != adminUsername || sign.Password != adminPassword {
		newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	accessT, refresh, err := h.service.Auth.SignIn(c, sign.Username, sign.Password, sign.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(cookieName, refresh, maxAge, "/admin", domain, false, true)

	c.JSON(http.StatusOK, models.AdminSignInResponse{Access: accessT})

}

func (h *handler) adminRefresh(c *gin.Context) {

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cookie not sent")
		return
	}

	access, refresh, err := h.service.Auth.AdminRefresh(c, cookie)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(cookieName, refresh, maxAge, path, domain, false, true)

	c.JSON(http.StatusCreated, models.AdminSignInResponse{Access: access})
}

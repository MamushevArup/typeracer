package handlers

import (
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary		Sign in for admins
// @Tags		admin
// @Description	This endpoint is used for user authentication.
// @ID				sign-in-admin
// @Accept			json
// @Produce		json
// @Param			models.AdminSignIn	body		models.AdminSignIn	true	"Sign In"
// @Success		201				{object}	models.SignInHandler
// @Failure		400				{object}	errorResponse
// @Failure		500				{object}	errorResponse
//
//	@Security		ApiKeyAuth
//
// @Router			/admin/auth/sign-in [post]
func (h *handler) adminSignIn(c *gin.Context) {
	var sign models.AdminSignIn

	if err := c.BindJSON(&sign); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	adminUsername := h.cfg.Admin.Username
	adminPassword := h.cfg.Admin.Password

	if sign.Username != adminUsername || sign.Password != adminPassword {
		newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	adminInfo, err := h.service.Auth.SignIn(c, sign.Username, sign.Password, sign.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ad := models.SignInHandler{
		Access:   adminInfo.Access,
		Username: adminInfo.Username,
		Avatar:   "",
	}

	c.SetCookie(cookieName, adminInfo.Refresh, maxAge, "/admin", domain, false, true)

	c.JSON(http.StatusOK, ad)

}

// @Summary		refresh for admins
// @Tags		admin
// @Description	Admin can refresh their token
// @ID				refresh_admin
// @Accept			json
// @Produce		json
// @Success		201				{object}	models.AdminSignInResponse
// @Failure		400				{object}	errorResponse
// @Failure		500				{object}	errorResponse
// @Router			/admin/auth/refresh [post]
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

	c.SetCookie(cookieName, refresh, maxAge, "/admin", domain, false, true)

	c.JSON(http.StatusCreated, models.AdminSignInResponse{Access: access})
}

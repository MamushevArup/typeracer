package handlers

import (
	"errors"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
)

var (
	emptyUsername = errors.New("username cannot be empty")
	emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	maxAge        = int(time.Now().Add(60 * 24 * time.Hour).Unix())
)

const (
	path       = "/api/auth"
	domain     = "localhost"
	cookieName = "refresh_token"
)

// @Summary Sign in
// @Tags auth
// @Description This endpoint is used for user authentication.
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param   models.SignIn     body    models.SignIn     true        "Sign In"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router  /api/auth/sign-in [post]
func (h *handler) signIn(c *gin.Context) {
	var sign models.SignIn

	if err := c.BindJSON(&sign); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if sign.Fingerprint == "" {
		newErrorResponse(c, http.StatusBadRequest, "fingerprint empty")
		return
	}

	if err := emailAndPasswdValidation(sign.Email, sign.Password); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	access, refresh, err := h.service.Auth.SignIn(c, sign.Email, sign.Password, sign.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var a models.AuthResponse
	a.Access = access

	c.SetCookie(cookieName, refresh, maxAge, path, domain, false, true)

	c.JSON(http.StatusOK, a)
}

// @Summary Sign up
// @Tags auth
// @Description This endpoint is used for user registration.
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param   models.SignUp     body    models.SignUp     true        "Sign Up"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router  /api/auth/sign-up [post]
func (h *handler) signUp(c *gin.Context) {

	var s models.SignUp

	if err := c.BindJSON(&s); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
	}

	err := emailAndPasswdValidation(s.Email, s.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if s.Username == "" {
		newErrorResponse(c, http.StatusBadRequest, emptyUsername.Error())
		return
	}

	err = h.service.Auth.CheckUserSignUp(c, s.Email, s.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	access, refresh, err := h.service.Auth.SignUp(c, s.Email, s.Username, s.Password, s.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var a models.AuthResponse
	a.Access = access

	c.SetCookie(cookieName, refresh, maxAge, path, domain, false, true)

	c.JSON(http.StatusCreated, a)
}

func emailAndPasswdValidation(email, password string) error {
	var (
		minLen = 6
		maxLen = 15
		digit  = `\d`
		char   = `[A-Za-z]`
		symbol = `[\x21-\x2F\x3A-\x40\x5B-\x60\x7B-\x7E]`
	)
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	// Check password length
	if len(password) < minLen || len(password) > maxLen {
		return fmt.Errorf("password length must be between %d and %d characters", minLen, maxLen)
	}

	hasDigit := regexp.MustCompile(digit).MatchString(password)
	hasChar := regexp.MustCompile(char).MatchString(password)
	hasSymbol := regexp.MustCompile(symbol).MatchString(password)

	if !hasDigit || !hasChar || !hasSymbol {
		return fmt.Errorf("password must contain at least one character, one digit, and one symbol")
	}

	return nil
}

// @Summary Log out
// @Tags auth
// @Description This endpoint is used for user logout.
// @ID log-out
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security Bearer
// @Router  /api/auth/logout [delete]
func (h *handler) logOut(c *gin.Context) {

	// here I will get refresh token and remove by this refresh token
	refresh, err := c.Cookie(cookieName)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	err = h.service.Auth.Logout(c, refresh)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(200)
}

// @Summary Refresh token
// @Tags auth
// @Description This endpoint is used to refresh the access token.
// @ID refresh
// @Accept  json
// @Produce  json
// @Param models.RefreshS     body    models.RefreshS     true        "Refresh"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security Bearer
// @Router  /api/auth/refresh [post]
func (h *handler) refresh(c *gin.Context) {

	var r models.RefreshS

	if err := c.BindJSON(&r); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cookie not sent")
		return
	}

	access, refresh, err := h.service.Auth.RefreshToken(c, cookie, r.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var a models.AuthResponse
	a.Access = access

	c.SetCookie(cookieName, refresh, maxAge, path, domain, false, true)

	c.JSON(http.StatusCreated, a)
}

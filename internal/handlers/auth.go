package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
)

var (
	mail          = errors.New("invalid email format")
	emptyUsername = errors.New("username cannot be empty")
	emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	maxAge        = int(time.Now().Add(60 * 24 * time.Hour).Unix())
)

type authResponse struct {
	Access string `json:"access"`
}
type signIn struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type signUp struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type refreshS struct {
	Fingerprint string `json:"fingerprint"`
}

// @Summary Sign in
// @Tags auth
// @Description This endpoint is used for user authentication.
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param   signIn     body    signIn     true        "Sign In"
// @Success 201 {object} authResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router  /api/auth/sign-in [post]
func (h *handler) signIn(c *gin.Context) {
	var sign signIn

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

	access, refresh, err := h.service.Auth.SignIn(context.TODO(), sign.Email, sign.Password, sign.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var a authResponse
	a.Access = access

	c.SetCookie("refresh_token", refresh, maxAge, "/api/auth", "localhost", false, true)

	c.JSON(http.StatusOK, a)
}

// @Summary Sign up
// @Tags auth
// @Description This endpoint is used for user registration.
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param   signUp     body    signUp     true        "Sign Up"
// @Success 201 {object} authResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router  /api/auth/sign-up [post]
func (h *handler) signUp(c *gin.Context) {

	var s signUp

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
	err = h.service.Auth.CheckUserSignUp(context.TODO(), s.Email, s.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadGateway, err.Error())
		return
	}

	access, refresh, err := h.service.Auth.SignUp(context.TODO(), s.Email, s.Username, s.Password, s.Fingerprint)

	if err != nil {
		return
	}

	var a authResponse
	a.Access = access

	c.SetCookie("refresh_token", refresh, maxAge, "/api/auth", "localhost", false, true)
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
		return mail
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
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	err = h.service.Auth.Logout(context.TODO(), refresh)
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
// @Param refreshS     body    refreshS     true        "Refresh"
// @Success 201 {object} authResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security Bearer
// @Router  /api/auth/refresh [post]
func (h *handler) refresh(c *gin.Context) {
	var r refreshS
	if err := c.BindJSON(&r); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cookie not sent")
		return
	}
	access, refresh, err := h.service.Auth.RefreshToken(context.TODO(), cookie, r.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var a authResponse
	a.Access = access

	c.SetCookie("refresh_token", refresh, maxAge, "/api/auth", "localhost", false, true)
	c.JSON(http.StatusCreated, a)
}

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
	symbol        = errors.New("password must contain at least one special symbol, one character and one digit")
	mail          = errors.New("invalid email format")
	emptyUsername = errors.New("username cannot be empty")
	emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func (h *handler) signUp(c *gin.Context) {
	/*
		{
			"email" : "arupmamushev@gmail.com",
			"username"  : "Cicada_3301",
			"Password" : "Hello world"
		}
	*/
	var r struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&r); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
	}

	err := emailAndPasswdValidation(r.Email, r.Password, r.Username)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.Auth.CheckUserSignUp(context.TODO(), r.Email, r.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadGateway, err.Error())
		return
	}

	access, refresh, err := h.service.Auth.SignUp(context.TODO(), r.Email, r.Username, r.Password)

	if err != nil {
		return
	}
	var res struct {
		Access string `json:"access_token"`
	}
	res.Access = access
	maxAge := int(time.Now().Add(60 * 24 * time.Hour).Unix())
	c.SetCookie("refresh_token", refresh, maxAge, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, &res)
}

func emailAndPasswdValidation(email, password, username string) error {
	var (
		minLen    = 6
		maxLen    = 15
		hasDigit  bool
		hasChar   bool
		hasSymbol bool
	)
	if username == "" {
		return emptyUsername
	}
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		return mail
	}

	// Check password length
	if len(password) < minLen || len(password) > maxLen {
		return fmt.Errorf("password length must be between %d and %d characters", minLen, maxLen)
	}

	// Check password composition
	for _, char := range password {
		switch {
		case '0' <= char && char <= '9':
			hasDigit = true
		case 'a' <= char && char <= 'z', 'A' <= char && char <= 'Z':
			hasChar = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&':
			hasSymbol = true
		}
	}
	if !hasDigit || !hasChar || !hasSymbol {
		return symbol
	}
	return nil
}

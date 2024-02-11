package handlers

import (
	"context"
	"github.com/MamushevArup/typeracer/internal/middleware"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// all validation will make here

type Handler interface {
	InitRoutes() *gin.Engine
}

type handler struct {
	service *services.Service
}

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// this middleware check for jwt token valid and expiry
	mdl := middleware.NewMiddleware(h.service)

	router.Use(mdl.AuthMiddleware())

	sgl := router.Group("/single")
	{
		sgl.POST("/race", h.startRace)
		sgl.POST("/end-race", h.endRace)
		sgl.POST("/curr-wpm", h.currWPM)
	}
	router.POST("/contribute", h.contribute)
	router.GET("/moderation/:id", h.moderation)

	auth := router.Group("/api/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/logout", h.logOut)
		auth.POST("/refresh", h.refresh)
	}

	return router
}

func (h *handler) logOut(context *gin.Context) {

}

func (h *handler) refresh(c *gin.Context) {
	var f struct {
		Fingerprint string `json:"fingerprint"`
	}
	if err := c.BindJSON(&f); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cookie not sent")
		return
	}
	access, refresh, err := h.service.Auth.RefreshToken(context.TODO(), cookie, f.Fingerprint)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.SetCookie("refresh_token", refresh, maxAge, "/api/auth", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"access": access,
	})
}

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

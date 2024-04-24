package handlers

import (
	"github.com/MamushevArup/typeracer/internal/config"
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-gonic/gin"
)

func NewHandler(service *services.Service, cfg *config.Config) Handler {
	return &handler{service: service, cfg: cfg}
}

type Handler interface {
	InitRoutes() *gin.Engine
}

type handler struct {
	service *services.Service
	cfg     *config.Config
}

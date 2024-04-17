package handlers

import (
	"github.com/MamushevArup/typeracer/internal/services"
	"github.com/gin-gonic/gin"
)

func NewHandler(service *services.Service) Handler {
	return &handler{service: service}
}

type Handler interface {
	InitRoutes() *gin.Engine
}

type handler struct {
	service *services.Service
}

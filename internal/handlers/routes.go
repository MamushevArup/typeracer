package handlers

import (
	_ "github.com/MamushevArup/typeracer/docs"
	"github.com/MamushevArup/typeracer/internal/middleware/auth/http/access"
	"github.com/MamushevArup/typeracer/internal/middleware/auth/http/admin"
	"github.com/MamushevArup/typeracer/internal/middleware/auth/http/endpoint"
	validation "github.com/MamushevArup/typeracer/internal/middleware/auth/http/token-validation"
	validationWs "github.com/MamushevArup/typeracer/internal/middleware/auth/ws/token-ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Use(validation.TokenInspector())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	contribute := router.Group("/content")
	{
		contribute.POST("/contribute", endpoint.Access(h.service), h.contribute)
	}

	moder := router.Group("/admin")
	{
		moderAuth := moder.Group("/auth")
		{
			moderAuth.POST("/sign-in", h.adminSignIn)
			moderAuth.POST("/refresh", h.adminRefresh)
		}
		moderation := moder.Group("/moderation")
		moderation.Use(admin.Confirm())
		{
			moderation.GET("/all", h.showContentToModerate)
			moderation.GET("/:moderation_id", h.moderationText)
			moderation.POST("/content/:moderation_id/approve", h.approveContent)
		}
	}

	sgl := router.Group("/single")
	sgl.Use(access.OnlyGuestOrRacer())
	{
		sgl.GET("/race", h.startRace)
		sgl.POST("/end-race", h.endRace)
		sgl.POST("/curr-wpm", h.currWPM)
	}

	//router.POST("/contribute", h.contribute)
	//router.GET("/moderation/:id", h.moderation)

	auth := router.Group("/api/auth")

	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.DELETE("/logout", h.logOut)
		auth.POST("/refresh", h.refresh)
	}
	// this route stands for create racetrack and start a multiple race
	mlt := router.Group("/track")
	mlt.Use(validationWs.TokenVerifier())
	mlt.Use(access.OnlyGuestOrRacer())
	{
		mlt.POST("/link", h.createLink)
		// this racetrack will look like this. /race/link?t=<endpoint token>
		mlt.GET("/race/:link", h.raceTrack)
	}

	return router
}

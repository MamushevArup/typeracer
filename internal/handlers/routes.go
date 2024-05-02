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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin", "X-Requested-With", "X-Access-Token", "X-Refresh-Token"},
		AllowCredentials: true,
		MaxAge:           0,
		AllowWildcard:    true,
		AllowWebSockets:  true,
		AllowFiles:       true,
	}))

	router.Use(validation.TokenInspector(h.cfg))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// documented with swagger now
	contribute := router.Group("/content")
	{
		contribute.POST("/contribute", endpoint.Access(h.service), h.contribute)
	}

	// documented with swagger now
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

			content := moderation.Group("/content")
			{
				content.POST("/:moderation_id/approve", h.approveContent)
				content.POST("/:moderation_id/reject", h.rejectContent)
			}

		}
		moder.POST("/add-cars", h.addCars)
	}

	sgl := router.Group("/single")
	sgl.Use(access.OnlyGuestOrRacer())
	{
		sgl.GET("/race", h.startRace)
		sgl.POST("/end-race", h.endRace)
		sgl.POST("/curr-wpm", h.currWPM)
	}

	auth := router.Group("/api/auth")

	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.DELETE("/logout", h.logOut)
		auth.POST("/refresh", h.refresh)
	}
	// this route stands for create racetrack and start a multiple race
	mlt := router.Group("/track")
	mlt.Use(access.OnlyGuestOrRacer())
	{
		mlt.POST("/link", h.createLink)
		// this racetrack will look like this. /race/link?t=<endpoint token>
		mlt.GET("/race/:link", validationWs.TokenVerifier(h.cfg), h.raceTrack)
	}

	profile := router.Group("/profile")
	profile.Use(access.OnlyRacer())
	{
		profile.GET("/info", h.profileInfo)
		profile.GET("/avatars", h.avatars)
		profile.PUT("/update", h.updateProfile)
		avatar := profile.Group("/avatar")
		{
			avatar.PUT("/update", h.updateAvatar)
		}
		history := profile.Group("/history")
		{
			history.GET("/single", h.singleHistory)

		}
	}

	return router
}

package rest

import (
	"github.com/baydogan/lnk/auth"
	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/rest/middleware"
	"github.com/baydogan/lnk/user"
	"github.com/gin-gonic/gin"
)

func NewRouter(mode string, h *URLHandler, userHandler *UserHandler, authSvc *auth.Service, userSvc *user.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", Health)

	api := router.Group("/api/v1")
	api.Use(middleware.Auth(authSvc))
	if mode == domain.ModeMulti {
		api.Use(middleware.WithRole(userSvc))
	}
	{
		api.POST("/shorten", h.ShortenURL)
		api.DELETE("/:code", h.DeleteURL)
		api.GET("/urls", h.ListURLs)
		api.GET("/urls/:code", h.StatsURL)

		if mode == domain.ModeMulti {
			api.GET("/me", userHandler.Whoami)

			users := api.Group("/users")
			users.Use(middleware.AdminOnly(userSvc))
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.ListUsers)
				users.DELETE("/:username", userHandler.DeleteUser)
			}
		}
	}

	router.NoRoute(h.RedirectURL)

	return router
}

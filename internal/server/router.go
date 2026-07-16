package server

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/middleware"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(mode string, h *handler.URLHandler, userHandler *handler.UserHandler, authSvc *service.AuthService, userSvc *service.UserService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.Health)

	api := router.Group("/api/v1")
	api.Use(middleware.Auth(authSvc))
	if mode == models.ModeMulti {
		api.Use(middleware.WithRole(userSvc))
	}
	{
		api.POST("/shorten", h.ShortenURL)
		api.DELETE("/:code", h.DeleteURL)
		api.GET("/urls", h.ListURLs)
		api.GET("/urls/:code", h.StatsURL)

		if mode == models.ModeMulti {
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

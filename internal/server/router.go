package server

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/middleware"
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.URLHandler, authSvc *service.AuthService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.Health)

	api := router.Group("/api/v1")
	api.Use(middleware.Auth(authSvc))
	{
		api.POST("/shorten", h.ShortenURL)
		api.DELETE("/:code", h.DeleteURL)
		api.GET("/urls", h.ListURLs)
		api.GET("/urls/:code", h.StatsURL)
	}

	router.NoRoute(h.RedirectURL)

	return router
}

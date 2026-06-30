package server

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.URLHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.Health)

	api := router.Group("/api/v1")
	{
		api.POST("/shorten", h.ShortenURL)
		api.DELETE("/:code", h.DeleteURL)
		api.GET("/urls", h.ListURLs)
	}

	// Public redirect: any unmatched root path is a short code/alias.
	router.NoRoute(h.RedirectURL)

	return router
}

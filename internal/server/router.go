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

	api := router.Group("/api")
	{
		api.POST("/shorten", h.ShortenURL)
		api.GET("/:code", h.RedirectURL)
		api.DELETE("/:code", h.DeleteURL)
		api.GET("/urls", h.ListURLs)
	}

	return router
}

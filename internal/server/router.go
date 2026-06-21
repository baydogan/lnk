package server

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", handler.Health)

	return router
}

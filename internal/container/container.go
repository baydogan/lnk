package container

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/repository"
	"github.com/baydogan/lnk/internal/service"
)

type Container struct {
	URLHandler  *handler.URLHandler
	AuthService *service.AuthService
}

func New(cfg models.ServerConfig) *Container {
	urlRepo := repository.NewURLRepository()
	keyRepo := repository.NewAPIKeyRepository()
	urlService := service.NewURLService(urlRepo, cfg.BaseURL)

	return &Container{
		URLHandler:  handler.NewHTTPHandler(urlService),
		AuthService: service.NewAuthService(keyRepo),
	}
}

package container

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/repository"
	"github.com/baydogan/lnk/internal/service"
)

type Container struct {
	URLHandler  *handler.URLHandler
	AuthService *service.AuthService
}

func New() *Container {
	urlRepo := repository.NewURLRepository()
	keyRepo := repository.NewAPIKeyRepository()
	urlService := service.NewURLService(urlRepo)

	return &Container{
		URLHandler:  handler.NewHTTPHandler(urlService),
		AuthService: service.NewAuthService(keyRepo),
	}
}

package container

import (
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/repository"
	"github.com/baydogan/lnk/internal/service"
)

type Container struct {
	URLHandler *handler.URLHandler
}

func New() *Container {
	urlRepo := repository.NewURLRepository()
	urlService := service.NewURLService(urlRepo)

	return &Container{
		URLHandler: handler.NewHTTPHandler(urlService),
	}
}

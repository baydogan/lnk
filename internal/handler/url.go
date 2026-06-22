package handler

import (
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	svc *service.URLService
}

func NewHTTPHandler(svc *service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {}

func (h *URLHandler) RedirectURL(c *gin.Context) {}

func (h *URLHandler) DeleteURL(c *gin.Context) {}

func (h *URLHandler) ListURLs(c *gin.Context) {}

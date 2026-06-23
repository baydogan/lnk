package handler

import (
	"net/http"

	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	svc *service.URLService
}

func NewHTTPHandler(svc *service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body. 'url' field is required and must be a valid URL.",
		})
		return
	}

	resp, err := h.svc.ShortenURL(&req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {}

func (h *URLHandler) DeleteURL(c *gin.Context) {}

func (h *URLHandler) ListURLs(c *gin.Context) {}

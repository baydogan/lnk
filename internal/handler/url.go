package handler

import (
	"net/http"
	"strings"

	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	resp, err := h.svc.ShortenURL(c.Request.Context(), &req, callerID(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func callerID(c *gin.Context) *bson.ObjectID {
	if v, ok := c.Get("user_id"); ok {
		if id, ok := v.(bson.ObjectID); ok {
			return &id
		}
	}
	return nil
}

func callerScope(c *gin.Context) *bson.ObjectID {
	if b, ok := c.Get("is_admin"); ok && b == true {
		return nil
	}
	return callerID(c)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	code := strings.Trim(c.Request.URL.Path, "/")
	if code == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	target, err := h.svc.ResolveURL(c.Request.Context(), code)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Redirect(http.StatusFound, target)
}

func (h *URLHandler) DeleteURL(c *gin.Context) {
	if err := h.svc.DeleteURL(c.Request.Context(), c.Param("code"), callerScope(c)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *URLHandler) ListURLs(c *gin.Context) {
	urls, err := h.svc.ListURLs(c.Request.Context(), callerScope(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, urls)
}

func (h *URLHandler) StatsURL(c *gin.Context) {
	u, err := h.svc.GetURL(c.Request.Context(), c.Param("code"), callerScope(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

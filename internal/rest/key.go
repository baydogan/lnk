package rest

import (
	"net/http"

	"github.com/baydogan/lnk/auth"
	"github.com/baydogan/lnk/domain"
	"github.com/gin-gonic/gin"
)

type KeyHandler struct {
	svc *auth.Service
}

func NewKeyHandler(svc *auth.Service) *KeyHandler {
	return &KeyHandler{svc: svc}
}

func (h *KeyHandler) Rotate(c *gin.Context) {
	v, ok := c.Get("api_key")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	key := v.(*domain.APIKey)

	plaintext, err := h.svc.RotateKey(c.Request.Context(), key.ID, key.UserID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.KeyResponse{APIKey: plaintext})
}

package rest

import (
	"errors"
	"net/http"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/gin-gonic/gin"
)

var errStatus = map[error]int{
	domain.ErrInvalidURL:        http.StatusBadRequest,
	domain.ErrExpireFormat:      http.StatusBadRequest,
	domain.ErrInvalidUsername:   http.StatusBadRequest,
	domain.ErrInvalidRole:       http.StatusBadRequest,
	domain.ErrAliasExists:       http.StatusConflict,
	domain.ErrAlreadyExists:     http.StatusConflict,
	domain.ErrURLLimit:          http.StatusForbidden,
	domain.ErrCannotDeleteAdmin: http.StatusForbidden,
	domain.ErrNotFound:          http.StatusNotFound,
}

func respondError(c *gin.Context, err error) {
	for known, status := range errStatus {
		if errors.Is(err, known) {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
	}
	logger.Error().Err(err).Msg("unhandled error")
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

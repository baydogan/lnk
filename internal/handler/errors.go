package handler

import (
	"errors"
	"net/http"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/gin-gonic/gin"
)

var errStatus = map[error]int{
	errs.ErrInvalidURL:        http.StatusBadRequest,
	errs.ErrExpireFormat:      http.StatusBadRequest,
	errs.ErrInvalidUsername:   http.StatusBadRequest,
	errs.ErrInvalidRole:       http.StatusBadRequest,
	errs.ErrAliasExists:       http.StatusConflict,
	errs.ErrAlreadyExists:     http.StatusConflict,
	errs.ErrURLLimit:          http.StatusForbidden,
	errs.ErrCannotDeleteAdmin: http.StatusForbidden,
	errs.ErrNotFound:          http.StatusNotFound,
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

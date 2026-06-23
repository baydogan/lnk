package handler

import (
	"errors"
	"net/http"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/gin-gonic/gin"
)

var errStatus = map[error]int{
	errs.ErrInvalidURL:    http.StatusBadRequest, // 400
	errs.ErrExpireFormat:  http.StatusBadRequest, // 400
	errs.ErrAliasExists:   http.StatusConflict,   // 409
	errs.ErrAlreadyExists: http.StatusConflict,   // 409
	errs.ErrURLLimit:      http.StatusForbidden,  // 403
	errs.ErrNotFound:      http.StatusNotFound,   // 404
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

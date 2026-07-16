package middleware

import (
	"net/http"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func WithRole(userSvc *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user_id")
		if !ok {
			c.Next()
			return
		}
		id, ok := v.(bson.ObjectID)
		if !ok {
			c.Next()
			return
		}
		if user, err := userSvc.GetUser(c.Request.Context(), id); err == nil && user.Role == domain.RoleAdmin {
			c.Set("is_admin", true)
		}
		c.Next()
	}
}

func AdminOnly(userSvc *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		id, ok := v.(bson.ObjectID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		user, err := userSvc.GetUser(c.Request.Context(), id)
		if err != nil || user.Role != domain.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

package rest

import (
	"net/http"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role,omitempty"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	role := req.Role
	if role == "" {
		role = domain.RoleUser
	}

	user, key, err := h.svc.CreateUser(c.Request.Context(), req.Username, role)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, domain.CreateUserResponse{User: *user, APIKey: key})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	if err := h.svc.DeleteUser(c.Request.Context(), c.Param("username")); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Whoami(c *gin.Context) {
	id := callerID(c)
	if id == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	user, err := h.svc.GetUser(c.Request.Context(), *id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

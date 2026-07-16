package rest_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/baydogan/lnk/auth"
	authmocks "github.com/baydogan/lnk/auth/mocks"
	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/rest"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func newKeyEngine(withKey bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := rest.NewKeyHandler(auth.NewService(authmocks.NewKeyRepository()))
	r := gin.New()
	key := &domain.APIKey{ID: bson.NewObjectID()}
	r.POST("/api/v1/keys/rotate", func(c *gin.Context) {
		if withKey {
			c.Set("api_key", key)
		}
	}, h.Rotate)
	return r
}

func TestRotateKeyHandler(t *testing.T) {
	w := do(newKeyEngine(true), "POST", "/api/v1/keys/rotate", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "api_key") {
		t.Fatalf("body missing api_key: %s", w.Body.String())
	}
}

func TestRotateKeyHandlerUnauthenticated(t *testing.T) {
	w := do(newKeyEngine(false), "POST", "/api/v1/keys/rotate", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

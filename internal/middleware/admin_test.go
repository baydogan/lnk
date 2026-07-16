package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type stubUserStore struct {
	user *models.User
	err  error
}

func (s stubUserStore) Create(context.Context, *models.User) error { return nil }
func (s stubUserStore) GetByUsername(context.Context, string) (*models.User, error) {
	return nil, errs.ErrNotFound
}
func (s stubUserStore) GetByID(context.Context, bson.ObjectID) (*models.User, error) {
	return s.user, s.err
}
func (s stubUserStore) List(context.Context) ([]models.User, error)    { return nil, nil }
func (s stubUserStore) EnsureIndexes(context.Context) error            { return nil }
func (s stubUserStore) DeleteByUsername(context.Context, string) error { return nil }

func runAdminOnly(t *testing.T, user *models.User, lookupErr error, setID bool) int {
	t.Helper()
	gin.SetMode(gin.TestMode)
	svc := service.NewUserService(stubUserStore{user: user, err: lookupErr}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	if setID {
		c.Set("user_id", bson.NewObjectID())
	}
	AdminOnly(svc)(c)
	if !c.IsAborted() {
		c.Status(http.StatusOK)
	}
	return w.Code
}

func TestAdminOnlyAllowsAdmin(t *testing.T) {
	if code := runAdminOnly(t, &models.User{Role: models.RoleAdmin}, nil, true); code != http.StatusOK {
		t.Fatalf("admin got %d, want 200", code)
	}
}

func TestAdminOnlyRejectsNonAdmin(t *testing.T) {
	if code := runAdminOnly(t, &models.User{Role: models.RoleUser}, nil, true); code != http.StatusForbidden {
		t.Fatalf("user got %d, want 403", code)
	}
}

func TestAdminOnlyRejectsMissingUserID(t *testing.T) {
	if code := runAdminOnly(t, &models.User{Role: models.RoleAdmin}, nil, false); code != http.StatusForbidden {
		t.Fatalf("no user_id got %d, want 403", code)
	}
}

func TestAdminOnlyRejectsLookupError(t *testing.T) {
	if code := runAdminOnly(t, nil, errs.ErrNotFound, true); code != http.StatusForbidden {
		t.Fatalf("lookup error got %d, want 403", code)
	}
}

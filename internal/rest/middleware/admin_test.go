package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baydogan/lnk/domain"
	usersvc "github.com/baydogan/lnk/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type stubUserStore struct {
	user *domain.User
	err  error
}

func (s stubUserStore) Create(context.Context, *domain.User) error { return nil }
func (s stubUserStore) GetByUsername(context.Context, string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}
func (s stubUserStore) GetByID(context.Context, bson.ObjectID) (*domain.User, error) {
	return s.user, s.err
}
func (s stubUserStore) List(context.Context) ([]domain.User, error)    { return nil, nil }
func (s stubUserStore) EnsureIndexes(context.Context) error            { return nil }
func (s stubUserStore) DeleteByUsername(context.Context, string) error { return nil }

func runAdminOnly(t *testing.T, user *domain.User, lookupErr error, setID bool) int {
	t.Helper()
	gin.SetMode(gin.TestMode)
	svc := usersvc.NewService(stubUserStore{user: user, err: lookupErr}, nil)

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
	if code := runAdminOnly(t, &domain.User{Role: domain.RoleAdmin}, nil, true); code != http.StatusOK {
		t.Fatalf("admin got %d, want 200", code)
	}
}

func TestAdminOnlyRejectsNonAdmin(t *testing.T) {
	if code := runAdminOnly(t, &domain.User{Role: domain.RoleUser}, nil, true); code != http.StatusForbidden {
		t.Fatalf("user got %d, want 403", code)
	}
}

func TestAdminOnlyRejectsMissingUserID(t *testing.T) {
	if code := runAdminOnly(t, &domain.User{Role: domain.RoleAdmin}, nil, false); code != http.StatusForbidden {
		t.Fatalf("no user_id got %d, want 403", code)
	}
}

func TestAdminOnlyRejectsLookupError(t *testing.T) {
	if code := runAdminOnly(t, nil, domain.ErrNotFound, true); code != http.StatusForbidden {
		t.Fatalf("lookup error got %d, want 403", code)
	}
}

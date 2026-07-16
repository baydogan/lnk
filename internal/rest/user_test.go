package rest_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/rest"
	"github.com/baydogan/lnk/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeUserStore struct {
	list      []domain.User
	getUser   *domain.User
	getByID   *domain.User
	createErr error
}

func (f *fakeUserStore) Create(_ context.Context, u *domain.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	u.ID = bson.NewObjectID()
	return nil
}
func (f *fakeUserStore) GetByUsername(context.Context, string) (*domain.User, error) {
	if f.getUser != nil {
		return f.getUser, nil
	}
	return nil, domain.ErrNotFound
}
func (f *fakeUserStore) GetByID(context.Context, bson.ObjectID) (*domain.User, error) {
	if f.getByID != nil {
		return f.getByID, nil
	}
	return nil, domain.ErrNotFound
}
func (f *fakeUserStore) List(context.Context) ([]domain.User, error)    { return f.list, nil }
func (f *fakeUserStore) EnsureIndexes(context.Context) error            { return nil }
func (f *fakeUserStore) DeleteByUsername(context.Context, string) error { return nil }

type fakeKeyStore struct{}

func (fakeKeyStore) Count(context.Context) (int64, error)         { return 0, nil }
func (fakeKeyStore) Create(context.Context, *domain.APIKey) error { return nil }
func (fakeKeyStore) EnsureIndexes(context.Context) error          { return nil }
func (fakeKeyStore) GetByHash(context.Context, string) (*domain.APIKey, error) {
	return nil, domain.ErrNotFound
}
func (fakeKeyStore) TouchLastUsed(context.Context, bson.ObjectID) error  { return nil }
func (fakeKeyStore) DeleteByUserID(context.Context, bson.ObjectID) error { return nil }

func newUserEngine(users user.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := rest.NewUserHandler(user.NewService(users, fakeKeyStore{}))
	r := gin.New()
	r.POST("/api/v1/users", h.CreateUser)
	r.GET("/api/v1/users", h.ListUsers)
	r.DELETE("/api/v1/users/:username", h.DeleteUser)
	return r
}

func TestCreateUserHandlerReturnsKey(t *testing.T) {
	w := do(newUserEngine(&fakeUserStore{}), "POST", "/api/v1/users", `{"username":"bob"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "api_key") {
		t.Fatalf("body missing api_key: %s", w.Body.String())
	}
}

func TestCreateUserHandlerBadBody(t *testing.T) {
	w := do(newUserEngine(&fakeUserStore{}), "POST", "/api/v1/users", `{"nope":1}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestCreateUserHandlerDuplicate(t *testing.T) {
	w := do(newUserEngine(&fakeUserStore{createErr: domain.ErrAlreadyExists}), "POST", "/api/v1/users", `{"username":"dup"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
}

func TestListUsersHandler(t *testing.T) {
	store := &fakeUserStore{list: []domain.User{{Username: "a", Role: domain.RoleUser}, {Username: "b", Role: domain.RoleAdmin}}}
	w := do(newUserEngine(store), "GET", "/api/v1/users", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	store := &fakeUserStore{getUser: &domain.User{Username: "bob", Role: domain.RoleUser}}
	w := do(newUserEngine(store), "DELETE", "/api/v1/users/bob", "")
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
}

func TestDeleteUserHandlerNotFound(t *testing.T) {
	w := do(newUserEngine(&fakeUserStore{}), "DELETE", "/api/v1/users/ghost", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestDeleteUserHandlerAdminForbidden(t *testing.T) {
	store := &fakeUserStore{getUser: &domain.User{Username: "root", Role: domain.RoleAdmin}}
	w := do(newUserEngine(store), "DELETE", "/api/v1/users/root", "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestWhoamiHandler(t *testing.T) {
	id := bson.NewObjectID()
	store := &fakeUserStore{getByID: &domain.User{ID: id, Username: "bob", Role: domain.RoleUser}}
	gin.SetMode(gin.TestMode)
	h := rest.NewUserHandler(user.NewService(store, fakeKeyStore{}))
	r := gin.New()
	r.GET("/api/v1/me", func(c *gin.Context) { c.Set("user_id", id) }, h.Whoami)

	w := do(r, "GET", "/api/v1/me", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "bob") {
		t.Fatalf("body missing username: %s", w.Body.String())
	}
}

func TestWhoamiHandlerNoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := rest.NewUserHandler(user.NewService(&fakeUserStore{}, fakeKeyStore{}))
	r := gin.New()
	r.GET("/api/v1/me", h.Whoami)

	w := do(r, "GET", "/api/v1/me", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

package handler_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/handler"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeUserStore struct {
	list      []models.User
	getUser   *models.User
	getByID   *models.User
	createErr error
}

func (f *fakeUserStore) Create(_ context.Context, u *models.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	u.ID = bson.NewObjectID()
	return nil
}
func (f *fakeUserStore) GetByUsername(context.Context, string) (*models.User, error) {
	if f.getUser != nil {
		return f.getUser, nil
	}
	return nil, errs.ErrNotFound
}
func (f *fakeUserStore) GetByID(context.Context, bson.ObjectID) (*models.User, error) {
	if f.getByID != nil {
		return f.getByID, nil
	}
	return nil, errs.ErrNotFound
}
func (f *fakeUserStore) List(context.Context) ([]models.User, error)    { return f.list, nil }
func (f *fakeUserStore) EnsureIndexes(context.Context) error            { return nil }
func (f *fakeUserStore) DeleteByUsername(context.Context, string) error { return nil }

type fakeKeyStore struct{}

func (fakeKeyStore) Count(context.Context) (int64, error)         { return 0, nil }
func (fakeKeyStore) Create(context.Context, *models.APIKey) error { return nil }
func (fakeKeyStore) EnsureIndexes(context.Context) error          { return nil }
func (fakeKeyStore) GetByHash(context.Context, string) (*models.APIKey, error) {
	return nil, errs.ErrNotFound
}
func (fakeKeyStore) TouchLastUsed(context.Context, bson.ObjectID) error  { return nil }
func (fakeKeyStore) DeleteByUserID(context.Context, bson.ObjectID) error { return nil }

func newUserEngine(users service.UserStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := handler.NewUserHandler(service.NewUserService(users, fakeKeyStore{}))
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
	w := do(newUserEngine(&fakeUserStore{createErr: errs.ErrAlreadyExists}), "POST", "/api/v1/users", `{"username":"dup"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
}

func TestListUsersHandler(t *testing.T) {
	store := &fakeUserStore{list: []models.User{{Username: "a", Role: models.RoleUser}, {Username: "b", Role: models.RoleAdmin}}}
	w := do(newUserEngine(store), "GET", "/api/v1/users", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	store := &fakeUserStore{getUser: &models.User{Username: "bob", Role: models.RoleUser}}
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
	store := &fakeUserStore{getUser: &models.User{Username: "root", Role: models.RoleAdmin}}
	w := do(newUserEngine(store), "DELETE", "/api/v1/users/root", "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestWhoamiHandler(t *testing.T) {
	id := bson.NewObjectID()
	store := &fakeUserStore{getByID: &models.User{ID: id, Username: "bob", Role: models.RoleUser}}
	gin.SetMode(gin.TestMode)
	h := handler.NewUserHandler(service.NewUserService(store, fakeKeyStore{}))
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
	h := handler.NewUserHandler(service.NewUserService(&fakeUserStore{}, fakeKeyStore{}))
	r := gin.New()
	r.GET("/api/v1/me", h.Whoami)

	w := do(r, "GET", "/api/v1/me", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

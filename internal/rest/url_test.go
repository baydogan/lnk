package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/rest"
	"github.com/baydogan/lnk/url"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeStore struct {
	url       *domain.URL
	all       []domain.URL
	exists    bool
	getErr    error
	createErr error
	deleteErr error
}

func (f *fakeStore) CreateURL(context.Context, *domain.URL) error { return f.createErr }

func (f *fakeStore) GetByCodeOrAlias(context.Context, string) (*domain.URL, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.url, nil
}

func (f *fakeStore) IncrementClickCount(context.Context, string) error { return nil }
func (f *fakeStore) GetURLsByOwner(context.Context, *bson.ObjectID) ([]domain.URL, error) {
	return f.all, nil
}
func (f *fakeStore) CodeExists(context.Context, string) (bool, error) { return f.exists, nil }
func (f *fakeStore) DeleteByCode(context.Context, string) error       { return f.deleteErr }

func newEngine(store url.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := rest.NewHTTPHandler(url.NewService(store, "http://x"))
	r := gin.New()
	r.GET("/health", rest.Health)
	r.POST("/api/v1/shorten", h.ShortenURL)
	r.DELETE("/api/v1/:code", h.DeleteURL)
	r.GET("/api/v1/urls", h.ListURLs)
	r.GET("/api/v1/urls/:code", h.StatsURL)
	r.NoRoute(h.RedirectURL)
	return r
}

func do(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealthHandler(t *testing.T) {
	if w := do(newEngine(&fakeStore{}), "GET", "/health", ""); w.Code != http.StatusOK {
		t.Fatalf("health status = %d", w.Code)
	}
}

func TestShortenHandlerCreated(t *testing.T) {
	w := do(newEngine(&fakeStore{}), "POST", "/api/v1/shorten", `{"url":"example.com"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "short_url") {
		t.Fatalf("body missing short_url: %s", w.Body.String())
	}
}

func TestShortenHandlerBadBody(t *testing.T) {
	w := do(newEngine(&fakeStore{}), "POST", "/api/v1/shorten", `{"nope":1}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestShortenHandlerAliasConflict(t *testing.T) {
	w := do(newEngine(&fakeStore{exists: true}), "POST", "/api/v1/shorten", `{"url":"a.com","alias":"taken"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
}

func TestRedirectHandlerFound(t *testing.T) {
	store := &fakeStore{url: &domain.URL{Code: "abc", OriginalURL: "https://target.com"}}
	w := do(newEngine(store), "GET", "/abc", "")
	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "https://target.com" {
		t.Fatalf("Location = %q", loc)
	}
}

func TestRedirectHandlerNotFound(t *testing.T) {
	store := &fakeStore{getErr: domain.ErrNotFound}
	w := do(newEngine(store), "GET", "/missing", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestRedirectHandlerRejectsNonGet(t *testing.T) {
	w := do(newEngine(&fakeStore{}), "POST", "/abc", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestDeleteHandler(t *testing.T) {
	if w := do(newEngine(&fakeStore{}), "DELETE", "/api/v1/abc", ""); w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
	store := &fakeStore{deleteErr: domain.ErrNotFound}
	if w := do(newEngine(store), "DELETE", "/api/v1/gone", ""); w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestListHandler(t *testing.T) {
	store := &fakeStore{all: []domain.URL{{Code: "a"}, {Code: "b"}}}
	w := do(newEngine(store), "GET", "/api/v1/urls", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestStatsHandler(t *testing.T) {
	store := &fakeStore{url: &domain.URL{Code: "abc"}}
	if w := do(newEngine(store), "GET", "/api/v1/urls/abc", ""); w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	miss := &fakeStore{getErr: domain.ErrNotFound}
	if w := do(newEngine(miss), "GET", "/api/v1/urls/none", ""); w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

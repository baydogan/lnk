package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetsBearerWhenTokenPresent(t *testing.T) {
	var gotAuth, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.SetToken("lnk_secret")
	if _, _, err := c.Get("/ping"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if gotAuth != "Bearer lnk_secret" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Fatalf("Content-Type = %q", gotContentType)
	}
}

func TestNoAuthHeaderWithoutToken(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
	}))
	defer srv.Close()

	if _, _, err := New(srv.URL).Get("/ping"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if gotAuth != "" {
		t.Fatalf("Authorization = %q, want empty", gotAuth)
	}
}

func TestPostSendsJSONBody(t *testing.T) {
	var gotMethod string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	body, status, err := New(srv.URL).Post("/shorten", map[string]string{"url": "https://x.com"})
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q", gotMethod)
	}
	if status != http.StatusCreated {
		t.Fatalf("status = %d, want 201", status)
	}
	if string(gotBody) != `{"url":"https://x.com"}` {
		t.Fatalf("body = %s", gotBody)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("resp body = %s", body)
	}
}

func TestDeleteMethodAndStatusPassthrough(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	_, status, err := New(srv.URL).Delete("/abc")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if status != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", status)
	}
}

func TestRequestErrorOnBadURL(t *testing.T) {
	if _, _, err := New("http://127.0.0.1:1").Get("/x"); err == nil {
		t.Fatal("expected error hitting unreachable server")
	}
}

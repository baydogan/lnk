package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
)

func TestShortenURLPrependsScheme(t *testing.T) {
	store := newFakeURLStore()
	svc := NewURLService(store, "http://localhost:8080")

	resp, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: "example.com"})
	if err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if resp.OriginalURL != "https://example.com" {
		t.Fatalf("OriginalURL = %q, want https://example.com", resp.OriginalURL)
	}
	if resp.ShortURL != "http://localhost:8080/"+resp.Code {
		t.Fatalf("ShortURL = %q", resp.ShortURL)
	}
	if len(store.created) != 1 || store.created[0].OriginalURL != "https://example.com" {
		t.Fatalf("stored url unexpected: %+v", store.created)
	}
}

func TestShortenURLKeepsExistingScheme(t *testing.T) {
	svc := NewURLService(newFakeURLStore(), "http://x")
	resp, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: "http://a.com"})
	if err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if resp.OriginalURL != "http://a.com" {
		t.Fatalf("OriginalURL = %q, want unchanged", resp.OriginalURL)
	}
}

func TestShortenURLEmpty(t *testing.T) {
	svc := NewURLService(newFakeURLStore(), "http://x")
	for _, in := range []string{"", "   "} {
		if _, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: in}); !errors.Is(err, errs.ErrInvalidURL) {
			t.Fatalf("ShortenURL(%q) err = %v, want ErrInvalidURL", in, err)
		}
	}
}

func TestShortenURLAliasExists(t *testing.T) {
	store := newFakeURLStore()
	store.existing["taken"] = true
	svc := NewURLService(store, "http://x")

	_, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: "a.com", Alias: "taken"})
	if !errors.Is(err, errs.ErrAliasExists) {
		t.Fatalf("err = %v, want ErrAliasExists", err)
	}
}

func TestShortenURLWithAliasUsesAliasInShortURL(t *testing.T) {
	store := newFakeURLStore()
	svc := NewURLService(store, "http://x")

	resp, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: "a.com", Alias: "mylink"})
	if err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if resp.ShortURL != "http://x/mylink" {
		t.Fatalf("ShortURL = %q, want http://x/mylink", resp.ShortURL)
	}
	if store.created[0].Alias == nil || *store.created[0].Alias != "mylink" {
		t.Fatalf("stored alias unexpected: %+v", store.created[0])
	}
}

func TestShortenURLInvalidExpiry(t *testing.T) {
	svc := NewURLService(newFakeURLStore(), "http://x")
	_, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: "a.com", Expires: "nope"})
	if !errors.Is(err, errs.ErrExpireFormat) {
		t.Fatalf("err = %v, want ErrExpireFormat", err)
	}
}

func TestShortenURLTrimsTrailingSlashInBaseURL(t *testing.T) {
	svc := NewURLService(newFakeURLStore(), "http://x/")
	resp, err := svc.ShortenURL(context.Background(), &models.ShortenRequest{URL: "a.com"})
	if err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if resp.ShortURL != "http://x/"+resp.Code {
		t.Fatalf("ShortURL = %q", resp.ShortURL)
	}
}

func TestResolveURLHappy(t *testing.T) {
	store := newFakeURLStore()
	store.byKey["abc"] = &models.URL{Code: "abc", OriginalURL: "https://target.com"}
	svc := NewURLService(store, "http://x")

	got, err := svc.ResolveURL(context.Background(), "abc")
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if got != "https://target.com" {
		t.Fatalf("target = %q", got)
	}
	if len(store.incremented) != 1 || store.incremented[0] != "abc" {
		t.Fatalf("click not incremented: %+v", store.incremented)
	}
}

func TestResolveURLExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	store := newFakeURLStore()
	store.byKey["abc"] = &models.URL{Code: "abc", OriginalURL: "https://target.com", ExpiresAt: &past}
	svc := NewURLService(store, "http://x")

	if _, err := svc.ResolveURL(context.Background(), "abc"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if len(store.incremented) != 0 {
		t.Fatalf("expired link should not increment clicks: %+v", store.incremented)
	}
}

func TestResolveURLNotFound(t *testing.T) {
	svc := NewURLService(newFakeURLStore(), "http://x")
	if _, err := svc.ResolveURL(context.Background(), "missing"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestResolveURLIncrementErrorIgnored(t *testing.T) {
	store := newFakeURLStore()
	store.byKey["abc"] = &models.URL{Code: "abc", OriginalURL: "https://target.com"}
	store.incErr = errors.New("mongo down")
	svc := NewURLService(store, "http://x")

	got, err := svc.ResolveURL(context.Background(), "abc")
	if err != nil {
		t.Fatalf("ResolveURL should ignore increment error, got %v", err)
	}
	if got != "https://target.com" {
		t.Fatalf("target = %q", got)
	}
}

func TestListURLsBuildsShortURL(t *testing.T) {
	alias := "vanity"
	store := newFakeURLStore()
	store.all = []models.URL{
		{Code: "code1"},
		{Code: "code2", Alias: &alias},
	}
	svc := NewURLService(store, "http://x")

	out, err := svc.ListURLs(context.Background())
	if err != nil {
		t.Fatalf("ListURLs: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2", len(out))
	}
	if out[0].ShortURL != "http://x/code1" {
		t.Fatalf("out[0].ShortURL = %q", out[0].ShortURL)
	}
	if out[1].ShortURL != "http://x/vanity" {
		t.Fatalf("out[1].ShortURL = %q", out[1].ShortURL)
	}
}

func TestGetURL(t *testing.T) {
	store := newFakeURLStore()
	store.byKey["abc"] = &models.URL{Code: "abc", OriginalURL: "https://t.com"}
	svc := NewURLService(store, "http://x")

	r, err := svc.GetURL(context.Background(), "abc")
	if err != nil {
		t.Fatalf("GetURL: %v", err)
	}
	if r.ShortURL != "http://x/abc" {
		t.Fatalf("ShortURL = %q", r.ShortURL)
	}
}

func TestDeleteURL(t *testing.T) {
	store := newFakeURLStore()
	svc := NewURLService(store, "http://x")
	if err := svc.DeleteURL(context.Background(), "abc"); err != nil {
		t.Fatalf("DeleteURL: %v", err)
	}
	if len(store.deleted) != 1 || store.deleted[0] != "abc" {
		t.Fatalf("deleted = %+v", store.deleted)
	}

	store.deleteErr = errs.ErrNotFound
	if err := svc.DeleteURL(context.Background(), "gone"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

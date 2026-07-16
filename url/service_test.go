package url

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/url/mocks"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestParseExpiryEmpty(t *testing.T) {
	for _, in := range []string{"", "   ", "\t"} {
		got, err := parseExpiry(in)
		if err != nil || got != nil {
			t.Fatalf("parseExpiry(%q) = %v, %v", in, got, err)
		}
	}
}

func TestParseExpiryValid(t *testing.T) {
	cases := map[string]time.Duration{
		"30m": 30 * time.Minute, "1h": time.Hour, "2h30m": 2*time.Hour + 30*time.Minute,
		"1.5h": 90 * time.Minute, "7d": 7 * 24 * time.Hour, "2w": 14 * 24 * time.Hour, "1d": 24 * time.Hour,
	}
	for in, want := range cases {
		before := time.Now().Add(want)
		got, err := parseExpiry(in)
		after := time.Now().Add(want)
		if err != nil || got == nil {
			t.Fatalf("parseExpiry(%q) = %v, %v", in, got, err)
		}
		if got.Before(before.Add(-time.Second)) || got.After(after.Add(time.Second)) {
			t.Fatalf("parseExpiry(%q) = %v, want ~%v", in, got, before)
		}
	}
}

func TestParseExpiryInvalid(t *testing.T) {
	for _, in := range []string{"0d", "-1h", "abc", "10", "d", "w", "1x", "0", "-3d", "1.5"} {
		if _, err := parseExpiry(in); !errors.Is(err, domain.ErrExpireFormat) {
			t.Fatalf("parseExpiry(%q) err = %v, want ErrExpireFormat", in, err)
		}
	}
}

func TestShortenURLPrependsScheme(t *testing.T) {
	repo := mocks.NewRepository()
	svc := NewService(repo, "http://localhost:8080")

	resp, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: "example.com"}, nil)
	if err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if resp.OriginalURL != "https://example.com" {
		t.Fatalf("OriginalURL = %q", resp.OriginalURL)
	}
	if resp.ShortURL != "http://localhost:8080/"+resp.Code {
		t.Fatalf("ShortURL = %q", resp.ShortURL)
	}
}

func TestShortenURLKeepsExistingScheme(t *testing.T) {
	svc := NewService(mocks.NewRepository(), "http://x")
	resp, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: "http://a.com"}, nil)
	if err != nil || resp.OriginalURL != "http://a.com" {
		t.Fatalf("OriginalURL = %q, err %v", resp.OriginalURL, err)
	}
}

func TestShortenURLEmpty(t *testing.T) {
	svc := NewService(mocks.NewRepository(), "http://x")
	for _, in := range []string{"", "   "} {
		if _, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: in}, nil); !errors.Is(err, domain.ErrInvalidURL) {
			t.Fatalf("err = %v, want ErrInvalidURL", err)
		}
	}
}

func TestShortenURLAliasExists(t *testing.T) {
	repo := mocks.NewRepository()
	repo.Existing["taken"] = true
	svc := NewService(repo, "http://x")
	if _, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: "a.com", Alias: "taken"}, nil); !errors.Is(err, domain.ErrAliasExists) {
		t.Fatalf("err = %v, want ErrAliasExists", err)
	}
}

func TestShortenURLWithAliasUsesAlias(t *testing.T) {
	repo := mocks.NewRepository()
	svc := NewService(repo, "http://x")
	resp, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: "a.com", Alias: "mylink"}, nil)
	if err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if resp.ShortURL != "http://x/mylink" {
		t.Fatalf("ShortURL = %q", resp.ShortURL)
	}
}

func TestShortenURLInvalidExpiry(t *testing.T) {
	svc := NewService(mocks.NewRepository(), "http://x")
	if _, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: "a.com", Expires: "nope"}, nil); !errors.Is(err, domain.ErrExpireFormat) {
		t.Fatalf("err = %v, want ErrExpireFormat", err)
	}
}

func TestShortenURLSetsOwner(t *testing.T) {
	repo := mocks.NewRepository()
	svc := NewService(repo, "http://x")
	owner := bson.NewObjectID()
	if _, err := svc.ShortenURL(context.Background(), &domain.ShortenRequest{URL: "a.com"}, &owner); err != nil {
		t.Fatalf("ShortenURL: %v", err)
	}
	if repo.Created[0].UserID == nil || *repo.Created[0].UserID != owner {
		t.Fatalf("UserID = %v, want %v", repo.Created[0].UserID, owner)
	}
}

func TestResolveURLHappy(t *testing.T) {
	repo := mocks.NewRepository()
	repo.ByKey["abc"] = &domain.URL{Code: "abc", OriginalURL: "https://target.com"}
	svc := NewService(repo, "http://x")
	got, err := svc.ResolveURL(context.Background(), "abc")
	if err != nil || got != "https://target.com" {
		t.Fatalf("got %q, err %v", got, err)
	}
	if len(repo.Incremented) != 1 {
		t.Fatalf("click not incremented")
	}
}

func TestResolveURLExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	repo := mocks.NewRepository()
	repo.ByKey["abc"] = &domain.URL{Code: "abc", OriginalURL: "https://t.com", ExpiresAt: &past}
	svc := NewService(repo, "http://x")
	if _, err := svc.ResolveURL(context.Background(), "abc"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if len(repo.Incremented) != 0 {
		t.Fatal("expired should not increment")
	}
}

func TestListURLsFiltersByOwner(t *testing.T) {
	alice := bson.NewObjectID()
	bob := bson.NewObjectID()
	repo := mocks.NewRepository()
	repo.All = []domain.URL{{Code: "a1", UserID: &alice}, {Code: "a2", UserID: &alice}, {Code: "b1", UserID: &bob}}
	svc := NewService(repo, "http://x")

	mine, err := svc.ListURLs(context.Background(), &alice)
	if err != nil || len(mine) != 2 {
		t.Fatalf("alice sees %d, err %v", len(mine), err)
	}
	all, err := svc.ListURLs(context.Background(), nil)
	if err != nil || len(all) != 3 {
		t.Fatalf("all %d, err %v", len(all), err)
	}
}

func TestGetURLHidesOtherOwners(t *testing.T) {
	alice := bson.NewObjectID()
	bob := bson.NewObjectID()
	repo := mocks.NewRepository()
	repo.ByKey["abc"] = &domain.URL{Code: "abc", UserID: &alice}
	svc := NewService(repo, "http://x")

	if _, err := svc.GetURL(context.Background(), "abc", &bob); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("bob err = %v, want ErrNotFound", err)
	}
	if _, err := svc.GetURL(context.Background(), "abc", &alice); err != nil {
		t.Fatalf("owner err = %v", err)
	}
	if _, err := svc.GetURL(context.Background(), "abc", nil); err != nil {
		t.Fatalf("admin err = %v", err)
	}
}

func TestDeleteURLDeniesOtherOwners(t *testing.T) {
	alice := bson.NewObjectID()
	bob := bson.NewObjectID()
	repo := mocks.NewRepository()
	repo.ByKey["abc"] = &domain.URL{Code: "abc", UserID: &alice}
	svc := NewService(repo, "http://x")

	if err := svc.DeleteURL(context.Background(), "abc", &bob); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("bob err = %v, want ErrNotFound", err)
	}
	if len(repo.Deleted) != 0 {
		t.Fatal("non-owner delete should not remove")
	}
	if err := svc.DeleteURL(context.Background(), "abc", &alice); err != nil {
		t.Fatalf("owner delete err = %v", err)
	}
	if len(repo.Deleted) != 1 {
		t.Fatal("owner delete should remove")
	}
}

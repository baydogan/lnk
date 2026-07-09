//go:build integration

package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestURLCreateAndGet(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	u := &models.URL{Code: "abc", OriginalURL: "https://x.com"}
	if err := repo.CreateURL(ctx, u); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}
	if u.ID.IsZero() {
		t.Fatal("CreateURL did not set ID")
	}
	if u.CreatedAt.IsZero() || u.UpdatedAt.IsZero() {
		t.Fatal("CreateURL did not set timestamps")
	}

	got, err := repo.GetURLByCode(ctx, "abc")
	if err != nil {
		t.Fatalf("GetURLByCode: %v", err)
	}
	if got.OriginalURL != "https://x.com" {
		t.Fatalf("OriginalURL = %q", got.OriginalURL)
	}
}

func TestGetURLByCodeNotFound(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	if _, err := repo.GetURLByCode(context.Background(), "missing"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestGetByCodeOrAlias(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	alias := "vanity"
	if err := repo.CreateURL(ctx, &models.URL{Code: "code1", OriginalURL: "https://x.com", Alias: &alias}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}

	byCode, err := repo.GetByCodeOrAlias(ctx, "code1")
	if err != nil || byCode.Code != "code1" {
		t.Fatalf("by code: %+v err %v", byCode, err)
	}
	byAlias, err := repo.GetByCodeOrAlias(ctx, "vanity")
	if err != nil || byAlias.Code != "code1" {
		t.Fatalf("by alias: %+v err %v", byAlias, err)
	}
	if _, err := repo.GetByCodeOrAlias(ctx, "nope"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("missing err = %v, want ErrNotFound", err)
	}
}

func TestCodeExists(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	alias := "al"
	if err := repo.CreateURL(ctx, &models.URL{Code: "c", OriginalURL: "https://x.com", Alias: &alias}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}

	for _, k := range []string{"c", "al"} {
		exists, err := repo.CodeExists(ctx, k)
		if err != nil || !exists {
			t.Fatalf("CodeExists(%q) = %v, %v; want true", k, exists, err)
		}
	}
	exists, err := repo.CodeExists(ctx, "free")
	if err != nil || exists {
		t.Fatalf("CodeExists(free) = %v, %v; want false", exists, err)
	}
}

func TestIncrementClickCount(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	if err := repo.CreateURL(ctx, &models.URL{Code: "hit", OriginalURL: "https://x.com"}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}
	for range 3 {
		if err := repo.IncrementClickCount(ctx, "hit"); err != nil {
			t.Fatalf("IncrementClickCount: %v", err)
		}
	}
	got, err := repo.GetURLByCode(ctx, "hit")
	if err != nil {
		t.Fatalf("GetURLByCode: %v", err)
	}
	if got.ClickCount != 3 {
		t.Fatalf("ClickCount = %d, want 3", got.ClickCount)
	}
}

func TestGetAllURLs(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	for _, c := range []string{"a", "b", "c"} {
		if err := repo.CreateURL(ctx, &models.URL{Code: c, OriginalURL: "https://x.com"}); err != nil {
			t.Fatalf("CreateURL: %v", err)
		}
	}
	all, err := repo.GetAllURLs(ctx)
	if err != nil {
		t.Fatalf("GetAllURLs: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len = %d, want 3", len(all))
	}
}

func TestDeleteByCode(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	alias := "al"
	if err := repo.CreateURL(ctx, &models.URL{Code: "c", OriginalURL: "https://x.com", Alias: &alias}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}

	if err := repo.DeleteByCode(ctx, "al"); err != nil {
		t.Fatalf("DeleteByCode by alias: %v", err)
	}
	if _, err := repo.GetURLByCode(ctx, "c"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("after delete err = %v, want ErrNotFound", err)
	}
	if err := repo.DeleteByCode(ctx, "gone"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("delete missing err = %v, want ErrNotFound", err)
	}
}

func TestCountByUserID(t *testing.T) {
	clearCollection(t, "urls")
	repo := NewURLRepository()
	ctx := context.Background()

	userID := bson.NewObjectID()
	other := bson.NewObjectID()
	if err := repo.CreateURL(ctx, &models.URL{Code: "u1", OriginalURL: "https://x.com", UserID: &userID}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}
	if err := repo.CreateURL(ctx, &models.URL{Code: "u2", OriginalURL: "https://x.com", UserID: &userID}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}
	if err := repo.CreateURL(ctx, &models.URL{Code: "o1", OriginalURL: "https://x.com", UserID: &other}); err != nil {
		t.Fatalf("CreateURL: %v", err)
	}

	n, err := repo.CountByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("CountByUserID: %v", err)
	}
	if n != 2 {
		t.Fatalf("count = %d, want 2", n)
	}
}

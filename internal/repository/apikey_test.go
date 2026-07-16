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

func TestAPIKeyCreateAndGetByHash(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()

	k := &models.APIKey{KeyHash: "hash1", Prefix: "lnk_abc12345"}
	if err := repo.Create(ctx, k); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if k.ID.IsZero() || k.CreatedAt.IsZero() {
		t.Fatal("Create did not set ID/CreatedAt")
	}

	got, err := repo.GetByHash(ctx, "hash1")
	if err != nil {
		t.Fatalf("GetByHash: %v", err)
	}
	if got.Prefix != "lnk_abc12345" {
		t.Fatalf("Prefix = %q", got.Prefix)
	}
}

func TestAPIKeyGetByHashNotFound(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	if _, err := repo.GetByHash(context.Background(), "nope"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestAPIKeyCount(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()

	n, err := repo.Count(ctx)
	if err != nil || n != 0 {
		t.Fatalf("Count empty = %d, %v; want 0", n, err)
	}
	if err := repo.Create(ctx, &models.APIKey{KeyHash: "h", Prefix: "p"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if n, _ := repo.Count(ctx); n != 1 {
		t.Fatalf("Count = %d, want 1", n)
	}
}

func TestAPIKeyTouchLastUsed(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()

	k := &models.APIKey{KeyHash: "h", Prefix: "p"}
	if err := repo.Create(ctx, k); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if k.LastUsedAt != nil {
		t.Fatal("LastUsedAt should be nil before touch")
	}
	if err := repo.TouchLastUsed(ctx, k.ID); err != nil {
		t.Fatalf("TouchLastUsed: %v", err)
	}
	got, err := repo.GetByHash(ctx, "h")
	if err != nil {
		t.Fatalf("GetByHash: %v", err)
	}
	if got.LastUsedAt == nil {
		t.Fatal("LastUsedAt still nil after touch")
	}
}

func TestAPIKeyDuplicateHash(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("EnsureIndexes: %v", err)
	}

	first := &models.APIKey{KeyHash: "same", Prefix: "p1"}
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	dup := &models.APIKey{KeyHash: "same", Prefix: "p2"}
	if err := repo.Create(ctx, dup); !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("duplicate hash err = %v, want ErrAlreadyExists", err)
	}
}

func TestAPIKeyOneKeyPerNilUser(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("EnsureIndexes: %v", err)
	}

	if err := repo.Create(ctx, &models.APIKey{KeyHash: "h1", Prefix: "p1"}); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	err := repo.Create(ctx, &models.APIKey{KeyHash: "h2", Prefix: "p2"})
	if !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("second nil-user key err = %v, want ErrAlreadyExists (uniq_user_id)", err)
	}
}

func TestAPIKeyDeleteByUserID(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()

	uid := bson.NewObjectID()
	other := bson.NewObjectID()
	if err := repo.Create(ctx, &models.APIKey{KeyHash: "h1", Prefix: "p1", UserID: &uid}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := repo.Create(ctx, &models.APIKey{KeyHash: "h2", Prefix: "p2", UserID: &other}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := repo.DeleteByUserID(ctx, uid); err != nil {
		t.Fatalf("DeleteByUserID: %v", err)
	}
	if _, err := repo.GetByHash(ctx, "h1"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("deleted user's key still present, err = %v", err)
	}
	if _, err := repo.GetByHash(ctx, "h2"); err != nil {
		t.Fatalf("other user's key should remain: %v", err)
	}
}

func TestAPIKeyEnsureIndexesIdempotent(t *testing.T) {
	clearCollection(t, "api_keys")
	repo := NewAPIKeyRepository()
	ctx := context.Background()
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("first EnsureIndexes: %v", err)
	}
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("second EnsureIndexes: %v", err)
	}
}

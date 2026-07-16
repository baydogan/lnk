//go:build integration

package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
)

func TestUserCreateAndGetByUsername(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	ctx := context.Background()

	u := &models.User{Username: "alice", Role: models.RoleAdmin}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u.ID.IsZero() || u.CreatedAt.IsZero() {
		t.Fatal("Create did not set ID/CreatedAt")
	}

	got, err := repo.GetByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("GetByUsername: %v", err)
	}
	if got.Role != models.RoleAdmin {
		t.Fatalf("Role = %q, want admin", got.Role)
	}
}

func TestUserGetByUsernameNotFound(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	if _, err := repo.GetByUsername(context.Background(), "ghost"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestUserGetByID(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	ctx := context.Background()

	u := &models.User{Username: "bob", Role: models.RoleUser}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := repo.GetByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Username != "bob" {
		t.Fatalf("Username = %q", got.Username)
	}
}

func TestUserDuplicateUsername(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	ctx := context.Background()
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("EnsureIndexes: %v", err)
	}

	if err := repo.Create(ctx, &models.User{Username: "dup", Role: models.RoleUser}); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	err := repo.Create(ctx, &models.User{Username: "dup", Role: models.RoleUser})
	if !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("duplicate username err = %v, want ErrAlreadyExists", err)
	}
}

func TestUserList(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	ctx := context.Background()

	for _, name := range []string{"a", "b", "c"} {
		if err := repo.Create(ctx, &models.User{Username: name, Role: models.RoleUser}); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}
	users, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(users) != 3 {
		t.Fatalf("len = %d, want 3", len(users))
	}
}

func TestUserDeleteByUsername(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	ctx := context.Background()

	if err := repo.Create(ctx, &models.User{Username: "bob", Role: models.RoleUser}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := repo.DeleteByUsername(ctx, "bob"); err != nil {
		t.Fatalf("DeleteByUsername: %v", err)
	}
	if _, err := repo.GetByUsername(ctx, "bob"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("after delete err = %v, want ErrNotFound", err)
	}
	if err := repo.DeleteByUsername(ctx, "ghost"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("delete missing err = %v, want ErrNotFound", err)
	}
}

func TestUserEnsureIndexesIdempotent(t *testing.T) {
	clearCollection(t, "users")
	repo := NewUserRepository()
	ctx := context.Background()
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("first EnsureIndexes: %v", err)
	}
	if err := repo.EnsureIndexes(ctx); err != nil {
		t.Fatalf("second EnsureIndexes: %v", err)
	}
}

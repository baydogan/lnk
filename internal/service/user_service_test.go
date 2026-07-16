package service

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
)

func TestCreateUserIssuesKeyBoundToUser(t *testing.T) {
	users := newFakeUserStore()
	keys := newFakeKeyStore()
	svc := NewUserService(users, keys)

	user, pt, err := svc.CreateUser(context.Background(), "alice", models.RoleUser)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if pt == "" {
		t.Fatal("expected non-empty plaintext key")
	}
	if user.Role != models.RoleUser {
		t.Fatalf("Role = %q", user.Role)
	}
	if len(keys.created) != 1 {
		t.Fatalf("expected one key issued, got %d", len(keys.created))
	}
	k := keys.created[0]
	if k.UserID == nil || *k.UserID != user.ID {
		t.Fatalf("key UserID = %v, want %v", k.UserID, user.ID)
	}
	if k.KeyHash != hashKey(pt) {
		t.Fatal("issued key hash does not match returned plaintext")
	}
}

func TestCreateUserTrimsAndValidates(t *testing.T) {
	svc := NewUserService(newFakeUserStore(), newFakeKeyStore())

	if _, _, err := svc.CreateUser(context.Background(), "   ", models.RoleUser); !errors.Is(err, errs.ErrInvalidUsername) {
		t.Fatalf("empty username err = %v, want ErrInvalidUsername", err)
	}
	if _, _, err := svc.CreateUser(context.Background(), "bob", "superuser"); !errors.Is(err, errs.ErrInvalidRole) {
		t.Fatalf("bad role err = %v, want ErrInvalidRole", err)
	}
}

func TestCreateUserDuplicatePropagates(t *testing.T) {
	users := newFakeUserStore()
	users.createErr = errs.ErrAlreadyExists
	svc := NewUserService(users, newFakeKeyStore())

	if _, _, err := svc.CreateUser(context.Background(), "dup", models.RoleUser); !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("err = %v, want ErrAlreadyExists", err)
	}
}

func TestListUsers(t *testing.T) {
	users := newFakeUserStore()
	svc := NewUserService(users, newFakeKeyStore())
	ctx := context.Background()

	if _, _, err := svc.CreateUser(ctx, "a", models.RoleUser); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if _, _, err := svc.CreateUser(ctx, "b", models.RoleAdmin); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	out, err := svc.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2", len(out))
	}
}

func TestEnsureAdminCreatesWhenAbsent(t *testing.T) {
	users := newFakeUserStore()
	keys := newFakeKeyStore()
	svc := NewUserService(users, keys)

	pt, created, err := svc.EnsureAdmin(context.Background(), "root")
	if err != nil {
		t.Fatalf("EnsureAdmin: %v", err)
	}
	if !created || pt == "" {
		t.Fatalf("created = %v, pt = %q", created, pt)
	}
	if u := users.byName["root"]; u == nil || u.Role != models.RoleAdmin {
		t.Fatalf("admin user not created with admin role: %+v", u)
	}
	if len(keys.created) != 1 || keys.created[0].UserID == nil {
		t.Fatal("admin key not issued bound to user")
	}
}

func TestDeleteUserRemovesUserAndKeys(t *testing.T) {
	users := newFakeUserStore()
	keys := newFakeKeyStore()
	svc := NewUserService(users, keys)
	ctx := context.Background()

	u, _, err := svc.CreateUser(ctx, "bob", models.RoleUser)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if err := svc.DeleteUser(ctx, "bob"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if _, err := users.GetByUsername(ctx, "bob"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("user still present, err = %v", err)
	}
	if len(keys.deletedUsers) != 1 || keys.deletedUsers[0] != u.ID {
		t.Fatalf("user's keys not deleted: %+v", keys.deletedUsers)
	}
}

func TestDeleteUserRefusesAdmin(t *testing.T) {
	users := newFakeUserStore()
	keys := newFakeKeyStore()
	svc := NewUserService(users, keys)
	ctx := context.Background()
	if _, _, err := svc.CreateUser(ctx, "root", models.RoleAdmin); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := svc.DeleteUser(ctx, "root"); !errors.Is(err, errs.ErrCannotDeleteAdmin) {
		t.Fatalf("err = %v, want ErrCannotDeleteAdmin", err)
	}
	if len(keys.deletedUsers) != 0 {
		t.Fatal("admin guard should not delete keys")
	}
}

func TestGetUser(t *testing.T) {
	users := newFakeUserStore()
	svc := NewUserService(users, newFakeKeyStore())
	ctx := context.Background()
	u, _, err := svc.CreateUser(ctx, "bob", models.RoleUser)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	got, err := svc.GetUser(ctx, u.ID)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Username != "bob" {
		t.Fatalf("Username = %q", got.Username)
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	svc := NewUserService(newFakeUserStore(), newFakeKeyStore())
	if err := svc.DeleteUser(context.Background(), "ghost"); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestEnsureAdminSkipsWhenPresent(t *testing.T) {
	users := newFakeUserStore()
	keys := newFakeKeyStore()
	svc := NewUserService(users, keys)
	if _, _, err := svc.CreateUser(context.Background(), "root", models.RoleAdmin); err != nil {
		t.Fatalf("seed CreateUser: %v", err)
	}

	pt, created, err := svc.EnsureAdmin(context.Background(), "root")
	if err != nil || created || pt != "" {
		t.Fatalf("EnsureAdmin = (%q, %v, %v), want empty/false/nil", pt, created, err)
	}
	if len(keys.created) != 1 {
		t.Fatalf("no extra key should be issued, got %d", len(keys.created))
	}
}

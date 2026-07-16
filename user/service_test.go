package user

import (
	"context"
	"errors"
	"testing"

	authmocks "github.com/baydogan/lnk/auth/mocks"
	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/user/mocks"
)

func TestCreateUserIssuesKeyBoundToUser(t *testing.T) {
	users := mocks.NewRepository()
	keys := authmocks.NewKeyRepository()
	svc := NewService(users, keys)

	u, pt, err := svc.CreateUser(context.Background(), "alice", domain.RoleUser)
	if err != nil || pt == "" {
		t.Fatalf("CreateUser = %v, %q", err, pt)
	}
	if len(keys.Created) != 1 {
		t.Fatalf("keys created = %d", len(keys.Created))
	}
	k := keys.Created[0]
	if k.UserID == nil || *k.UserID != u.ID {
		t.Fatalf("key UserID = %v, want %v", k.UserID, u.ID)
	}
	if k.KeyHash != domain.HashKey(pt) {
		t.Fatal("hash mismatch")
	}
}

func TestCreateUserValidates(t *testing.T) {
	svc := NewService(mocks.NewRepository(), authmocks.NewKeyRepository())
	if _, _, err := svc.CreateUser(context.Background(), "  ", domain.RoleUser); !errors.Is(err, domain.ErrInvalidUsername) {
		t.Fatalf("err = %v, want ErrInvalidUsername", err)
	}
	if _, _, err := svc.CreateUser(context.Background(), "bob", "root"); !errors.Is(err, domain.ErrInvalidRole) {
		t.Fatalf("err = %v, want ErrInvalidRole", err)
	}
}

func TestCreateUserDuplicate(t *testing.T) {
	users := mocks.NewRepository()
	users.CreateErr = domain.ErrAlreadyExists
	svc := NewService(users, authmocks.NewKeyRepository())
	if _, _, err := svc.CreateUser(context.Background(), "dup", domain.RoleUser); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("err = %v", err)
	}
}

func TestDeleteUserRemovesUserAndKeys(t *testing.T) {
	users := mocks.NewRepository()
	keys := authmocks.NewKeyRepository()
	svc := NewService(users, keys)
	u, _, _ := svc.CreateUser(context.Background(), "bob", domain.RoleUser)

	if err := svc.DeleteUser(context.Background(), "bob"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if _, err := users.GetByUsername(context.Background(), "bob"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatal("user still present")
	}
	if len(keys.DeletedUsers) != 1 || keys.DeletedUsers[0] != u.ID {
		t.Fatal("keys not deleted")
	}
}

func TestDeleteUserRefusesAdmin(t *testing.T) {
	users := mocks.NewRepository()
	keys := authmocks.NewKeyRepository()
	svc := NewService(users, keys)
	_, _, _ = svc.CreateUser(context.Background(), "root", domain.RoleAdmin)
	if err := svc.DeleteUser(context.Background(), "root"); !errors.Is(err, domain.ErrCannotDeleteAdmin) {
		t.Fatalf("err = %v, want ErrCannotDeleteAdmin", err)
	}
	if len(keys.DeletedUsers) != 0 {
		t.Fatal("admin guard should not delete keys")
	}
}

func TestEnsureAdminCreatesWhenAbsent(t *testing.T) {
	users := mocks.NewRepository()
	keys := authmocks.NewKeyRepository()
	svc := NewService(users, keys)
	pt, created, err := svc.EnsureAdmin(context.Background(), "root")
	if err != nil || !created || pt == "" {
		t.Fatalf("= (%q, %v, %v)", pt, created, err)
	}
	if u := users.ByName["root"]; u == nil || u.Role != domain.RoleAdmin {
		t.Fatal("admin not created")
	}
}

func TestEnsureAdminSkipsWhenPresent(t *testing.T) {
	users := mocks.NewRepository()
	keys := authmocks.NewKeyRepository()
	svc := NewService(users, keys)
	_, _, _ = svc.CreateUser(context.Background(), "root", domain.RoleAdmin)
	pt, created, err := svc.EnsureAdmin(context.Background(), "root")
	if err != nil || created || pt != "" {
		t.Fatalf("= (%q, %v, %v)", pt, created, err)
	}
}

func TestGetUser(t *testing.T) {
	users := mocks.NewRepository()
	svc := NewService(users, authmocks.NewKeyRepository())
	u, _, _ := svc.CreateUser(context.Background(), "bob", domain.RoleUser)
	got, err := svc.GetUser(context.Background(), u.ID)
	if err != nil || got.Username != "bob" {
		t.Fatalf("GetUser = %v, %v", got, err)
	}
}

func TestListUsers(t *testing.T) {
	users := mocks.NewRepository()
	svc := NewService(users, authmocks.NewKeyRepository())
	_, _, _ = svc.CreateUser(context.Background(), "a", domain.RoleUser)
	_, _, _ = svc.CreateUser(context.Background(), "b", domain.RoleAdmin)
	out, err := svc.ListUsers(context.Background())
	if err != nil || len(out) != 2 {
		t.Fatalf("len %d, err %v", len(out), err)
	}
}

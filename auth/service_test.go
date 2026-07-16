package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/auth/mocks"
	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestAuthenticateValid(t *testing.T) {
	repo := mocks.NewKeyRepository()
	id := bson.NewObjectID()
	repo.ByHash[domain.HashKey("lnk_secret")] = &domain.APIKey{ID: id, KeyHash: domain.HashKey("lnk_secret")}
	svc := NewService(repo)

	key, err := svc.Authenticate(context.Background(), "lnk_secret")
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if key.ID != id {
		t.Fatal("wrong key")
	}
	if len(repo.Touched) != 1 || repo.Touched[0] != id {
		t.Fatal("last_used not touched")
	}
}

func TestAuthenticateInvalid(t *testing.T) {
	svc := NewService(mocks.NewKeyRepository())
	if _, err := svc.Authenticate(context.Background(), "lnk_wrong"); !errors.Is(err, domain.ErrInvalidKey) {
		t.Fatalf("err = %v, want ErrInvalidKey", err)
	}
}

func TestAuthenticateRepoErrorPassthrough(t *testing.T) {
	repo := mocks.NewKeyRepository()
	repo.GetErr = errors.New("mongo down")
	svc := NewService(repo)
	if _, err := svc.Authenticate(context.Background(), "lnk_x"); err == nil || errors.Is(err, domain.ErrInvalidKey) {
		t.Fatalf("err = %v, want raw repo error", err)
	}
}

func TestEnsureAdminKeyCreates(t *testing.T) {
	repo := mocks.NewKeyRepository()
	svc := NewService(repo)
	pt, created, err := svc.EnsureAdminKey(context.Background())
	if err != nil || !created || pt == "" {
		t.Fatalf("= (%q, %v, %v)", pt, created, err)
	}
	if len(repo.Created) != 1 || repo.Created[0].KeyHash != domain.HashKey(pt) {
		t.Fatal("stored key mismatch")
	}
}

func TestEnsureAdminKeyExisting(t *testing.T) {
	repo := mocks.NewKeyRepository()
	repo.CountVal = 1
	svc := NewService(repo)
	pt, created, err := svc.EnsureAdminKey(context.Background())
	if err != nil || created || pt != "" || len(repo.Created) != 0 {
		t.Fatalf("= (%q, %v, %v), created keys %d", pt, created, err, len(repo.Created))
	}
}

func TestEnsureAdminKeyDuplicateRace(t *testing.T) {
	repo := mocks.NewKeyRepository()
	repo.CreateErr = domain.ErrAlreadyExists
	svc := NewService(repo)
	pt, created, err := svc.EnsureAdminKey(context.Background())
	if err != nil || created || pt != "" {
		t.Fatalf("duplicate race = (%q, %v, %v)", pt, created, err)
	}
}

func TestEnsureIndexesDelegates(t *testing.T) {
	repo := mocks.NewKeyRepository()
	repo.IdxErr = errors.New("index failed")
	svc := NewService(repo)
	if err := svc.EnsureIndexes(context.Background()); err != repo.IdxErr {
		t.Fatalf("err = %v", err)
	}
}

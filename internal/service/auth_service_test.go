package service

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestAuthenticateValid(t *testing.T) {
	store := newFakeKeyStore()
	id := bson.NewObjectID()
	store.byHash[hashKey("lnk_secret")] = &models.APIKey{ID: id, KeyHash: hashKey("lnk_secret")}
	svc := NewAuthService(store)

	key, err := svc.Authenticate(context.Background(), "lnk_secret")
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if key.ID != id {
		t.Fatalf("returned wrong key: %+v", key)
	}
	if len(store.touched) != 1 || store.touched[0] != id {
		t.Fatalf("last_used not touched: %+v", store.touched)
	}
}

func TestAuthenticateInvalid(t *testing.T) {
	svc := NewAuthService(newFakeKeyStore())
	if _, err := svc.Authenticate(context.Background(), "lnk_wrong"); !errors.Is(err, errs.ErrInvalidKey) {
		t.Fatalf("err = %v, want ErrInvalidKey", err)
	}
}

func TestAuthenticateRepoErrorPassthrough(t *testing.T) {
	store := newFakeKeyStore()
	store.getErr = errors.New("mongo down")
	svc := NewAuthService(store)

	_, err := svc.Authenticate(context.Background(), "lnk_x")
	if err == nil || errors.Is(err, errs.ErrInvalidKey) {
		t.Fatalf("err = %v, want raw repo error", err)
	}
}

func TestAuthenticateTouchErrorIgnored(t *testing.T) {
	store := newFakeKeyStore()
	store.byHash[hashKey("lnk_x")] = &models.APIKey{ID: bson.NewObjectID(), KeyHash: hashKey("lnk_x")}
	store.touchErr = errors.New("touch failed")
	svc := NewAuthService(store)

	if _, err := svc.Authenticate(context.Background(), "lnk_x"); err != nil {
		t.Fatalf("Authenticate should ignore touch error, got %v", err)
	}
}

func TestEnsureAdminKeyCreates(t *testing.T) {
	store := newFakeKeyStore()
	svc := NewAuthService(store)

	pt, created, err := svc.EnsureAdminKey(context.Background())
	if err != nil {
		t.Fatalf("EnsureAdminKey: %v", err)
	}
	if !created || pt == "" {
		t.Fatalf("created = %v, plaintext = %q", created, pt)
	}
	if len(store.created) != 1 {
		t.Fatalf("expected one key created, got %d", len(store.created))
	}
	if store.created[0].KeyHash != hashKey(pt) {
		t.Fatal("stored hash does not match returned plaintext")
	}
	if store.created[0].Prefix != pt[:12] {
		t.Fatalf("prefix %q != first 12 of %q", store.created[0].Prefix, pt)
	}
}

func TestEnsureAdminKeyExisting(t *testing.T) {
	store := newFakeKeyStore()
	store.count = 1
	svc := NewAuthService(store)

	pt, created, err := svc.EnsureAdminKey(context.Background())
	if err != nil || created || pt != "" {
		t.Fatalf("EnsureAdminKey = (%q, %v, %v), want empty/false/nil", pt, created, err)
	}
	if len(store.created) != 0 {
		t.Fatal("should not create a key when one exists")
	}
}

func TestEnsureAdminKeyDuplicateRace(t *testing.T) {
	store := newFakeKeyStore()
	store.createErr = errs.ErrAlreadyExists
	svc := NewAuthService(store)

	pt, created, err := svc.EnsureAdminKey(context.Background())
	if err != nil {
		t.Fatalf("duplicate insert should not be fatal, got %v", err)
	}
	if created || pt != "" {
		t.Fatalf("EnsureAdminKey = (%q, %v), want empty/false", pt, created)
	}
}

func TestEnsureAdminKeyCountError(t *testing.T) {
	store := newFakeKeyStore()
	store.countErr = errors.New("mongo down")
	svc := NewAuthService(store)

	if _, _, err := svc.EnsureAdminKey(context.Background()); err == nil {
		t.Fatal("expected count error to propagate")
	}
}

func TestEnsureIndexesDelegates(t *testing.T) {
	store := newFakeKeyStore()
	store.idxErr = errors.New("index failed")
	svc := NewAuthService(store)

	if err := svc.EnsureIndexes(context.Background()); err != store.idxErr {
		t.Fatalf("EnsureIndexes err = %v, want delegated error", err)
	}
}

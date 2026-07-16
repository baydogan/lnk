//go:build integration

package mongo

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/domain"
)

func TestSettingsGetModeNotSet(t *testing.T) {
	clearCollection(t, "settings")
	repo := NewSettingsRepository(testDB)
	if _, err := repo.GetMode(context.Background()); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestSettingsSetAndGetMode(t *testing.T) {
	clearCollection(t, "settings")
	repo := NewSettingsRepository(testDB)
	ctx := context.Background()

	if err := repo.SetMode(ctx, "multi"); err != nil {
		t.Fatalf("SetMode: %v", err)
	}
	got, err := repo.GetMode(ctx)
	if err != nil || got != "multi" {
		t.Fatalf("GetMode = %q, %v", got, err)
	}
	if err := repo.SetMode(ctx, "single"); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("second SetMode err = %v, want ErrAlreadyExists (singleton)", err)
	}
}

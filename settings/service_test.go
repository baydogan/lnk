package settings

import (
	"context"
	"errors"
	"testing"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/settings/mocks"
)

func TestEnsureModePinsWhenUnset(t *testing.T) {
	repo := &mocks.Repository{}
	svc := NewService(repo)
	if err := svc.EnsureMode(context.Background(), "multi"); err != nil {
		t.Fatalf("EnsureMode: %v", err)
	}
	if len(repo.SetCalls) != 1 || repo.SetCalls[0] != "multi" {
		t.Fatalf("mode not pinned: %+v", repo.SetCalls)
	}
}

func TestEnsureModeMatch(t *testing.T) {
	m := "multi"
	repo := &mocks.Repository{Stored: &m}
	svc := NewService(repo)
	if err := svc.EnsureMode(context.Background(), "multi"); err != nil {
		t.Fatalf("EnsureMode: %v", err)
	}
	if len(repo.SetCalls) != 0 {
		t.Fatal("should not re-pin on match")
	}
}

func TestEnsureModeMismatch(t *testing.T) {
	m := "multi"
	repo := &mocks.Repository{Stored: &m}
	svc := NewService(repo)
	if err := svc.EnsureMode(context.Background(), "single"); !errors.Is(err, domain.ErrModeMismatch) {
		t.Fatalf("err = %v, want ErrModeMismatch", err)
	}
}

func TestEnsureModeRacePinnedByOther(t *testing.T) {
	other := "multi"
	repo := &mocks.Repository{SetErr: domain.ErrAlreadyExists, RaceValue: &other}
	svc := NewService(repo)
	if err := svc.EnsureMode(context.Background(), "multi"); err != nil {
		t.Fatalf("race match: %v", err)
	}
}

func TestEnsureModeRaceMismatch(t *testing.T) {
	other := "single"
	repo := &mocks.Repository{SetErr: domain.ErrAlreadyExists, RaceValue: &other}
	svc := NewService(repo)
	if err := svc.EnsureMode(context.Background(), "multi"); !errors.Is(err, domain.ErrModeMismatch) {
		t.Fatalf("err = %v, want ErrModeMismatch", err)
	}
}

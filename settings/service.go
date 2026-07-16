package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/baydogan/lnk/domain"
)

type Repository interface {
	GetMode(ctx context.Context) (string, error)
	SetMode(ctx context.Context, mode string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) EnsureMode(ctx context.Context, want string) error {
	got, err := s.repo.GetMode(ctx)
	if err == nil {
		return compare(got, want)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	if err := s.repo.SetMode(ctx, want); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			got, err := s.repo.GetMode(ctx)
			if err != nil {
				return err
			}
			return compare(got, want)
		}
		return err
	}
	return nil
}

func compare(got, want string) error {
	if got != want {
		return fmt.Errorf("%w: DB=%q, config/env=%q — kasıtlı migration gerekir (settings dokümanını elle güncelle)", domain.ErrModeMismatch, got, want)
	}
	return nil
}

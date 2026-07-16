package mocks

import (
	"context"

	"github.com/baydogan/lnk/domain"
)

type Repository struct {
	Stored    *string
	RaceValue *string
	SetCalls  []string
	GetErr    error
	SetErr    error
}

func (f *Repository) GetMode(context.Context) (string, error) {
	if f.GetErr != nil {
		return "", f.GetErr
	}
	if f.Stored == nil {
		return "", domain.ErrNotFound
	}
	return *f.Stored, nil
}

func (f *Repository) SetMode(_ context.Context, mode string) error {
	if f.SetErr != nil {
		if f.RaceValue != nil {
			f.Stored = f.RaceValue
		}
		return f.SetErr
	}
	f.SetCalls = append(f.SetCalls, mode)
	f.Stored = &mode
	return nil
}

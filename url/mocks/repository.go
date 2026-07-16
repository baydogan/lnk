package mocks

import (
	"context"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository struct {
	Created     []*domain.URL
	ByKey       map[string]*domain.URL
	Existing    map[string]bool
	All         []domain.URL
	Incremented []string
	Deleted     []string

	CreateErr error
	GetErr    error
	ExistsErr error
	ListErr   error
	IncErr    error
	DeleteErr error
}

func NewRepository() *Repository {
	return &Repository{ByKey: map[string]*domain.URL{}, Existing: map[string]bool{}}
}

func (f *Repository) CreateURL(_ context.Context, url *domain.URL) error {
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = append(f.Created, url)
	return nil
}

func (f *Repository) GetByCodeOrAlias(_ context.Context, s string) (*domain.URL, error) {
	if f.GetErr != nil {
		return nil, f.GetErr
	}
	if u, ok := f.ByKey[s]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (f *Repository) IncrementClickCount(_ context.Context, code string) error {
	f.Incremented = append(f.Incremented, code)
	return f.IncErr
}

func (f *Repository) GetURLsByOwner(_ context.Context, ownerID *bson.ObjectID) ([]domain.URL, error) {
	if f.ListErr != nil {
		return nil, f.ListErr
	}
	if ownerID == nil {
		return f.All, nil
	}
	var out []domain.URL
	for _, u := range f.All {
		if u.UserID != nil && *u.UserID == *ownerID {
			out = append(out, u)
		}
	}
	return out, nil
}

func (f *Repository) CodeExists(_ context.Context, code string) (bool, error) {
	if f.ExistsErr != nil {
		return false, f.ExistsErr
	}
	return f.Existing[code], nil
}

func (f *Repository) DeleteByCode(_ context.Context, code string) error {
	f.Deleted = append(f.Deleted, code)
	return f.DeleteErr
}

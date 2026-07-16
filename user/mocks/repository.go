package mocks

import (
	"context"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository struct {
	Created []*domain.User
	ByName  map[string]*domain.User
	Deleted []string

	CreateErr error
	GetErr    error
	ListErr   error
}

func NewRepository() *Repository {
	return &Repository{ByName: map[string]*domain.User{}}
}

func (f *Repository) Create(_ context.Context, u *domain.User) error {
	if f.CreateErr != nil {
		return f.CreateErr
	}
	u.ID = bson.NewObjectID()
	f.Created = append(f.Created, u)
	f.ByName[u.Username] = u
	return nil
}

func (f *Repository) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	if f.GetErr != nil {
		return nil, f.GetErr
	}
	if u, ok := f.ByName[username]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (f *Repository) GetByID(_ context.Context, id bson.ObjectID) (*domain.User, error) {
	for _, u := range f.ByName {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *Repository) List(context.Context) ([]domain.User, error) {
	if f.ListErr != nil {
		return nil, f.ListErr
	}
	out := make([]domain.User, 0, len(f.Created))
	for _, u := range f.Created {
		out = append(out, *u)
	}
	return out, nil
}

func (f *Repository) EnsureIndexes(context.Context) error { return nil }

func (f *Repository) DeleteByUsername(_ context.Context, username string) error {
	if _, ok := f.ByName[username]; !ok {
		return domain.ErrNotFound
	}
	delete(f.ByName, username)
	f.Deleted = append(f.Deleted, username)
	return nil
}

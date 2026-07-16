package mocks

import (
	"context"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type KeyRepository struct {
	CountVal     int64
	Created      []*domain.APIKey
	ByHash       map[string]*domain.APIKey
	Touched      []bson.ObjectID
	DeletedUsers []bson.ObjectID

	CountErr  error
	CreateErr error
	GetErr    error
	TouchErr  error
	IdxErr    error
	DeleteErr error
}

func NewKeyRepository() *KeyRepository {
	return &KeyRepository{ByHash: map[string]*domain.APIKey{}}
}

func (f *KeyRepository) Count(context.Context) (int64, error) { return f.CountVal, f.CountErr }

func (f *KeyRepository) Create(_ context.Context, k *domain.APIKey) error {
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Created = append(f.Created, k)
	return nil
}

func (f *KeyRepository) EnsureIndexes(context.Context) error { return f.IdxErr }

func (f *KeyRepository) GetByHash(_ context.Context, hash string) (*domain.APIKey, error) {
	if f.GetErr != nil {
		return nil, f.GetErr
	}
	if k, ok := f.ByHash[hash]; ok {
		return k, nil
	}
	return nil, domain.ErrNotFound
}

func (f *KeyRepository) TouchLastUsed(_ context.Context, id bson.ObjectID) error {
	f.Touched = append(f.Touched, id)
	return f.TouchErr
}

func (f *KeyRepository) DeleteByUserID(_ context.Context, id bson.ObjectID) error {
	f.DeletedUsers = append(f.DeletedUsers, id)
	return f.DeleteErr
}

package auth

import (
	"context"
	"errors"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, k *domain.APIKey) error
	EnsureIndexes(ctx context.Context) error
	GetByHash(ctx context.Context, hash string) (*domain.APIKey, error)
	TouchLastUsed(ctx context.Context, id bson.ObjectID) error
	DeleteByID(ctx context.Context, id bson.ObjectID) error
}

type Service struct {
	keys Repository
}

func NewService(keys Repository) *Service {
	return &Service{keys: keys}
}

func (s *Service) EnsureAdminKey(ctx context.Context) (plaintext string, created bool, err error) {
	n, err := s.keys.Count(ctx)
	if err != nil {
		return "", false, err
	}
	if n > 0 {
		return "", false, nil
	}
	pt, hash, prefix, err := domain.GenerateAPIKey()
	if err != nil {
		return "", false, err
	}
	if err := s.keys.Create(ctx, &domain.APIKey{KeyHash: hash, Prefix: prefix}); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return "", false, nil
		}
		return "", false, err
	}
	return pt, true, nil
}

func (s *Service) EnsureIndexes(ctx context.Context) error {
	return s.keys.EnsureIndexes(ctx)
}

func (s *Service) RotateKey(ctx context.Context, oldKeyID bson.ObjectID, userID *bson.ObjectID) (string, error) {
	if err := s.keys.DeleteByID(ctx, oldKeyID); err != nil {
		return "", err
	}
	plaintext, hash, prefix, err := domain.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	if err := s.keys.Create(ctx, &domain.APIKey{KeyHash: hash, Prefix: prefix, UserID: userID}); err != nil {
		return "", err
	}
	return plaintext, nil
}

func (s *Service) Authenticate(ctx context.Context, plaintext string) (*domain.APIKey, error) {
	key, err := s.keys.GetByHash(ctx, domain.HashKey(plaintext))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidKey
		}
		return nil, err
	}
	_ = s.keys.TouchLastUsed(ctx, key.ID)
	return key, nil
}

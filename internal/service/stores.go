package service

import (
	"context"

	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type URLStore interface {
	CreateURL(ctx context.Context, url *models.URL) error
	GetByCodeOrAlias(ctx context.Context, s string) (*models.URL, error)
	IncrementClickCount(ctx context.Context, code string) error
	GetAllURLs(ctx context.Context) ([]models.URL, error)
	CodeExists(ctx context.Context, code string) (bool, error)
	DeleteByCode(ctx context.Context, code string) error
}

type APIKeyStore interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, k *models.APIKey) error
	EnsureIndexes(ctx context.Context) error
	GetByHash(ctx context.Context, hash string) (*models.APIKey, error)
	TouchLastUsed(ctx context.Context, id bson.ObjectID) error
}

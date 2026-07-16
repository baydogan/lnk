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
	GetURLsByOwner(ctx context.Context, ownerID *bson.ObjectID) ([]models.URL, error)
	CodeExists(ctx context.Context, code string) (bool, error)
	DeleteByCode(ctx context.Context, code string) error
}

type APIKeyStore interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, k *models.APIKey) error
	EnsureIndexes(ctx context.Context) error
	GetByHash(ctx context.Context, hash string) (*models.APIKey, error)
	TouchLastUsed(ctx context.Context, id bson.ObjectID) error
	DeleteByUserID(ctx context.Context, userID bson.ObjectID) error
}

type UserStore interface {
	Create(ctx context.Context, u *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	EnsureIndexes(ctx context.Context) error
	DeleteByUsername(ctx context.Context, username string) error
}

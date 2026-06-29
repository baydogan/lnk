package repository

import (
	"context"
	"time"

	"github.com/baydogan/lnk/internal/database"
	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type APIKeyRepository struct{ col *mongo.Collection }

func NewAPIKeyRepository() *APIKeyRepository {
	return &APIKeyRepository{col: database.Collection("api_keys")}
}

func (r *APIKeyRepository) Count() (int64, error) {
	return r.col.CountDocuments(context.Background(), bson.M{})
}

func (r *APIKeyRepository) Create(k *models.APIKey) error {
	k.ID = bson.NewObjectID()
	k.CreatedAt = time.Now()
	_, err := r.col.InsertOne(context.Background(), k)
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrAlreadyExists
	}
	return err
}

func (r *APIKeyRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key_hash", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_key_hash"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_user_id"),
		},
	})
	return err
}

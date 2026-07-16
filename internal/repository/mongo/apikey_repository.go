package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type APIKeyRepository struct{ col *mongo.Collection }

func NewAPIKeyRepository(db *mongo.Database) *APIKeyRepository {
	return &APIKeyRepository{col: db.Collection("api_keys")}
}

func (r *APIKeyRepository) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}

func (r *APIKeyRepository) Create(ctx context.Context, k *domain.APIKey) error {
	k.ID = bson.NewObjectID()
	k.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, k)
	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrAlreadyExists
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

func (r *APIKeyRepository) GetByHash(ctx context.Context, hash string) (*domain.APIKey, error) {
	var k domain.APIKey
	err := r.col.FindOne(ctx, bson.M{"key_hash": hash}).Decode(&k)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *APIKeyRepository) TouchLastUsed(ctx context.Context, id bson.ObjectID) error {
	_, err := r.col.UpdateByID(ctx, id,
		bson.M{"$set": bson.M{"last_used_at": time.Now()}})
	return err
}

func (r *APIKeyRepository) DeleteByUserID(ctx context.Context, userID bson.ObjectID) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}

func (r *APIKeyRepository) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

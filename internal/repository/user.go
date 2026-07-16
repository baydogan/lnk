package repository

import (
	"context"
	"errors"
	"time"

	"github.com/baydogan/lnk/internal/database"
	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct{ col *mongo.Collection }

func NewUserRepository() *UserRepository {
	return &UserRepository{col: database.Collection("users")}
}

func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_username"),
	})
	return err
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	u.ID = bson.NewObjectID()
	u.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrAlreadyExists
	}
	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errs.ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	var u models.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errs.ErrNotFound
	}
	return &u, err
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) DeleteByUsername(ctx context.Context, username string) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errs.ErrNotFound
	}
	return nil
}

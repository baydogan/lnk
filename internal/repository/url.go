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
)

type URLRepository struct {
	col *mongo.Collection
}

func NewURLRepository() *URLRepository {
	return &URLRepository{col: database.Collection("urls")}
}

func (r *URLRepository) CreateURL(url *models.URL) error {
	url.ID = bson.NewObjectID()
	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()

	_, err := r.col.InsertOne(context.Background(), url)
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrAlreadyExists
	}
	return err
}

func (r *URLRepository) GetURLByCode(code string) (*models.URL, error) {
	var url models.URL
	err := r.col.FindOne(context.Background(), bson.M{"code": code}).Decode(&url)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errs.ErrNotFound
	}
	return &url, err
}

func (r *URLRepository) GetByCodeOrAlias(s string) (*models.URL, error) {
	var url models.URL
	err := r.col.FindOne(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{"code": s},
			bson.M{"alias": s},
		},
	}).Decode(&url)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errs.ErrNotFound
	}
	return &url, err
}

func (r *URLRepository) IncrementClickCount(code string) error {
	_, err := r.col.UpdateOne(
		context.Background(),
		bson.M{"code": code},
		bson.M{"$inc": bson.M{"click_count": 1}, "$set": bson.M{"updated_at": time.Now()}},
	)
	return err
}

func (r *URLRepository) GetAllURLs() ([]models.URL, error) {
	cursor, err := r.col.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var urls []models.URL
	if err = cursor.All(context.Background(), &urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *URLRepository) CodeExists(code string) (bool, error) {
	count, err := r.col.CountDocuments(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{"code": code},
			bson.M{"alias": code},
		},
	})
	return count > 0, err
}

func (r *URLRepository) CountByUserID(userID bson.ObjectID) (int64, error) {
	return r.col.CountDocuments(context.Background(), bson.M{"user_id": userID})
}

func (r *URLRepository) DeleteByCode(code string) error {
	res, err := r.col.DeleteOne(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{"code": code},
			bson.M{"alias": code},
		},
	})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errs.ErrNotFound
	}
	return nil
}

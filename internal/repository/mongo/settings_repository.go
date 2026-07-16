package mongo

import (
	"context"
	"errors"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const deploymentSettingsID = "deployment"

type SettingsRepository struct{ col *mongo.Collection }

func NewSettingsRepository(db *mongo.Database) *SettingsRepository {
	return &SettingsRepository{col: db.Collection("settings")}
}

func (r *SettingsRepository) GetMode(ctx context.Context) (string, error) {
	var doc struct {
		Mode string `bson:"mode"`
	}
	err := r.col.FindOne(ctx, bson.M{"_id": deploymentSettingsID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return doc.Mode, nil
}

func (r *SettingsRepository) SetMode(ctx context.Context, mode string) error {
	_, err := r.col.InsertOne(ctx, bson.M{"_id": deploymentSettingsID, "mode": mode})
	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrAlreadyExists
	}
	return err
}

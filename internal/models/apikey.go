package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type APIKey struct {
	ID         bson.ObjectID  `bson:"_id,omitempty"`
	KeyHash    string         `bson:"key_hash"`
	Prefix     string         `bson:"prefix"`
	UserID     *bson.ObjectID `bson:"user_id,omitempty"` // single'da nil
	CreatedAt  time.Time      `bson:"created_at"`
	LastUsedAt *time.Time     `bson:"last_used_at,omitempty"`
}

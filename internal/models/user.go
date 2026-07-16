package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string        `bson:"username" json:"username"`
	Role      string        `bson:"role" json:"role"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type CreateUserResponse struct {
	User   User   `json:"user"`
	APIKey string `json:"api_key"`
}

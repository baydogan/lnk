package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type URL struct {
	ID          bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Code        string         `bson:"code" json:"code"`
	OriginalURL string         `bson:"original_url" json:"original_url"`
	UserID      *bson.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	ClickCount  int            `bson:"click_count" json:"click_count"`
	Alias       *string        `bson:"alias,omitempty" json:"alias,omitempty"`
	ExpiresAt   *time.Time     `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	CreatedAt   time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at" json:"updated_at"`
}

type URLResponse struct {
	ID          bson.ObjectID  `json:"id"`
	Code        string         `json:"code"`
	ShortURL    string         `json:"short_url"`
	OriginalURL string         `json:"original_url"`
	UserID      *bson.ObjectID `json:"user_id,omitempty"`
	ClickCount  int            `json:"click_count"`
	Alias       *string        `json:"alias,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (u *URL) ToResponse() URLResponse {
	return URLResponse{
		ID:          u.ID,
		Code:        u.Code,
		OriginalURL: u.OriginalURL,
		UserID:      u.UserID,
		ClickCount:  u.ClickCount,
		Alias:       u.Alias,
		ExpiresAt:   u.ExpiresAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type ShortenRequest struct {
	URL     string `json:"url" binding:"required"`
	Alias   string `json:"alias,omitempty"`
	Expires string `json:"expires,omitempty"`
}

type ShortenResponse struct {
	Code        string `json:"code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

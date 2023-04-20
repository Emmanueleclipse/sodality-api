package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User Model
type User struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username        string             `json:"username,omitempty" bson:"username,omitempty"`
	Email           string             `json:"email,omitempty" bson:"email,omitempty"`
	Password        string             `json:"password,omitempty" bson:"password,omitempty"`
	SubscriberCount int64              `json:"subscriber_count,omitempty" bson:"subscriber_count,omitempty"`
	// Role            int                `json:"role,omitempty" bson:"role,omitempty"`
	Bio    string `json:"bio,omitempty" bson:"bio,omitempty"`
	Avatar string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Dash   string `json:"dash,omitempty" bson:"dash,omitempty"`
}

type Followers struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatorID string             `json:"creator_id,omitempty" bson:"creator_id"`
	UserID    string             `json:"user_id,omitempty" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Content Model
type Content struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContentID    string             `json:"content_id,omitempty" bson:"content_id,omitempty"`
	UserID       string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CategoryName string             `json:"category_name,omitempty" bson:"category_name"`
	Title        string             `json:"title,omitempty" bson:"title,omitempty"`
	Body         string             `json:"body,omitempty" bson:"body,omitempty"`
	Type         string             `json:"type,omitempty" bson:"type,omitempty"`
	VideoURL     string             `json:"video,omitempty" bson:"video,omitempty"`
	AudioURL     string             `json:"audio,omitempty" bson:"audio,omitempty"`
	Fund         float64            `json:"fund,omitempty" bson:"fund,omitempty"`
	CurrencyType string             `json:"currency_type,omitempty" bson:"currency_type"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type GetAllContentWithCreatorResp struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ContentID    string             `json:"content_id,omitempty" bson:"content_id,omitempty"`
	UserID       string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	User         User               `json:"user_detail,omitempty" bson:"user_detail,omitempty"`
	CategoryName string             `json:"category_name,omitempty" bson:"category_name"`
	Title        string             `json:"title,omitempty" bson:"title,omitempty"`
	Body         string             `json:"body,omitempty" bson:"body,omitempty"`
	Type         string             `json:"type,omitempty" bson:"type,omitempty"`
	Fund         float64            `json:"fund,omitempty" bson:"fund,omitempty"`
	CurrencyType string             `json:"currency_type,omitempty" bson:"currency_type"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type FollowersCount struct {
	Count int64 `json:"count,omitempty" bson:"count"`
}

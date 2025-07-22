package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewUser creates a new user entity
func NewUser(username, email string) *User {
	now := time.Now()
	return &User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

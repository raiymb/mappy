package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleGuest Role = "guest"
	RoleUser  Role = "user"
	RoleMod   Role = "moderator"
	RoleAdmin Role = "admin"
)

// User represents an authenticated account in MongoDB.
type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"     json:"id"`
	Email        string               `bson:"email"             json:"email"       validate:"required,email"`
	PasswordHash string               `bson:"pwd,omitempty"     json:"-"           validate:"required"` // never expose
	Name         string               `bson:"name,omitempty"    json:"name"        validate:"required,min=2"`
	Role         Role                 `bson:"role"              json:"role"        validate:"oneof=guest user moderator admin"`
	AvatarURL    string               `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty"`
	SavedQuizzes []primitive.ObjectID `bson:"saved_quizzes,omitempty" json:"savedQuizzes,omitempty"`
	ViewedMaps   []primitive.ObjectID `bson:"viewed_maps,omitempty"   json:"viewedMaps,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

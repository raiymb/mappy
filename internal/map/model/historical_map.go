package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HistoricalMap is the static background layer (image / SVG).
type HistoricalMap struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title"         json:"title"`
	Period    string             `bson:"period"        json:"period"`    // e.g. "XVIII century"
	ImageURL  string             `bson:"image_url"     json:"imageUrl"`
	Tags      []string           `bson:"tags"          json:"tags"`
	CreatedAt time.Time          `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at"    json:"updatedAt"`
}

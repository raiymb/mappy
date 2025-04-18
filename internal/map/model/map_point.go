package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// PointType enumerates what a marker represents.
type PointType string

const (
	TypeEvent  PointType = "event"
	TypeFigure PointType = "figure"
	TypePlace  PointType = "place"
)

// MapPoint is a single marker on a historical map layer.
type MapPoint struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type           PointType          `bson:"type"          json:"type"`   // event | figure | place
	Title          string             `bson:"title"         json:"title"`
	X              float64            `bson:"x"             json:"x"`      // 0‑100  % of width
	Y              float64            `bson:"y"             json:"y"`      // 0‑100  % of height
	StartYear      int                `bson:"start_year"    json:"startYear"`
	EndYear        int                `bson:"end_year"      json:"endYear"`
	LinkedEntityID primitive.ObjectID `bson:"linked_id"     json:"linkedEntityId"`
}

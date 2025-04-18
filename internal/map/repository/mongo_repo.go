package repository

import (
	"context"

	"github.com/raiymb/mappy/internal/map/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository API exposed to services.
type Repository interface {
	PointsByYear(ctx context.Context, year int) ([]model.MapPoint, error)
	InsertPoint(ctx context.Context, mp *model.MapPoint) error
	UpdatePoint(ctx context.Context, mp *model.MapPoint) error
	DeletePoint(ctx context.Context, id string) error
}

// mongoRepo implements Repository.
type mongoRepo struct {
	col *mongo.Collection
}

// NewMongo returns a repo bound to db.map_points.
func NewMongo(db *mongo.Database) Repository {
	return &mongoRepo{col: db.Collection("map_points")}
}

func (r *mongoRepo) PointsByYear(ctx context.Context, year int) ([]model.MapPoint, error) {
	filter := bson.M{
		"start_year": bson.M{"$lte": year},
		"end_year":   bson.M{"$gte": year},
	}
	opts := options.Find().SetProjection(bson.M{
		"linked_id": 0, // strip heavy field for list response
	})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var out []model.MapPoint
	if err = cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *mongoRepo) InsertPoint(ctx context.Context, mp *model.MapPoint) error {
	_, err := r.col.InsertOne(ctx, mp)
	return err
}

func (r *mongoRepo) UpdatePoint(ctx context.Context, mp *model.MapPoint) error {
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": mp.ID}, mp)
	return err
}

func (r *mongoRepo) DeletePoint(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

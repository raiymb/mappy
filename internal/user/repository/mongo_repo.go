package repository

import (
	"context"
	"time"

	"github.com/raiymb/mappy/internal/user/model"
	"github.com/raiymb/mappy/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository declares all persistence operations.
type UserRepository interface {
	EnsureIndexes(ctx context.Context) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
}

// mongoRepo implements UserRepository.
type mongoRepo struct {
	col *mongo.Collection
}

// NewMongo creates a repo bound to db.users.
func NewMongo(db *mongo.Database) UserRepository {
	return &mongoRepo{col: db.Collection("users")}
}

// EnsureIndexes creates a unique index on email.
func (r *mongoRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	return err
}

func (r *mongoRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &u, err
}

func (r *mongoRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var u model.User
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &u, err
}

func (r *mongoRepo) Create(ctx context.Context, u *model.User) error {
	now := time.Now().UTC()
	u.ID = primitive.NewObjectID()
	u.CreatedAt, u.UpdatedAt = now, now
	_, err := r.col.InsertOne(ctx, u)
	return err
}

func (r *mongoRepo) Update(ctx context.Context, u *model.User) error {
	u.UpdatedAt = time.Now().UTC()
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": u.ID}, u)
	return err
}

func (r *mongoRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

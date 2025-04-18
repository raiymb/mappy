package storage

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoOnce   sync.Once
	mongoClient *mongo.Client
)

// NewMongo returns a singleton *mongo.Client.
func NewMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	var err error
	mongoOnce.Do(func() {
		opt := options.Client().ApplyURI(uri).
			SetConnectTimeout(10 * time.Second)
		mongoClient, err = mongo.Connect(ctx, opt)
	})
	return mongoClient, err
}

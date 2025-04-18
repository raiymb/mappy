package storage

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// We expose exactly one *mongo.Client for the whole process.
var (
	mu         sync.Mutex
	mongoCli   *mongo.Client
	mongoError error
)

func NewMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if mongoCli != nil || mongoError != nil {
		return mongoCli, mongoError
	}

	opt := options.Client().ApplyURI(uri).
		SetConnectTimeout(10 * time.Second)

	mongoCli, mongoError = mongo.Connect(ctx, opt)
	return mongoCli, mongoError
}

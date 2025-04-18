package token

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Blacklist struct {
	rdb *redis.Client
}

func NewBlacklist(rdb *redis.Client) *Blacklist {
	return &Blacklist{rdb: rdb}
}

// Revoke pushes the token JTI into Redis for its remaining TTL.
func (b *Blacklist) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	return b.rdb.Set(ctx, key(jti), "revoked", ttl).Err()
}

// IsBlacklisted returns true if the JTI exists in Redis.
func (b *Blacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	exists, err := b.rdb.Exists(ctx, key(jti)).Result()
	return exists == 1, err
}

func key(jti string) string {
	return "jwt:bl:" + jti
}

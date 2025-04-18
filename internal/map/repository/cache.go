package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/raiymb/mappy/internal/map/model"
	"github.com/redis/go-redis/v9"
)

// CachedRepo decorates Repository with Redis cache‑aside.
type CachedRepo struct {
	inner Repository
	rdb   *redis.Client
	ttl   time.Duration
}

// NewCached wraps any Repository with Redis caching.
func NewCached(inner Repository, rdb *redis.Client, ttl time.Duration) Repository {
	return &CachedRepo{inner: inner, rdb: rdb, ttl: ttl}
}

func key(year int) string { return fmt.Sprintf("map:points:%d", year) }

func (c *CachedRepo) PointsByYear(ctx context.Context, year int) ([]model.MapPoint, error) {
	// 1. attempt cache
	if data, err := c.rdb.Get(ctx, key(year)).Bytes(); err == nil {
		var cached []model.MapPoint
		if json.Unmarshal(data, &cached) == nil {
			return cached, nil
		}
	}

	// 2. miss → load from Mongo
	points, err := c.inner.PointsByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	// 3. write‑back (fire‑and‑forget)
	if b, _ := json.Marshal(points); len(b) > 0 {
		_ = c.rdb.Set(ctx, key(year), b, c.ttl).Err()
	}
	return points, nil
}

/* passthrough mutating methods — invalidate cache */

func (c *CachedRepo) InsertPoint(ctx context.Context, mp *model.MapPoint) error {
	if err := c.inner.InsertPoint(ctx, mp); err != nil {
		return err
	}
	c.flushByYear(ctx, mp.StartYear, mp.EndYear)
	return nil
}

func (c *CachedRepo) UpdatePoint(ctx context.Context, mp *model.MapPoint) error {
	if err := c.inner.UpdatePoint(ctx, mp); err != nil {
		return err
	}
	c.flushByYear(ctx, mp.StartYear, mp.EndYear)
	return nil
}

func (c *CachedRepo) DeletePoint(ctx context.Context, id string) error {
	// lazy: delete all cached years (small key‑space)
	_ = c.rdb.FlushDB(ctx).Err()
	return c.inner.DeletePoint(ctx, id)
}

func (c *CachedRepo) flushByYear(ctx context.Context, from, to int) {
	for y := from; y <= to; y++ {
		_ = c.rdb.Del(ctx, key(y)).Err()
	}
}

package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheManager struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func NewRedisCacheManager(addr string, password string, db int) *RedisCacheManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCacheManager{
		client:     rdb,
		defaultTTL: DefaultTTL,
	}
}

func (r *RedisCacheManager) Get(ctx context.Context, key CacheKey) (string, error) {
	val, err := r.client.Get(ctx, key.Key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisCacheManager) Set(ctx context.Context, key CacheKey, value string) error {
	ttl := r.defaultTTL
	if key.Time > 0 {
		ttl = time.Duration(key.Time) * time.Minute
	}

	return r.client.Set(ctx, key.Key, value, ttl).Err()
}

func (r *RedisCacheManager) Remove(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCacheManager) RemoveByPrefix(ctx context.Context, prefix string) error {
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (r *RedisCacheManager) Clear(ctx context.Context) error {
	// WARNING: This flushes the entire Redis DB!
	return r.client.FlushDB(ctx).Err()
}

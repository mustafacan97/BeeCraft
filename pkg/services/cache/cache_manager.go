package cache

import "context"

type CacheManager interface {
	Get(ctx context.Context, key CacheKey) (string, error)
	Set(ctx context.Context, key CacheKey, value string) error
	Remove(ctx context.Context, key CacheKey) error
	RemoveByPrefix(ctx context.Context, prefix string) error
	Clear(ctx context.Context) error
}

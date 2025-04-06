package cache

import "context"

type CacheManager interface {
	Get(ctx context.Context, key cacheKey) (string, error)
	Set(ctx context.Context, key cacheKey, value string) error
	Remove(ctx context.Context, key cacheKey) error
	RemoveByPrefix(ctx context.Context, prefix string) error
	Clear(ctx context.Context) error
}

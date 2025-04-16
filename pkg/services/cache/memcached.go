package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheManager struct {
	client     *memcache.Client
	defaultTTL time.Duration
}

func NewMemcacheManager(server string) *MemcacheManager {
	return &MemcacheManager{
		client:     memcache.New(server),
		defaultTTL: defaultTTL,
	}
}

func (m *MemcacheManager) Get(ctx context.Context, key cacheKey) (string, error) {
	item, err := m.client.Get(key.Key)
	if err == memcache.ErrCacheMiss {
		return "", ErrKeyNotFound
	} else if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func (m *MemcacheManager) Set(ctx context.Context, key cacheKey, value string) error {
	ttl := m.defaultTTL
	if key.CacheTime > 0 {
		ttl = time.Duration(key.CacheTime) * time.Minute
	}

	item := &memcache.Item{
		Key:        key.Key,
		Value:      []byte(value),
		Expiration: int32(ttl.Seconds()), // TTL in seconds
	}

	return m.client.Set(item)
}

func (m *MemcacheManager) Remove(ctx context.Context, key cacheKey) error {
	return m.client.Delete(key.Key)
}

// Memcached does not support key scanning or prefix deletes natively
func (m *MemcacheManager) RemoveByPrefix(ctx context.Context, prefix string) error {
	return fmt.Errorf("RemoveByPrefix is not supported in memcached")
}

func (m *MemcacheManager) Clear(ctx context.Context) error {
	return m.client.FlushAll()
}

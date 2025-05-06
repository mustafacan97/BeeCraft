package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCacheManager_SetAndGet(t *testing.T) {
	ctx := context.Background()
	cache := NewInMemoryCacheManager()

	key := CacheKey{
		Key:  "test-key",
		Time: 1, // 1 minutes
	}
	value := "test-value"

	err := cache.Set(ctx, key, value)
	assert.NoError(t, err)

	got, err := cache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestInMemoryCacheManager_Expiration(t *testing.T) {
	ctx := context.Background()
	cache := NewInMemoryCacheManager()

	key := CacheKey{
		Key:  "short-expire",
		Time: 0, // default TTL
	}
	value := "short-lived"

	// set TTL short in order to test
	cache.defaultTTL = 1 * time.Second
	_ = cache.Set(ctx, key, value)

	time.Sleep(2 * time.Second)

	got, err := cache.Get(ctx, key)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestInMemoryCacheManager_Remove(t *testing.T) {
	ctx := context.Background()
	cache := NewInMemoryCacheManager()

	key := CacheKey{
		Key:  "remove-key",
		Time: 5,
	}
	value := "to-be-removed"

	_ = cache.Set(ctx, key, value)
	_ = cache.Remove(ctx, key)

	got, err := cache.Get(ctx, key)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestInMemoryCacheManager_RemoveByPrefix(t *testing.T) {
	ctx := context.Background()
	cache := NewInMemoryCacheManager()

	_ = cache.Set(ctx, CacheKey{"prefix:1", 5}, "val1")
	_ = cache.Set(ctx, CacheKey{"prefix:2", 5}, "val2")
	_ = cache.Set(ctx, CacheKey{"other:1", 5}, "val3")

	err := cache.RemoveByPrefix(ctx, "prefix:")
	assert.NoError(t, err)

	_, err1 := cache.Get(ctx, CacheKey{"prefix:1", 0})
	_, err2 := cache.Get(ctx, CacheKey{"prefix:2", 0})
	_, err3 := cache.Get(ctx, CacheKey{"other:1", 0})

	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.NoError(t, err3)
}

func TestInMemoryCacheManager_Clear(t *testing.T) {
	ctx := context.Background()
	cache := NewInMemoryCacheManager()

	_ = cache.Set(ctx, CacheKey{"key1", 5}, "val1")
	_ = cache.Set(ctx, CacheKey{"key2", 5}, "val2")

	err := cache.Clear(ctx)
	assert.NoError(t, err)

	_, err1 := cache.Get(ctx, CacheKey{"key1", 0})
	_, err2 := cache.Get(ctx, CacheKey{"key2", 0})

	assert.Error(t, err1)
	assert.Error(t, err2)
}

package cache

import (
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

const (
	maxCacheTtlMinute = 1440
	defaultTTL        = time.Duration(10) * time.Minute
)

type cacheKey struct {
	Key       string
	CacheTime int // TTL in minutes
}

func (ck cacheKey) NewCacheKey(key string, cacheTime int) *cacheKey {
	if cacheTime > maxCacheTtlMinute {
		cacheTime = maxCacheTtlMinute
	}

	return &cacheKey{
		Key:       key,
		CacheTime: cacheTime,
	}
}

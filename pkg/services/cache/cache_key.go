package cache

import (
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

const (
	MaxCacheTtlMinute = time.Duration(1440) * time.Minute
	DefaultTTL        = time.Duration(10) * time.Minute
)

type CacheKey struct {
	Key  string
	Time time.Duration
}

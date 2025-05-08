package cache

import (
	"context"
	"strings"
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type InMemoryCacheManager struct {
	data       map[string]entry
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

const cleanupInterval = 15 * time.Minute

func NewInMemoryCacheManager() *InMemoryCacheManager {
	manager := &InMemoryCacheManager{
		data:       make(map[string]entry),
		defaultTTL: DefaultTTL,
	}

	go manager.startCleanupTask()

	return manager
}

func (c *InMemoryCacheManager) Get(ctx context.Context, key CacheKey) (any, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.data[key.Key]
	if !ok || time.Now().After(item.expiresAt) {
		return nil, ErrKeyNotFound
	}

	return item.value, nil
}

func (c *InMemoryCacheManager) Set(ctx context.Context, key CacheKey, value any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ttl := c.defaultTTL
	if key.Time > 0 {
		ttl = time.Duration(key.Time) * time.Minute
	}

	c.data[key.Key] = entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (c *InMemoryCacheManager) Remove(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	return nil
}

func (c *InMemoryCacheManager) RemoveByPrefix(ctx context.Context, prefix string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for k := range c.data {
		if strings.HasPrefix(k, prefix) {
			delete(c.data, k)
		}
	}

	return nil
}

func (c *InMemoryCacheManager) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]entry)
	return nil
}

func (c *InMemoryCacheManager) startCleanupTask() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.cleanupExpiredEntries()
	}
}

func (c *InMemoryCacheManager) cleanupExpiredEntries() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for k, v := range c.data {
		if now.After(v.expiresAt) {
			delete(c.data, k)
		}
	}
}

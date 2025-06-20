package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	Value      any
	Enable     bool
	ExpiryTime time.Time
}

type cacheQueue struct {
	m  sync.Map
	mu sync.Mutex
}

var cache1 = newCache().Expired(10 * time.Second)

func newCache() *cacheQueue {
	return &cacheQueue{}
}

func (c *cacheQueue) Expired(interval time.Duration) *cacheQueue {
	go func() {
		for {
			time.Sleep(interval)
			c.cleanupExpired()
		}
	}()

	return c
}

func (c *cacheQueue) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m.Range(func(key, value any) bool {
		cacheItemValue := value.(cacheItem)
		if cacheItemValue.Enable {
			if time.Now().After(cacheItemValue.ExpiryTime) {
				c.m.Delete(key)
			}
		}
		return true
	})
}

func Store(key any, value any) {
	cache1.mu.Lock()
	defer cache1.mu.Unlock()

	cache1.m.Store(key, cacheItem{Value: value, Enable: false})
}

func Set(key any, value any, duration time.Duration) {
	cache1.mu.Lock()
	defer cache1.mu.Unlock()

	expiryTime := time.Now().Add(duration)
	cache1.m.Store(key, cacheItem{Value: value, Enable: true, ExpiryTime: expiryTime})
}

func Load(key any) (any, bool) {
	cache1.mu.Lock()
	defer cache1.mu.Unlock()

	item, ok := cache1.m.Load(key)
	if !ok {
		return nil, false
	}

	cacheItemValue := item.(cacheItem)
	if cacheItemValue.Enable {
		if time.Now().After(cacheItemValue.ExpiryTime) {
			cache1.m.Delete(key)
			return nil, false
		}
	}

	return cacheItemValue.Value, true
}

func IsExist(key string) bool {
	cache1.mu.Lock()
	defer cache1.mu.Unlock()

	if _, ok := cache1.m.Load(key); ok {
		return ok
	}

	return false
}

func Delete(key string) {
	cache1.mu.Lock()
	defer cache1.mu.Unlock()

	cache1.m.Delete(key)
}

func Clear() {
	cache1.mu.Lock()
	defer cache1.mu.Unlock()

	cache1.m = sync.Map{}
}

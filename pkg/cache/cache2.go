package cache

import (
	"sync"
)

// https://github.com/patrickmn/go-cache
// https://github.com/allegro/bigcache
// https://github.com/dgraph-io/ristretto

type cache2 struct {
	mu    sync.RWMutex
	store map[string]int
}

var Cache = cache2{
	store: make(map[string]int),
}

func (c *cache2) Set(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *cache2) Get(key string) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.store[key]
	return value, ok
}

func (c *cache2) List() map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.store
}

func (c *cache2) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *cache2) IsExist(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.store[key]; ok {
		return ok
	}
	return false
}

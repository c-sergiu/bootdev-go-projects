package cache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu *sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt:time.Now(),
		val:val,
	}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if ok { return e.val, true
	} else { return nil, false }
}

func (c *Cache) reapLoop(interval time.Duration) {
	tic := time.NewTicker(interval)
	defer tic.Stop()
	for range tic.C {
		c.clear(interval)
	}
}

func (c *Cache) clear(interval time.Duration) {
	for k, e := range c.entries {
		c.mu.Lock()
		if time.Since(e.createdAt) >= interval {
			delete(c.entries, k)
		}
		c.mu.Unlock()
	}
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache{
	c := Cache{
		entries: make(map[string]cacheEntry),
		mu: &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return &c
}

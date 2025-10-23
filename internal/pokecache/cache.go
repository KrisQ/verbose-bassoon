// Package pokecache - handles cache
package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entries:  make(map[string]cacheEntry),
		mu:       sync.Mutex{},
		interval: interval,
	}
	c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			c.mu.Lock()
			fiveSecondsAgo := time.Now().Add(-5 * time.Second)
			for key, entry := range c.entries {
				if entry.createdAt.Before(fiveSecondsAgo) {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value with an expiry timestamp.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is a simple in-memory TTL cache safe for concurrent use.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	defaultTTL time.Duration
}

// New creates a Cache with the given default TTL.
func New(defaultTTL time.Duration) *Cache {
	return &Cache{
		items:      make(map[string]Entry),
		defaultTTL: defaultTTL,
	}
}

// Set stores a value under key with the default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value under key with an explicit TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value by key. Returns (value, true) on hit, (nil, false) on miss or expiry.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all expired entries.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.items {
		if now.After(v.ExpiresAt) {
			delete(c.items, k)
		}
	}
}

// Len returns the number of items currently stored (including expired ones not yet flushed).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

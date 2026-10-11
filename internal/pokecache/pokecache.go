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
	interval time.Duration
	mu       sync.Mutex
	entries  map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{interval, sync.Mutex{}, map[string]cacheEntry{}}

	go func() {
		c.reapLoop()
	}()

	return c
}

func (s *Cache) Add(key string, val []byte) {
	s.mu.Lock()
	s.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	s.mu.Unlock()
}

func (s *Cache) Get(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

func (s *Cache) reapLoop() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			s.reap()
		}
	}
}

func (s *Cache) reap() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for i, c := range s.entries {
		expiryAt := c.createdAt.Add(s.interval)
		if expiryAt.Before(now) {
			delete(s.entries, i)
		}
	}
}

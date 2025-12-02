package auth

import (
	"context"
	"sync"
	"time"
)

// TokenCache определяет интерфейс для кэширования blacklist токенов
type TokenCache interface {
	IsBlacklisted(ctx context.Context, tokenHash string) bool
	AddToBlacklist(ctx context.Context, tokenHash string, ttl time.Duration)
	RemoveFromBlacklist(ctx context.Context, tokenHash string)
	Clear(ctx context.Context)
}

type memoryTokenCache struct {
	mu      sync.RWMutex
	entries map[string]time.Time
}

func NewMemoryTokenCache() TokenCache {
	cache := &memoryTokenCache{
		entries: make(map[string]time.Time),
	}

	go cache.cleanup()

	return cache
}

func (c *memoryTokenCache) IsBlacklisted(ctx context.Context, tokenHash string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expiresAt, exists := c.entries[tokenHash]
	if !exists {
		return false
	}

	if time.Now().After(expiresAt) {
		return false
	}

	return true
}

func (c *memoryTokenCache) AddToBlacklist(ctx context.Context, tokenHash string, ttl time.Duration) {
	if ttl <= 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[tokenHash] = time.Now().Add(ttl)
}

func (c *memoryTokenCache) RemoveFromBlacklist(ctx context.Context, tokenHash string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, tokenHash)
}

func (c *memoryTokenCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]time.Time)
}

func (c *memoryTokenCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for tokenHash, expiresAt := range c.entries {
			if now.After(expiresAt) {
				delete(c.entries, tokenHash)
			}
		}
		c.mu.Unlock()
	}
}

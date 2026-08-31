package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService defines the interface for interacting with redis cache.
type CacheService interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// CacheStore provides a higher-level abstraction for interacting with the cache.
// It uses a CacheService implementation to perform the actual cache operations.
type CacheStore struct {
	storer CacheService
}

// NewCacheStore creates a new CacheStore instance.
// It takes a CacheService implementation as an argument, which will be used for cache operations.
// This allows for flexible cache service implementations to be passed in.
func NewCacheStore(cache CacheService) *CacheStore {
	return &CacheStore{
		storer: cache,
	}
}

// Get retrieves the value associated with the given key from the cache.
// It returns a StringCmd which can be used to obtain the actual value.
func (c CacheStore) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.storer.Get(ctx, key)
}

// Set stores the given value in the cache under the specified key, with an expiration time.
// It returns a StatusCmd which indicates the result of the operation.
func (c CacheStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.storer.Set(ctx, key, value, expiration)
}

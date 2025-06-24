package interfaces

import (
	"context"
	"time"
)

// CacheService defines the interface for caching operations
// This will be mocked using: mockery --name=CacheService --dir=./internal/interfaces --output=./mocks
type CacheService interface {
	// Set stores a value in the cache with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (interface{}, error)
	
	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)
}
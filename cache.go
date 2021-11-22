package gocache

import (
	"time"
)

type cacheConfig struct {
	size      int // size in MB
	ttl       time.Duration
	cleanFreq time.Duration
}

type Option interface {
	apply(config *cacheConfig)
}

// define optionFunc implement Option interface
type optionFunc func(config *cacheConfig)

func (opt optionFunc) apply(config *cacheConfig) {
	opt(config)
}

func WithSizeInMB(size int) Option {
	return optionFunc(func(config *cacheConfig) {
		config.size = size
	})
}

func WithTTL(ttl time.Duration) Option {
	return optionFunc(func(config *cacheConfig) {
		config.ttl = ttl
	})
}
func WithCleanFrequency(cleanFreq time.Duration) Option {
	return optionFunc(func(config *cacheConfig) {
		config.cleanFreq = cleanFreq
	})
}

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

func New(options ...Option) (Cache, error) {
	config := &cacheConfig{
		size:      512,
		ttl:       1 * time.Hour,
		cleanFreq: 5 * time.Minute,
	}
	for _, option := range options {
		option.apply(config)
	}
	return newBigCache(config)
}

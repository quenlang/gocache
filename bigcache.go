package gocache

import (
	"bytes"
	"encoding/gob"
	"github.com/allegro/bigcache/v3"
)

type bigCache struct {
	cache *bigcache.BigCache
}

func newBigCache(config *cacheConfig) (*bigCache, error) {
	cache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             1024,
		LifeWindow:         config.ttl,
		CleanWindow:        config.cleanFreq,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       128,
		Verbose:            false,
		HardMaxCacheSize:   config.size,
		StatsEnabled:       true,
	})
	return &bigCache{
		cache: cache,
	}, err
}

func (c *bigCache) Set(key string, value interface{}) error {
	b, err := serializeGOB(value)
	if err != nil {
		return err
	}
	return c.cache.Set(key, b)
}

func (c *bigCache) Get(key string) (interface{}, error) {
	b, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	return deserializeGOB(b)
}

func (c *bigCache) Delete(key string) error {
	return c.cache.Delete(key)
}

func (c *bigCache) Capacity() int {
	return c.cache.Capacity()
}

func (c *bigCache) Len() int {
	return c.cache.Len()
}

func (c *bigCache) Exist(key string) bool {
	v, err := c.cache.Get(key)
	return err == nil && v != nil
}

func serializeGOB(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)
	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserializeGOB(valueBytes []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

package cache

import (
	"sync"
	"sync/atomic"
)

const MAX_CACHE_SIZE = 10000

type Map[K comparable, V any] struct {
	count atomic.Uint64
	m     sync.Map
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
	m.count.Add(1)
}

func (m *Map[K, V]) Count() uint64 {
	return m.count.Load()
}

type CacheEntry struct {
	Url   string
	Album string
}

type ImageUrlCache struct {
	cache *Map[string, CacheEntry]
}

var imageUrlCache = &ImageUrlCache{cache: new(Map[string, CacheEntry])}

func GetImageUrlCache() *ImageUrlCache {
	return imageUrlCache
}

func (c *ImageUrlCache) Get(key string) (CacheEntry, bool) {
	value, ok := c.cache.Load(key)
	return value, ok
}

func (c *ImageUrlCache) Set(key string, value CacheEntry) {
	if c.cache.Count() > MAX_CACHE_SIZE {
		c.cache = new(Map[string, CacheEntry])
	}
	c.cache.Store(key, value)
}

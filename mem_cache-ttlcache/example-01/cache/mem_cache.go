package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemCache[K comparable, V any] struct {
	cache *ttlcache.Cache[K, V]
}

type IMemCache[K comparable, V any] interface {
	// Set [key, value] with default TTL, it means defaultTTL when init MemCache with effect this [key, value]
	Set(key K, value V)
	// Set [key, value] with TTL, it means defaultTTL when init MemCache with no effect this [key, value], instead will be TTL param in this function
	SetTTL(key K, value V, ttl time.Duration)
	// Get value by key and refresh this TTL of this [key, value]
	Get(key K) (V, bool)
	// Delete [key, value] by key
	Del(key K)
}

func NewMemCache[K comparable, V any](defaultTTL time.Duration) IMemCache[K, V] {
	cache := ttlcache.New(
		ttlcache.WithTTL[K, V](defaultTTL),
	)

	go cache.Start()
	return &MemCache[K, V]{
		cache: cache,
	}
}

func (mcache *MemCache[K, V]) Set(key K, value V) {
	mcache.cache.Set(key, value, ttlcache.DefaultTTL)
}

func (mcache *MemCache[K, V]) SetTTL(key K, value V, ttl time.Duration) {
	mcache.cache.Set(key, value, ttl)
}

func (mcache *MemCache[K, V]) Get(key K) (V, bool) {
	item := mcache.cache.Get(key)
	if item == nil {
		var zero V
		return zero, false
	}
	return item.Value(), true
}

func (mcache *MemCache[K, V]) Del(key K) {
	mcache.cache.Delete(key)
}

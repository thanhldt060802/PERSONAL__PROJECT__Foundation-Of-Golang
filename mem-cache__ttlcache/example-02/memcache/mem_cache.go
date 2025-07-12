package memcache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	log "github.com/sirupsen/logrus"
)

type MemCache[K comparable, V any] struct {
	cache *ttlcache.Cache[K, V]
}

type IMemCache[K comparable, V any] interface {
	Set(key K, value V)
	SetTTL(key K, value V, ttl time.Duration)
	Get(key K) (V, bool)
	Del(key K)
}

func NewMemCache[K comparable, V any](defaultTTL time.Duration) IMemCache[K, V] {
	cache := ttlcache.New(
		ttlcache.WithTTL[K, V](defaultTTL),
	)
	cache.OnInsertion(func(ctx context.Context, item *ttlcache.Item[K, V]) {
		log.Infof("memcache: inserted [%v, %v], expires at %v", item.Key(), item.Value(), item.ExpiresAt())
	})
	cache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[K, V]) {
		switch reason {
		case ttlcache.EvictionReasonDeleted:
			log.Infof("memcache-deleted: removed [%v, %v]", item.Key(), item.Value())
		case ttlcache.EvictionReasonCapacityReached:
			log.Infof("memcache-capacity-reached: removed [%v, %v]", item.Key(), item.Value())
		case ttlcache.EvictionReasonExpired:
			log.Infof("memcache-expired: removed [%v, %v]", item.Key(), item.Value())
		case ttlcache.EvictionReasonMaxCostExceeded:
			log.Infof("memcache-max-cost-exceeded: removed [%v, %v]", item.Key(), item.Value())
		}
	})

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

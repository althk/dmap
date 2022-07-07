package dmap

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

type Shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

type DMap[K comparable, V any] []*Shard[K, V]

func New[K comparable, V any](nShards int) DMap[K, V] {
	shards := make([]*Shard[K, V], nShards)
	for i := 0; i < nShards; i++ {
		shard := &Shard[K, V]{
			items: make(map[K]V),
		}
		shards[i] = shard
	}
	return shards
}

func (m DMap[K, V]) getShardIndex(key K) int {
	checksum := sha1.Sum([]byte(fmt.Sprintf("%v", key)))
	hash := int(checksum[7]<<1 | checksum[19])
	return hash % len(m)
}

func (m DMap[K, V]) getShard(key K) *Shard[K, V] {
	i := m.getShardIndex(key)
	return m[i]
}

func (m DMap[K, V]) Get(key K) (V, bool) {
	shard := m.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	v, e := shard.items[key]
	return v, e
}

func (m DMap[K, V]) Set(key K, val V) {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = val
}

func (m DMap[K, V]) Keys() []K {
	keys := make([]K, 0)

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	mu := sync.Mutex{}

	for _, shard := range m {
		go func(shard *Shard[K, V]) {
			shard.mu.RLock()
			defer shard.mu.RUnlock()

			for key := range shard.items {
				mu.Lock()
				keys = append(keys, key)
				mu.Unlock()
			}
			wg.Done()
		}(shard)
	}
	wg.Wait()
	return keys
}

func (m DMap[K, V]) Remove(key K) {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

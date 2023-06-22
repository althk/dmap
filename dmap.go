// Package dmap provides a thread-safe, generics-based horizontally distributed (sharded) map.
package dmap

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

// Shard represents one partition of the entire data.
type Shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
	count int
}

// DMap represents a simple map structure which shards
// its data for improved performance.
// The number of shards (partitions) is fixed, and is set
// on construction of the map.
// DMap supports heterogeneous values (when V is interface{}).
// DMap is thread-safe.
type DMap[K comparable, V any] []*Shard[K, V]

// New creates a new DMap with nShards number of shards.
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

// Get returns the value for the given key from the map.
// If a key is not found, ok is false.
func (m DMap[K, V]) Get(key K) (V, bool) {
	shard := m.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	v, ok := shard.items[key]
	return v, ok
}

// Set sets the given key, value in the map.
func (m DMap[K, V]) Set(key K, val V) {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = val
	shard.count += 1
}

// Keys returns a list of all keys in the map (from all shards).
func (m DMap[K, V]) Keys() []K {
	keys := make([]K, 0)

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	mu := sync.Mutex{}

	for _, shard := range m {
		go func(shard *Shard[K, V]) {
			shard.mu.RLock()
			defer shard.mu.RUnlock()

			mu.Lock()
			for key := range shard.items {
				keys = append(keys, key)
			}
			mu.Unlock()
			wg.Done()
		}(shard)
	}
	wg.Wait()
	return keys
}

// Remove deletes the key from the map (if found).
func (m DMap[K, V]) Remove(key K) {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

// Count returns the total number of items in the map (across all shards).
func (m DMap[K, V]) Count() int64 {
	count := 0
	for i := 0; i < len(m); i++ {
		shard := m[i]
		shard.mu.RLock()
		count += shard.count
		shard.mu.RUnlock()
	}
	return int64(count)
}

func (m DMap[K, V]) Has(key K) bool {
	_, ok := m.Get(key)
	return ok
}

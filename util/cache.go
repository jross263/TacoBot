package util

import (
	"fmt"
	"maps"
	"sync"
)

type Cache[K comparable, V any] struct {
	mu   sync.RWMutex
	name string
	data map[K]V
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.data[k]
	return v, ok
}

func (c *Cache[K, V]) AddOnce(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.data[k]
	if ok {
		panic(fmt.Sprintf("entry with key %v is already present in %s cache", k, c.name))
	}

	c.data[k] = v
}

func (c *Cache[K, V]) Set(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[k] = v
}

func (c *Cache[K, V]) Update(k K, fn func(V) V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.data[k]
	if !ok {
		return fmt.Errorf("key %v not found in %s cache", k, c.name)
	}

	c.data[k] = fn(v)
	return nil
}

// NOTE: These are shallow copies
func (c *Cache[K, V]) All() map[K]V {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return maps.Clone(c.data)
}

func NewCache[K comparable, V any](name string) *Cache[K, V] {
	return &Cache[K, V]{name: name, data: make(map[K]V)}
}

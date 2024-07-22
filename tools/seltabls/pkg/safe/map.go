package safe

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

// Map is a thread-safe map.
type Map[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSafeMap creates a new SafeMap.
func NewSafeMap[K comparable, V any]() *Map[K, V] {
	log.Debugf("safe.NewSafeMap called")
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

// Get returns the value for the given key.
func (sm *Map[K, V]) Get(key K) (V, bool) {
	log.Debugf("safe.Map.Get called with key: %v", key)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.m[key]
	return val, ok
}

// Set sets the value for the given key.
func (sm *Map[K, V]) Set(key K, value V) {
	log.Debugf("safe.Map.Set called with key: %v", key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

// Delete deletes the value for the given key.
func (sm *Map[K, V]) Delete(key K) {
	log.Debugf("safe.Map.Delete called with key: %v", key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

// Len returns the length of the map.
func (sm *Map[K, V]) Len() int {
	log.Debugf("safe.Map.Len called")
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.m)
}

// Clear clears the map.
func (sm *Map[K, V]) Clear() {
	log.Debugf("safe.Map.Clear called")
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m = make(map[K]V)
}

// String returns a string representation of the map.
func (sm *Map[K, V]) String() string {
	log.Debugf("safe.Map.String called")
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var b strings.Builder
	for k, v := range sm.m {
		b.WriteString(fmt.Sprintf("%v: %v\n", k, v))
	}
	return b.String()
}

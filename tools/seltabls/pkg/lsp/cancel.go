package lsp

import (
	"context"
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	// CancelMap is a map of cancel functions
	CancelMap = make(CCancelMap)
	mu        sync.Mutex
)

// CCancelMap is a map of cancel functions that is thread safe.
type CCancelMap map[int]context.CancelFunc

// Cancel cancels the cancel function for the given id.
func (c CCancelMap) Cancel(id int) {
	if cancel, ok := c[id]; ok {
		cancel()
	}
}

// Add adds a cancel function for the given id.
func (c CCancelMap) Add(id int, cancel context.CancelFunc) {
	mu.Lock()
	defer mu.Unlock()
	c[id] = cancel
	log.Debugf("added cancel function for id: %d", id)
}

// Remove removes the cancel function for the given id.
func (c CCancelMap) Remove(id int) {
	mu.Lock()
	defer mu.Unlock()
	delete(c, id)
	log.Debugf("removed cancel function for id: %d", id)
}

// Contains checks if the cancel map contains the given id.
func (c CCancelMap) Contains(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := c[id]
	return ok
}

// Len returns the length of the cancel map.
func (c CCancelMap) Len() int {
	mu.Lock()
	defer mu.Unlock()
	return len(c)
}

// Clear clears the cancel map.
func (c CCancelMap) Clear() {
	mu.Lock()
	defer mu.Unlock()
	for k := range c {
		delete(c, k)
	}
}

// String returns a string representation of the cancel map.
func (c CCancelMap) String() string {
	str := ""
	for k, v := range c {
		str += fmt.Sprintf("%d: %v\n", k, v)
	}
	return str
}

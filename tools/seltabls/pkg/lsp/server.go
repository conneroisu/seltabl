package lsp

import (
	"sync"
)

// LineRange returns a range of a line in a document
//
// line is the line number
//
// start is the start character of the range
//
// end is the end character of the range
func LineRange(line, start, end int) Range {
	return Range{
		Start: Position{
			Line:      line,
			Character: start,
		},
		End: Position{
			Line:      line,
			Character: end,
		},
	}
}

// HandlerFunc is a function that handles a request by returning a function
type HandlerFunc func(ResponseWriter, *Request)

// Handler is a function that handles a request
type Handler interface {
	ServeRPC(ResponseWriter, *Request)
}

// ServeRPC serves a request
func (f HandlerFunc) ServeRPC(w ResponseWriter, r *Request) {
	f(w, r)
}

// ResponseWriter is a writer for a response
type ResponseWriter interface {
	Write(data []byte) (int, error)
}

// Route is a interface for routing requests and notifications to a handler
type Route interface {
	Handler(method string, request interface{}) error
}

// Router is a struct for routing requests and notifications to handlers
type Router struct {
	handlers map[string]Handler
}

// ServeMux is a multiplexer that can be used to serve rpc requests
type ServeMux struct {
	mu   sync.RWMutex
	tree routingNode
}

// Handle registers the handler for the given pattern.
func (s *ServeMux) Handle(method string, handler Handler) {
	if s.tree.children.m == nil {
		s.tree.children.m = map[string]*routingNode{}
	}
	s.tree.children.m[method] = &routingNode{
		handler: handler,
	}
}

// routingNode is a node in the routing tree.
type routingNode struct {
	children mapping[string, *routingNode]
	handler  Handler
}

// DefaultMux is the default [ServeMux] used by [Serve].
var DefaultMux = &defaultMux
var defaultMux ServeMux

// A mapping is a collection of key-value pairs where the keys are unique.
// A zero mapping is empty and ready to use.
// A mapping tries to pick a representation that makes [mapping.find] most efficient.
type mapping[K comparable, V any] struct {
	s []entry[K, V] // for few pairs
	m map[K]V       // for many pairs
}

// entry is an entry in a mapping.
type entry[K comparable, V any] struct {
	key   K
	value V
}

// maxSlice is the maximum number of pairs for which a slice is used.
// It is a variable for benchmarking.
var maxSlice = 8

// add adds a key-value pair to the mapping.
func (h *mapping[K, V]) add(k K, v V) {
	if h.m == nil && len(h.s) < maxSlice {
		h.s = append(h.s, entry[K, V]{k, v})
	} else {
		if h.m == nil {
			h.m = map[K]V{}
			for _, e := range h.s {
				h.m[e.key] = e.value
			}
			h.s = nil
		}
		h.m[k] = v
	}
}

// find returns the value corresponding to the given key.
// The second return value is false if there is no value
// with that key.
func (h *mapping[K, V]) find(k K) (v V, found bool) {
	if h == nil {
		return v, false
	}
	if h.m != nil {
		v, found = h.m[k]
		return v, found
	}
	for _, e := range h.s {
		if e.key == k {
			return e.value, true
		}
	}
	return v, false
}

// eachPair calls f for each pair in the mapping.
// If f returns false, pairs returns immediately.
func (h *mapping[K, V]) eachPair(f func(k K, v V) bool) {
	if h == nil {
		return
	}
	if h.m != nil {
		for k, v := range h.m {
			if !f(k, v) {
				return
			}
		}
	} else {
		for _, e := range h.s {
			if !f(e.key, e.value) {
				return
			}
		}
	}
}

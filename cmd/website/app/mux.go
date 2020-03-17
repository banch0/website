package app

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// map["GET"] - map["/"] - handler GET
// map["POST"] - map["/"] - handler POST

// ExactMux new router
type ExactMux struct {
	mutex           sync.RWMutex
	routes          map[string]map[string]ExactMuxEntry
	routesSorted    map[string][]ExactMuxEntry
	notFoundHandler http.Handler
}

// NewExactMux ...
func NewExactMux() *ExactMux {
	return &ExactMux{}
}

// ServerHTTP interface implementation ...
func (m *ExactMux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handler, err := m.handler(request.Method, request.URL.Path); err == nil {
		handler.ServeHTTP(writer, request)
	}

	if m.notFoundHandler != nil {
		m.notFoundHandler.ServeHTTP(writer, request)
	}
}

// GET method
func (m *ExactMux) GET(pattern string, handlerFunc func(responseWriter http.ResponseWriter, request *http.Request)) {
	m.HandleFunc(http.MethodGet, pattern, handlerFunc)
}

// POST method
func (m *ExactMux) POST(pattern string, handlerFunc func(responseWriter http.ResponseWriter, request *http.Request)) {
	m.HandleFunc(http.MethodPost, pattern, handlerFunc)
}

// HandleFunc ...
func (m *ExactMux) HandleFunc(method string, pattern string, handlerFunc func(responseWriter http.ResponseWriter, request *http.Request)) {
	// pattern: "/..."
	if !strings.HasPrefix(pattern, "/") {
		panic(fmt.Errorf("pattern must start with /: %s", pattern))
	}

	if handlerFunc == nil { // ?
		panic(errors.New("handler can't be empty"))
	}

	// TODO: check method
	m.mutex.Lock()
	defer m.mutex.Unlock()
	entry := ExactMuxEntry{
		pattern: pattern,
		handler: http.HandlerFunc(handlerFunc),
		weight:  calculateWeight(pattern),
	}

	// запретить добавлять дубликаты
	if _, exists := m.routes[method][pattern]; exists {
		panic(fmt.Errorf("ambigious mapping: %s", pattern))
	}

	if m.routes == nil {
		m.routes = make(map[string]map[string]ExactMuxEntry)
	}

	if m.routes[method] == nil {
		m.routes[method] = make(map[string]ExactMuxEntry)
	}

	m.routes[method][pattern] = entry
	m.appendSorted(method, entry)
}

func (m *ExactMux) appendSorted(method string, entry ExactMuxEntry) {
	if m.routesSorted == nil {
		m.routesSorted = make(map[string][]ExactMuxEntry)
	}

	if m.routesSorted[method] == nil {
		m.routesSorted[method] = make([]ExactMuxEntry, 0)
	}
	// TODO: rewrite to append
	routes := append(m.routesSorted[method], entry)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].weight > routes[j].weight
	})
	m.routesSorted[method] = routes
}

func (m *ExactMux) handler(method string, path string) (handler http.Handler, err error) {
	entries, exists := m.routes[method]
	if !exists {
		return nil, fmt.Errorf("can't find handler for: %s, %s", method, path)
	}

	if entry, ok := entries[path]; ok {
		return entry.handler, nil
	}

	return nil, fmt.Errorf("can't find handler for: %s, %s", method, path)
}

// ExactMuxEntry ...
type ExactMuxEntry struct {
	pattern string
	handler http.Handler
	weight  int
}

func calculateWeight(pattern string) int {
	if pattern == "/" {
		return 0
	}

	count := (strings.Count(pattern, "/") - 1) * 2
	if !strings.HasSuffix(pattern, "/") {
		return count + 1
	}
	return count
}

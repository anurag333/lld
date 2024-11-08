package middlewarerouter

import (
	"fmt"
	"strings"
)

// Router interface defines the middleware router methods
type Router interface {
	AddRoute(path string, result string)
	CallRoute(path string) (string, error)
}

// SimpleRouter struct will implement the Router interface
type SimpleRouter struct {
	routes map[string]string
}

// NewRouter creates a new instance of SimpleRouter
func NewRouter() *SimpleRouter {
	return &SimpleRouter{
		routes: make(map[string]string),
	}
}

// AddRoute adds a new route and its associated result
func (r *SimpleRouter) AddRoute(path string, result string) {
	// Normalize the path by ensuring it starts with a "/"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	r.routes[path] = result
}

// CallRoute calls a route and returns its associated result
func (r *SimpleRouter) CallRoute(path string) (string, error) {
	// Normalize the path by ensuring it starts with a "/"
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Check if the path exists in the routes map
	if result, exists := r.routes[path]; exists {
		return result, nil
	}

	return "", fmt.Errorf("route not found: %s", path)
}

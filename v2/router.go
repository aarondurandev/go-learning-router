package router

import "net/http"

// Router is the public contract for the router. Any type that implements
// these methods can be used wherever a Router is expected.
type Router interface {
	// Embedding http.Handler means a Router can be passed directly to http.ListenAndServe.
	http.Handler

	// Handle registers a handler for the given pattern, matching any HTTP method.
	Handle(pattern string, h http.Handler)

	// HandleFunc registers a handler function for the given pattern, matching any HTTP method.
	HandleFunc(pattern string, h http.HandlerFunc)

	// Method registers a handler for the given method and pattern.
	Method(method, pattern string, h http.Handler)

	// MethodFunc registers a handler function for the given method and pattern.
	MethodFunc(method, pattern string, h http.HandlerFunc)

	// HTTP verb shortcuts.
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)

	// NotFound sets a custom handler for requests that match no route.
	NotFound(h http.HandlerFunc)

	// Use appends middleware to the router's middleware stack.
	Use(middlewares ...func(http.Handler) http.Handler)

	// Group registers a set of routes under a common prefix.
	Group(prefix string, fn func(Router))
}

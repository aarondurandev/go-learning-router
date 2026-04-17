package router

import "net/http"

// RouteBuilder holds a pattern and a reference to the parent Mux, allowing
// multiple HTTP methods to be registered on the same pattern via method chaining.
// Middleware registered with Use is scoped to this builder and applied only to
// handlers registered through it.
type RouteBuilder struct {
	*Mux
	pattern     string
	middlewares []func(http.Handler) http.Handler
}

// Route returns a RouteBuilder for the given pattern, enabling method chaining.
func (mx *Mux) Route(pattern string) *RouteBuilder {
	return &RouteBuilder{Mux: mx, pattern: pattern}
}

// Get registers a GET handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Get(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Method("GET", rb.pattern, chain(rb.middlewares, h))
	return rb
}

// Post registers a POST handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Post(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Method("POST", rb.pattern, chain(rb.middlewares, h))
	return rb
}

// Put registers a PUT handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Put(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Method("PUT", rb.pattern, chain(rb.middlewares, h))
	return rb
}

// Delete registers a DELETE handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Delete(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Method("DELETE", rb.pattern, chain(rb.middlewares, h))
	return rb
}

// Patch registers a PATCH handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Patch(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Method("PATCH", rb.pattern, chain(rb.middlewares, h))
	return rb
}

// Use appends middleware to this builder's stack. The middleware is applied only
// to handlers registered through this RouteBuilder, not to other routes on the Mux.
func (rb *RouteBuilder) Use(middlewares ...func(http.Handler) http.Handler) *RouteBuilder {
	rb.middlewares = append(rb.middlewares, middlewares...)
	return rb
}

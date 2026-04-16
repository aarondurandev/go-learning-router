package router

import "net/http"

// RouteBuilder holds a pattern and a reference to the parent Mux, allowing
// multiple HTTP methods to be registered on the same pattern via method chaining.
type RouteBuilder struct {
	*Mux
	pattern string
}

// Route returns a RouteBuilder for the given pattern, enabling method chaining.
func (mx *Mux) Route(pattern string) *RouteBuilder {
	return &RouteBuilder{Mux: mx, pattern: pattern}
}

// Get registers a GET handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Get(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Get(rb.pattern, h)
	return rb
}

// Post registers a POST handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Post(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Post(rb.pattern, h)
	return rb
}

// Put registers a PUT handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Put(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Put(rb.pattern, h)
	return rb
}

// Delete registers a DELETE handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Delete(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Delete(rb.pattern, h)
	return rb
}

// Patch registers a PATCH handler and returns the RouteBuilder for further chaining.
func (rb *RouteBuilder) Patch(h http.HandlerFunc) *RouteBuilder {
	rb.Mux.Patch(rb.pattern, h)
	return rb
}

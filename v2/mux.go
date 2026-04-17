package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// contextKey is an unexported type for context keys in this package.
// Using a named type prevents collisions with keys from other packages.
type contextKey string

// paramsKey is the context key under which URL parameters are stored.
const paramsKey contextKey = "params"

// Mux is the router. Routes are stored in a radix tree for O(log n) matching
// and dispatched to the appropriate handler on each request.
// RedirectTrailingSlash is true by default: requests with a mismatched trailing
// slash are redirected (301) to the alternate path when a matching route exists there.
type Mux struct {
	notFoundHandler       http.HandlerFunc
	middlewares           []func(http.Handler) http.Handler
	RedirectTrailingSlash bool
	root                  *node
}

// group is a set of routes sharing a common prefix and middleware stack.
// It implements Router by delegating to the parent Mux, prepending the
// prefix and wrapping handlers with group-level middleware at registration time.
type group struct {
	prefix      string
	mux         *Mux
	middlewares []func(http.Handler) http.Handler
}

// compile-time check that *Mux implements Router.
var _ Router = (*Mux)(nil)

// compile-time check that *group implements Router.
var _ Router = (*group)(nil)

// Method registers a handler for the given HTTP method and pattern.
// The method is normalized to uppercase before storing.
// All other registration methods delegate to this one.
func (mx *Mux) Method(method, pattern string, handler http.Handler) {
	parsedMethod := strings.ToUpper(method)
	mx.root.insert(parsedMethod, pattern, handler)
}

// MethodFunc registers a handler function for the given HTTP method and pattern.
func (mx *Mux) MethodFunc(method, pattern string, handlerFn http.HandlerFunc) {
	mx.Method(method, pattern, handlerFn)
}

// Handle registers a handler for the given pattern, matching any HTTP method.
// Internally the method is stored as an empty string to signal "any method".
func (mx *Mux) Handle(pattern string, handler http.Handler) {
	mx.Method("", pattern, handler)
}

// HandleFunc registers a handler function for the given pattern, matching any HTTP method.
func (mx *Mux) HandleFunc(pattern string, handlerFn http.HandlerFunc) {
	mx.Method("", pattern, handlerFn)
}

func (mx *Mux) Get(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("GET", pattern, handlerFn)
}
func (mx *Mux) Delete(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("DELETE", pattern, handlerFn)
}
func (mx *Mux) Post(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("POST", pattern, handlerFn)
}
func (mx *Mux) Patch(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("PATCH", pattern, handlerFn)
}
func (mx *Mux) Put(pattern string, handlerFn http.HandlerFunc) {
	mx.MethodFunc("PUT", pattern, handlerFn)
}

// NotFound sets a custom handler for requests that match no route.
// If not set, the default net/http not-found handler is used.
func (mx *Mux) NotFound(handlerFn http.HandlerFunc) {
	mx.notFoundHandler = handlerFn
}

// ServeHTTP dispatches the request to the matching route's handler.
// Returns 405 if the path matches but the method does not.
// Returns 404 if no route matches at all.
func (mx *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matched, params := mx.root.search(r.URL.Path)
	if matched != nil {
		h, ok := matched.handlers[r.Method]
		if !ok {
			h = matched.handlers[""]
		}
		if h != nil {
			ctx := context.WithValue(r.Context(), paramsKey, params)
			handler := chain(mx.middlewares, h)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if len(matched.handlers) > 0 {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
	{
		// If RedirectTrailingSlash is enabled, try the path with the slash toggled.
		// If a route matches the alternate path, redirect the client there (301).
		var altPath string
		if strings.HasSuffix(r.URL.Path, "/") {
			altPath = r.URL.Path[:len(r.URL.Path)-1]
		} else {
			altPath = r.URL.Path + "/"
		}
		if mx.RedirectTrailingSlash && altPath != "" {
			trailMatch, _ := mx.root.search(altPath)

			if trailMatch != nil {
				http.Redirect(w, r, altPath, http.StatusMovedPermanently)
				return
			}

		}
		if mx.notFoundHandler != nil {
			mx.notFoundHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}

}

// URLParam returns the value of the URL parameter with the given key
// from the request context. Returns an empty string if not found.
func URLParam(r *http.Request, key string) string {
	params, ok := r.Context().Value(paramsKey).(map[string]string)
	if !ok {
		return ""
	}
	return params[key]
}

// URLParamInt returns the URL parameter for the given key as an int.
// Returns an error if the params are not in the context or the value is not a valid integer.
func URLParamInt(r *http.Request, key string) (int, error) {
	params, ok := r.Context().Value(paramsKey).(map[string]string)
	if !ok {
		return 0, fmt.Errorf("url params not found in context")
	}
	return strconv.Atoi(params[key])
}

// NewMux creates and returns a new Mux instance.
// RedirectTrailingSlash is enabled by default.
func NewMux() *Mux {
	return &Mux{
		RedirectTrailingSlash: true,
		root:                  &node{},
	}
}

// Use appends one or more middleware functions to the router's middleware stack.
// Middleware is applied in registration order — the first registered runs first.
func (mx *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	mx.middlewares = append(mx.middlewares, middlewares...)
}

// chain wraps the final handler with each middleware in order.
// Middleware is applied in reverse so the first registered ends up outermost.
func chain(middlewares []func(http.Handler) http.Handler, final http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		final = middlewares[i](final)
	}
	return final
}

// Group creates a new route group under the given prefix and calls fn with it.
func (mx *Mux) Group(prefix string, fn func(Router)) {
	fn(&group{prefix: prefix, mux: mx})
}

// Method registers a handler for the given method and pattern, prepending the
// group prefix and wrapping the handler with the group's middleware chain.
func (g *group) Method(method, pattern string, h http.Handler) {
	g.mux.Method(method, g.prefix+pattern, chain(g.middlewares, h))
}

func (g *group) MethodFunc(method, pattern string, h http.HandlerFunc) {
	g.Method(method, pattern, h)
}

func (g *group) Handle(pattern string, h http.Handler) {
	g.Method("", pattern, h)
}

func (g *group) HandleFunc(pattern string, h http.HandlerFunc) {
	g.Method("", pattern, h)
}

func (g *group) Get(pattern string, h http.HandlerFunc)    { g.Method("GET", pattern, h) }
func (g *group) Post(pattern string, h http.HandlerFunc)   { g.Method("POST", pattern, h) }
func (g *group) Put(pattern string, h http.HandlerFunc)    { g.Method("PUT", pattern, h) }
func (g *group) Delete(pattern string, h http.HandlerFunc) { g.Method("DELETE", pattern, h) }
func (g *group) Patch(pattern string, h http.HandlerFunc)  { g.Method("PATCH", pattern, h) }

// NotFound delegates to the parent mux — not-found handling is router-wide.
func (g *group) NotFound(h http.HandlerFunc) { g.mux.NotFound(h) }

// ServeHTTP delegates to the parent mux.
func (g *group) ServeHTTP(w http.ResponseWriter, r *http.Request) { g.mux.ServeHTTP(w, r) }

// Use appends middleware to the group's own middleware stack.
// These middlewares only apply to routes registered through this group.
func (g *group) Use(middlewares ...func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Group creates a nested group, prepending this group's prefix to the new prefix.
func (g *group) Group(prefix string, fn func(Router)) {
	g.mux.Group(g.prefix+prefix, fn)
}

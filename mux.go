package router

import (
	"net/http"
	"strings"
)

type Route struct {
	method  string
	pattern string
	handler http.Handler
}

type Mux struct {
	routes          []Route
	notFoundHandler http.HandlerFunc
}

var _ Router = (*Mux)(nil)

func (mx *Mux) Method(method, pattern string, handler http.Handler) {
	parsedMethod := strings.ToUpper(method)
	newRoute := Route{
		method:  parsedMethod,
		pattern: pattern,
		handler: handler,
	}
	mx.routes = append(mx.routes, newRoute)
}

func (mx *Mux) MethodFunc(method, pattern string, handlerFn http.HandlerFunc) {
	mx.Method(method, pattern, handlerFn)
}

func (mx *Mux) Handle(pattern string, handler http.Handler) {
	mx.Method("", pattern, handler)
}

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

func (mx *Mux) NotFound(handlerFn http.HandlerFunc) {
	mx.notFoundHandler = handlerFn
}

func (mx *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var matchFound bool = false
	for _, route := range mx.routes {
		if route.pattern == r.URL.Path && (route.method == r.Method || route.method == "") {
			route.handler.ServeHTTP(w, r)
			return
		} else if route.pattern == r.URL.Path && (route.method != r.Method) {
			matchFound = true
		}

	}
	if matchFound == true {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		if mx.notFoundHandler != nil {
			mx.notFoundHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func NewMux() *Mux {
	return &Mux{}
}

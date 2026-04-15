# go-learning-router

A lightweight HTTP router for Go, built from scratch as a learning project. The goal is to understand how routers work internally — route registration, request dispatching, middleware, and URL parameters — by implementing each piece step by step.

## Installation

```bash
go get github.com/aarondurandev/go-learning-router
```

## Usage

```go
package main

import (
    "fmt"
    "net/http"

    router "github.com/aarondurandev/go-learning-router"
)

func main() {
    m := router.NewMux()

    m.Get("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "home")
    })

    m.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "hello there")
    })

    m.NotFound(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "not found")
    })

    http.ListenAndServe(":8080", m)
}
```

### Available methods

```go
m.Get(pattern, handler)
m.Post(pattern, handler)
m.Put(pattern, handler)
m.Patch(pattern, handler)
m.Delete(pattern, handler)

// Any method
m.Handle(pattern, handler)
m.HandleFunc(pattern, handler)

// Explicit method
m.Method("OPTIONS", pattern, handler)
m.MethodFunc("OPTIONS", pattern, handler)

// Custom not-found handler
m.NotFound(handler)
```

## Design notes

- Routes are stored as a slice and matched linearly — simple and easy to reason about.
- `Handle`/`HandleFunc` register a route that matches any HTTP method (stored as an empty string internally).
- The HTTP verb shortcuts (`Get`, `Post`, etc.) delegate down to `MethodFunc` → `Method`, so all registration goes through one place.
- A request that matches a pattern but not the method gets a `405 Method Not Allowed`, not a `404`.
- `*Mux` satisfies `http.Handler` directly, so it can be passed to `http.ListenAndServe` without any wrapping.

## Roadmap

- [x] Route registration
- [x] Request dispatching (404 / 405 handling)
- [ ] URL parameters (`/users/{id}`)
- [ ] Middleware
- [ ] Subrouters / route groups

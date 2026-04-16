# go-learning-router

A lightweight HTTP router for Go, built from scratch as a learning project. The goal is to understand how routers work internally — route registration, request dispatching, middleware, and URL parameters — by implementing each piece step by step.

## Installation

```bash
go get github.com/aarondurandev/go-learning-router
```

## Usage

See [_examples/](_examples/) for working code.

## Design notes

- Routes are stored as a slice and matched linearly — simple and easy to reason about.
- `Handle`/`HandleFunc` register a route that matches any HTTP method (stored as an empty string internally).
- The HTTP verb shortcuts (`Get`, `Post`, etc.) delegate down to `MethodFunc` → `Method`, so all registration goes through one place.
- A request that matches a pattern but not the method gets a `405 Method Not Allowed`, not a `404`.
- URL parameters use `{name}` syntax and are matched segment by segment. Captured values are stored in the request context and retrieved with `URLParam`.
- Middleware is registered with `Use` and applied to every matched route. Multiple middlewares run in registration order.
- Route groups share a common prefix and can have their own middleware stack via `Use`. Group middleware only applies to routes registered through that group.
- Groups can be nested — a nested group prepends both prefixes automatically.
- `*Mux` satisfies `http.Handler` directly, so it can be passed to `http.ListenAndServe` without any wrapping.
- Wildcard routes use a trailing `*` segment (`/files/*`) and match the rest of the path regardless of depth. The captured tail is retrieved with `URLParam(r, "*")`.
- `RedirectTrailingSlash` (default `true`) — a request to `/users/` redirects (301) to `/users` (and vice versa) if a matching route exists at the alternate path. Set to `false` to disable.
- `Route(pattern)` returns a `RouteBuilder` that supports method chaining — `m.Route("/users").Get(h).Post(h)` registers multiple methods on the same pattern in one expression.

## Roadmap

- [x] Route registration
- [x] Request dispatching (404 / 405 handling)
- [x] URL parameters (`/users/{id}`)
- [x] Middleware
- [x] Subrouters / route groups
- [x] Tests
- [x] Wildcards (`/files/*`)
- [x] Trailing slash redirect
- [x] Method chaining

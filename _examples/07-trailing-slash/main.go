// This example demonstrates trailing slash redirect. When RedirectTrailingSlash
// is true (the default), a request to /users/ is redirected (301) to /users if
// a matching route exists there. Use "curl -L -v localhost:8080/users/" to see
// both hops in the terminal.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.Use(redirectLoggerMiddleware)
	m.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "No users found")
	})
	http.ListenAndServe(":8080", m)
}

func redirectLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Requested: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

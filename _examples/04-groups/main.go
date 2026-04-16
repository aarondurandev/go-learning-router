// This example demonstrates route groups. Routes registered inside a group
// are automatically prefixed, making API versioning straightforward.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.Group("/api/v1", func(r router.Router) {
		r.Get("/users", usersHandlerV1)
	})
	m.Group("/api/v2", func(r router.Router) {
		r.Get("/users", usersHandlerV2)
	})
	http.ListenAndServe(":8080", m)
}

func usersHandlerV1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "users endpoint v1")
}

func usersHandlerV2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "users endpoint v2")
}

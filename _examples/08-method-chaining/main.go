// This example demonstrates method chaining and RouteBuilder middleware.
// Route returns a RouteBuilder that lets you register multiple HTTP methods on
// the same pattern in a single expression. Use registers middleware scoped only
// to that builder — it does not affect other routes on the mux.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()

	// /users has no route-level middleware — requests pass through unmodified.
	m.Route("/users").
		Get(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "GET /users")
		}).
		Post(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "POST /users")
		}).
		Put(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "PUT /users")
		})

	// /admin uses Use to attach middleware before registering handlers.
	// The middleware only applies to routes registered through this builder.
	m.Route("/admin").
		Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("admin middleware fired")
				next.ServeHTTP(w, r)
			})
		}).
		Get(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "GET /admin")
		}).
		Post(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "POST /admin")
		})

	http.ListenAndServe(":8080", m)
}

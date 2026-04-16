// This example demonstrates method chaining. Route returns a RouteBuilder that
// lets you register multiple HTTP methods on the same pattern in a single expression.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
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
	http.ListenAndServe(":8080", m)
}

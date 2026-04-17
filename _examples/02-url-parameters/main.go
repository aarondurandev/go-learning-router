// This example demonstrates URL parameters using {name} syntax.
// String values are retrieved with URLParam; URLParamInt parses the value as an
// integer and returns an error if the value is not a valid integer.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Can't find it")
	})
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "home")
	})

	m.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello there")
	})
	m.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, %s", router.URLParam(r, "name"))
	})
	m.Get("/hello/{name}/{age}", func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")
		age, err := router.URLParamInt(r, "age")
		if err != nil {
			http.Error(w, "age must be a number", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Hello, %s. You're %d years old", name, age)
	})
	http.ListenAndServe(":8080", m)
}

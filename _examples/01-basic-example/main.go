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
	http.ListenAndServe(":8080", m)
}

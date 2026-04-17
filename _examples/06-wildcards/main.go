// This example demonstrates wildcard routes. A trailing * segment matches the rest
// of the path and is retrieved with URLParam(r, "*"). The {name:*} syntax does the
// same but stores the captured tail under the given name, consistent with {name} params.
package main

import (
	"fmt"
	"net/http"

	router "github.com/aarondurandev/go-learning-router"
)

func main() {
	m := router.NewMux()
	m.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		file := router.URLParam(r, "*")
		if file != "" {
			fmt.Fprintf(w, "Requested file: %s", file)
		}
	})
	m.Get("/videos/{path:*}", func(w http.ResponseWriter, r *http.Request) {
		video := router.URLParam(r, "path")
		if video != "" {
			fmt.Fprintf(w, "Requested video: %s", video)
		}
	})
	http.ListenAndServe(":8080", m)
}

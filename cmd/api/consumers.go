// Filename: consumers.go
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// createConsumerHandler for the POST /v1/consumers endpoint
func (app *application) createConsumerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Consumer Handler")
}

// showConsumerHandler for the GET /v1/consumers/:id endpoint
func (app *application) showConsumerHandler(w http.ResponseWriter, r *http.Request) {
	// Use the "ParmsFromContext" function to get the parameters from the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	// GET the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Display the consumer ID
	fmt.Fprintf(w, "Consumer ID: %d\n", id)
}


// Filename: consumers.go
package main

import (
	"fmt"
	"net/http"
	
)

// createConsumerHandler for the POST /v1/consumers endpoint
func (app *application) createConsumerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Consumer Handler")
}

// showConsumerHandler for the GET /v1/consumers/:id endpoint
func (app *application) showConsumerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Display the consumer ID
	fmt.Fprintf(w, "Consumer ID: %d\n", id)
}


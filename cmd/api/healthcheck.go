// Filename: cmd/api/healthcheck.go
package main

import (
	"net/http"
	"fmt"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.Env, version)
	// Specify that we will be serve our responses using JSON
	w.Header().Set("Content-Type", "application/json")
	// write the JSON as the HTTP response body
	w.Write([]byte(js))
}
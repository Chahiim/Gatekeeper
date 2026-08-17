// Filename: cmd/api/healthcheck.go
package main

import (
	"net/http"
	"fmt"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status: 200 OK")
	fmt.Fprintf(w, "Environment: %s \n", app.config.Env)
	fmt.Fprintf(w, "Version: %s \n", version)
	fmt.Fprintf(w, "Current Time: %s \n", r.Context().Value(http.ServerContextKey).(*http.Server).IdleTimeout)
}
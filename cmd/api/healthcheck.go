// Filename: cmd/api/healthcheck.go
package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map to hold our healthcheck data
	data := map[string]string {
		"status": "available",
		"environment": app.config.Env,
		"version": version,
	}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
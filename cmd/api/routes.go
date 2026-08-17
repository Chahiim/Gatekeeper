// Filename: cmd/api/routes.go
package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Create a new httprouter router instance
	router := httprouter.New()
	router.HandleFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandleFunc(http.MethodPost, "/v1/consumers", app.createConsumerHandler)
	router.HandleFunc(http.MethodGet, "/v1/consumers/:id", app.getConsumerHandler)
	return router
}
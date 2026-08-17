// Filename: cmd/api/routes.go
package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Create a new httprouter router instance
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/consumers", app.createConsumerHandler)
	router.HandlerFunc(http.MethodGet, "/v1/consumers/:id", app.showConsumerHandler)
	return router
}
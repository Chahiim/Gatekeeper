package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("POST /v1/consumers", app.createConsumerHandler)
	mux.HandleFunc("GET /v1/consumers/{id}", app.showConsumerHandler)
	mux.HandleFunc("POST /v1/reports", app.createReportHandler)

	return mux
}

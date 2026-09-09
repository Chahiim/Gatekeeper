package main

import (
	"io/fs"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("POST /v1/images", app.createImageHandler)
	mux.HandleFunc("GET /v1/jobs/{id}", app.showJobHandler)
	mux.HandleFunc("GET /v1/images/{imageID}/variants/{name}", app.showVariantHandler)

	mux.HandleFunc("GET /v1/consumers", func(w http.ResponseWriter, r *http.Request) {
		app.writeJSON(w, http.StatusOK, envelope{"consumers": "use POST to create"}, nil)
	})
	mux.HandleFunc("POST /v1/consumers", app.createConsumerHandler)
	mux.HandleFunc("GET /v1/consumers/{id}", app.showConsumerHandler)
	mux.HandleFunc("POST /v1/reports", app.createReportHandler)

	frontendFS := http.Dir("./ui")
	fileServer := http.FileServer(frontendFS)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/index.html" {
			content, err := fs.ReadFile(os.DirFS("./ui"), "index.html")
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(content)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	return mux
}

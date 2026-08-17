//Filename: cmd/api/main.go
package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

// The application version number
const version = "1.0.0"
// The configuration settings
type config struct {
	Port int
	Env  string // development, staging, production
}
// Dependency Injection
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	// Flags needed to populate the config
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	//Create a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	//Create an instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
	}
	// Create our new servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	// Create a new HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// Start the server
	logger.Printf("Starting %s server on %s", cfg.Env, srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}


}
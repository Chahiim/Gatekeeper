// Filename: cmd/api/helpers.go
package main

import (
	"errors"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	// Use the "ParmsFromContext" function to get the parameters from the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	// GET the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	//convert our map into a JSON object
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Add a newline to make viewing on the terminal easier
	js = append(js, '\n')
	// Add the headers
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Specify that we will be serve our responses using JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Write the []byte slice containing the JSON response body
	w.Write(js)	
	return nil
}
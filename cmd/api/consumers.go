package main

import (
	"fmt"
	"net/http"

	"github.com/Chahiim/Gatekeeper/internal/data"
	"github.com/Chahiim/Gatekeeper/internal/validator"
)

func (app *application) createConsumerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string               `json:"name"`
		Email string               `json:"email"`
		Status data.ConsumerStatus `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	consumer := &data.Consumer{
		Name:   input.Name,
		Email:  input.Email,
		Status: input.Status,
	}

	if consumer.Status == "" {
		consumer.Status = data.ConsumerStatusActive
	}

	v := validateConsumer(input.Name, input.Email, input.Status)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Consumers.Insert(consumer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/consumers/%s", consumer.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"consumer": consumer}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func validateConsumer(name, email string, status data.ConsumerStatus) *validator.Validator {
	v := validator.New()
	v.Check(name != "", "name", "must be provided")
	v.Check(len(name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(status == "" || status == data.ConsumerStatusActive || status == data.ConsumerStatusSuspended || status == data.ConsumerStatusTerminated, "status", "must be a valid status")
	return v
}

func (app *application) showConsumerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	consumer, err := app.models.Consumers.Get(id)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"consumer": consumer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

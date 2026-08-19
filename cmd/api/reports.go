package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Chahiim/Gatekeeper/internal/data"
	"github.com/Chahiim/Gatekeeper/internal/validator"
)

func (app *application) createReportHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConsumerID string    `json:"consumer_id"`
		From       time.Time `json:"from"`
		To         time.Time `json:"to"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.ConsumerID != "", "consumer_id", "must be provided")
	v.Check(!input.From.IsZero(), "from", "must be provided")
	v.Check(!input.To.IsZero(), "to", "must be provided")
	v.Check(input.From.Before(input.To), "from", "must be earlier than to")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	app.logger.Info(
		"report generation started",
		"consumer_id", input.ConsumerID,
		"artificial_delay", app.config.reportDelay,
	)

	time.Sleep(app.config.reportDelay)

	report, err := app.models.Reports.Generate(input.ConsumerID, input.From, input.To)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.Info(
		"report generation finished",
		"consumer_id", input.ConsumerID,
	)

	err = app.writeJSON(w, http.StatusOK, envelope{"report": report}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

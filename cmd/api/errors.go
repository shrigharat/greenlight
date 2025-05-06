package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		path   = r.URL.Path
	)
	app.logger.Error(err.Error(), "method", method, "path", path)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message any, status int) {
	errorJson := envelope{"error": message}
	err := app.writeJSON(w, status, errorJson, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, message, http.StatusInternalServerError)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, message, http.StatusNotFound)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := err.Error()
	app.errorResponse(w, r, message, http.StatusBadRequest)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, errors, http.StatusUnprocessableEntity)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not allowed for this resource", r.Method)
	app.errorResponse(w, r, message, http.StatusMethodNotAllowed)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, message, http.StatusConflict)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, message, http.StatusTooManyRequests)
}
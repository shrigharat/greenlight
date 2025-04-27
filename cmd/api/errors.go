package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		path = r.URL.Path
	)
	app.logger.Error(err.Error(), "method", method, "path", path)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message string, status int) {
	errorJson := envelope{"error": message}
	err := app.writeJSON(w, status, errorJson, nil)
	if err!=nil {
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

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not allowed for this resource", r.Method)
	app.errorResponse(w, r, message, http.StatusMethodNotAllowed)
}
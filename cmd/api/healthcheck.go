package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available",
		"version": version,
		"env": app.config.env,
	}
	err := app.writeJSON(w, http.StatusOK, envelope{"data": data}, nil)
	if err!=nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
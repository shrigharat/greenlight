package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// panic handler goes here
		// this is a defered function, so it will always run in the event of a panic
		// as Go unwinds the stack
		defer func() {
			if err := recover(); err !=nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w,r)
	})
}